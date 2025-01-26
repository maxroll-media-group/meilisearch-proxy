import http from 'k6/http';

function getRandomInt(min, max) {
  return Math.floor(Math.random() * (max - min + 1)) + min;
}

function sleep(ms) {
	  return new Promise(resolve => setTimeout(resolve, ms));
}

export default async function () {
  const url = 'http://localhost:7700/indexes/test_index/search';

  const params = {
    headers: {
      'Content-Type': 'application/json',
	  'Authorization': 'Bearer test'
    },
  }

  while(true) {
	http.post(url, JSON.stringify({
		limit: getRandomInt(1, 100),
	  }), params);

	  await sleep(getRandomInt(1, 100));

	  // do a purge every 100 requests
	  if (getRandomInt(1, 100) == 1) {
		  http.post('http://localhost:7700/purge/test_index', null, params);
	  }
  }
}
