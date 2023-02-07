# Storage-deal Microservice
Generic Deal Making Service using whypfs + filclient + estuary_auth

## Features
- Creates a deal for files
- Shows all the deals made for specific user
- Built-in gateway to view the files while the background process creates the deals for each

This is stricly a deal-making service. It'll pin the files/CIDs but it won't keep the files. Once a deal is made, the CID will be removed from the blockstore. For retrievals, use the retrieval-deal microservice.
