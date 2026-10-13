package proxy_test

import (
	"io"
	"net/http"

	"github.com/alicebob/miniredis/v2"
	"github.com/eko/gocache/lib/v4/store"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/maxroll-media-group/meilisearch-proxy/pkg/config"
	"github.com/maxroll-media-group/meilisearch-proxy/pkg/proxy"
)

const testJSON = `{"name":"test"}`
const testIndexJSON = `
{
	"results": [
		{
			"uid": "profiles",
			"createdAt": "2024-08-07T22:32:45.025568204Z",
			"updatedAt": "2024-08-11T14:39:20.410931753Z",
			"primaryKey": "id"
		}]
}
`

var _ = Describe("Proxy", Ordered, func() {

	var redis *miniredis.Miniredis
	var fakeMeilisearch *http.Server
	var addr string
	var proxyServer *proxy.Proxy

	// start a fake Meilisearch server

	BeforeAll(func() {
		redis, _ = miniredis.Run()
		addr = "redis://" + redis.Addr()

		mux := http.NewServeMux()
		mux.HandleFunc("/indexes/test/search", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(testJSON))
			w.WriteHeader(http.StatusOK)
		})

		mux.HandleFunc("/indexes", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(testIndexJSON))
			w.WriteHeader(http.StatusOK)
		})

		fakeMeilisearch = &http.Server{
			Addr:    "localhost:7777",
			Handler: mux,
		}

		go fakeMeilisearch.ListenAndServe()

		cfg := &config.Config{
			MeilisearchHost:        "http://localhost:7777",
			MeilisearchMasterKey:   "masterKey",
			ProxyMasterKey:         "proxyMasterKey",
			ProxyPurgeToken:        "token",
			ProxyMasterKeyOverride: false,
			Port:                   "8888",
			MaxFailures:            5,
			MinSuccesses:           2,
			ExitOnMaxFailures:      false,
			CacheConfig: &config.CacheConfig{
				TTL:    300,
				Engine: "redis",
				Url:    addr,
			},
		}

		proxyServer = proxy.NewProxy(cfg)
		go proxyServer.Listen()
	})

	Context("ProxyCalls", func() {

		It("should proxy POST indexes", func() {
			// create a request
			req, _ := http.NewRequest("POST", "http://localhost:8888/indexes/test/search", nil)

			// req should contain "{"name":"test"}"
			resp, err := http.DefaultClient.Do(req)
			Expect(err).To(BeNil())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))

			resBody, err := io.ReadAll(resp.Body)

			Expect(err).To(BeNil())

			// resp should contain "{"name":"test"}"
			Expect(resBody).To(Equal([]byte(testJSON)))
		})
	})

	It("should have test data in cache", func() {

		// create a request
		_, _ = http.NewRequest("GET", "http://localhost:8888/indexes/test/search", nil)

		cached, err := proxyServer.GetCache().Get(proxyServer.Context, "69166bc619a4b7d7b518c76d46e73d10c2a1dae9baf31aaf4254906582534213")

		Expect(err).To(BeNil())

		Expect(cached).To(Equal(testJSON))
	})

	It("should simply proxy other requests", func() {

		// create a request
		req, _ := http.NewRequest("GET", "http://localhost:8888/indexes", nil)

		resp, err := http.DefaultClient.Do(req)
		Expect(err).To(BeNil())

		Expect(resp.StatusCode).To(Equal(http.StatusOK))

		resBody, err := io.ReadAll(resp.Body)

		Expect(err).To(BeNil())

		Expect(resBody).To(Equal([]byte(testIndexJSON)))
	})

	It("should purge cache on POST /purge", func() {

		cache := proxyServer.GetCache()
		cache.Set(proxyServer.Context, "key", "value")

		// create a request
		req, _ := http.NewRequest("POST", "http://localhost:8888/purge", nil)
		req.Header.Set("Authorization", "Bearer token")

		resp, err := http.DefaultClient.Do(req)

		Expect(err).To(BeNil())

		Expect(resp.StatusCode).To(Equal(http.StatusOK))

		_, err = cache.Get(proxyServer.Context, "key")

		Expect(err).ToNot(BeNil())
	})

	It("should 401 when a wrong purge token is used", func() {

		// create a request
		req, _ := http.NewRequest("POST", "http://localhost:8888/purge", nil)

		resp, err := http.DefaultClient.Do(req)

		Expect(err).To(BeNil())

		Expect(resp.StatusCode).To(Equal(http.StatusUnauthorized))
	})

	It("should only purge a specific index on POST /purge/test", func() {

		cache := proxyServer.GetCache()

		cache.Set(proxyServer.Context, "<index_test>", "value", store.WithTags([]string{"test"}))
		cache.Set(proxyServer.Context, "<index_test2>", "value", store.WithTags([]string{"test2"}))

		// create a request
		req, _ := http.NewRequest("POST", "http://localhost:8888/purge/test", nil)
		req.Header.Set("Authorization", "Bearer token")

		resp, err := http.DefaultClient.Do(req)

		Expect(err).To(BeNil())

		Expect(resp.StatusCode).To(Equal(http.StatusOK))

		_, err = cache.Get(proxyServer.Context, "<index_test>")

		Expect(err).ToNot(BeNil())

		val, err := cache.Get(proxyServer.Context, "<index_test2>")

		Expect(err).To(BeNil())
		Expect(val).To(Equal("value"))
	})

	AfterAll(func() {
		fakeMeilisearch.Close()
		redis.Close()
	})

})
