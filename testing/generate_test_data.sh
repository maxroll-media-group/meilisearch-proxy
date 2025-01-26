#!/bin/bash

# Set variables
MEILISEARCH_HOST="http://127.0.0.1:7700"
INDEX_NAME="test_index"
DOCUMENTS='[
  { "id": 1, "title": "Document 1", "content": "This is the first document." },
  { "id": 2, "title": "Document 2", "content": "This is the second document." },
  { "id": 3, "title": "Document 3", "content": "This is the third document." }
]'

# Create index
curl -X POST "$MEILISEARCH_HOST/indexes" \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer test" \-H "Authorization: Bearer test" \
     -d "{\"uid\":\"$INDEX_NAME\"}" &> /dev/null

# Add documents to the index
curl -X POST "$MEILISEARCH_HOST/indexes/$INDEX_NAME/documents" \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer test" \
     -d "$DOCUMENTS" &> /dev/null

echo "Index and documents created successfully."