#!/bin/bash

# Define the Meilisearch host and API key
MEILISEARCH_HOST="http://127.0.0.1:7700"

# Define the test index name
TEST_INDEX="test_index"

# Delete the test index
curl -X DELETE "$MEILISEARCH_HOST/indexes/$TEST_INDEX" \
  -H "Authorization: Bearer test" &> /dev/null

echo "Test index '$TEST_INDEX' deleted."
