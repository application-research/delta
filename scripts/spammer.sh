#!/bin/bash

# This script is used to spam the edge with 50kb files
while :
do
  ms=$(date +%s%N)
  dd if=/dev/random of=random_"$ms".dat bs=10000 count=500000
	curl --location --request POST 'http://localhost:1313/api/v1/content/add' \
  --header 'Authorization: Bearer EST18096aa7-e0cc-4a2e-8a03-e595fe534b14ARY' \
  --form 'data=@"./random_'${ms}'.dat"'
  rm random_$ms.dat
done