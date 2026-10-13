# :zap: Meilisearch Caching Proxy :zap:
----------------------------

## Description

This is a minimal caching proxy for Meilisearch. It's extremely light-weight and performant.

### Why did we build this?

We had issues scaling meilisearch. Meilisearch does not scale out properly, you can scale vertically but not horizontally.
This proxy can cache search requests in a compatible memory store, greatly reducing the amount of requests that hit the MeiliSearch instance.


### Features

* :zap: Fast and minimal reverse proxy
* :floppy_disk: Proxy and cache search requests (POST and GET)
* :takeout_box: Secure purge API per index or globally

It supports the following caching engines:

* in-memory (ristretto)
* Redis

Tested against the following MeiliSearch versions:
* v1.9

We only test with the self-hosted version, it will work with Meilisearch cloud but we cannot guarantee it.

### Usage

You can install the proxy in Kubernetes using our helm chart:

```
helm pull oci://ghcr.io/maxroll-media-group/meilisearch-proxy/charts/meilisearch-proxy --version 0.2.1

helm install meilisearch-proxy oci://ghcr.io/maxroll-media-group/meilisearch-proxy/charts/meilisearch-proxy --version 0.2.1 \
--set meilisearch.url=http://meilisearch:7700 \
--set meilisearch.masterKey=<master_key> \
--set proxy.purgeToken=<purge_secret> \
```

You can also run the docker container directly:

```
docker run -p 7700:7700 -e MEILISEARCH_HOST=http://meilisearch-endpoint MEILISEARCH_MASTER_KEY=xxxx  -it ghcr.io/maxroll-media-group/meilisearch-proxy:latest
```
