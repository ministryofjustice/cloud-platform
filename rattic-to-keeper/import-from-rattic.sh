#!/bin/bash
set -e
wget -q --no-check-certificate --config=wget-json -O rattic/groups.json 'https://rattic.service.dsd.io/api/v1/group/?limit=1000&format=json'
wget -q --no-check-certificate --config=wget-json -O rattic/secrets.json 'https://rattic.service.dsd.io/api/v1/cred/?limit=1000&format=json'
wget -q --no-check-certificate --config=wget-json -O rattic/tags.json 'https://rattic.service.dsd.io/api/v1/tag/?limit=1000&format=json'
for id in $(jq -rM '.objects[] | .id' secrets.json) ; do
  uri="https://rattic.service.dsd.io/api/v1/cred/$id/?format=json"
  wget -q --no-check-certificate --config=wget-json -O rattic/secret-$id.json $uri
  uri="https://rattic.service.dsd.io/cred/detail/$id/"
  wget -q --no-check-certificate --config=wget-html -O rattic/detail-$id.html $uri
done
