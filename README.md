# Storage-deal Microservice
Generic Deal Making Service using whypfs + filclient + estuary_auth

## Features
- Creates a deal for files that are 4GB and above
- Shows all the deals made for specific user
- Built-in gateway to view the files while the background process creates the deals for each

## What this service doesn't do
- It won't accept any files lower than 4GB
- It won't prepare CAR files or split files for users


This is stricly a deal-making service. It'll pin the files/CIDs but it won't keep the files. Once a deal is made, the CID will be removed from the blockstore. For retrievals, use the retrieval-deal microservice.
