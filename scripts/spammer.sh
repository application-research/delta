#!/bin/bash

# This script is used to spam the edge with 50kb files
while :
do
  ms=$(date +%s%N)
  dd if=/dev/random of=random_"$ms".dat bs=10000 count=200
	curl --location --request POST 'http://shuttle-4-bs1.estuary.tech:1414/api/v1/content/add' \
  --header 'Authorization: Bearer [REDACTED]' \
  --form 'data=@"./random_'${ms}'.dat"'
  rm random_$ms.dat
done