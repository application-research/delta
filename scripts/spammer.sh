#!/bin/bash

# This script is used to spam the edge with 50kb files
while :
do
  ms=$(date +%s%N)
  dd if=/dev/random of=random_"$ms".dat bs=10000 count=250000
	curl --location --request POST 'http://localhost:1414/api/v1/content/add' \
  --header 'Authorization: Bearer [ESTUARY_API_KEY]' \
  --form 'data=@"./random_'${ms}'.dat"' \
  --form 'connection_mode="online"' \
  rm -rf random_$ms.dat
done