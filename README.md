# Storage-deal Microservice
Generic Deal Making Service using whypfs + filclient + estuary_auth

## Features
- Creates a deal for large files. The recommended size is 1GB. 
- Shows all the deals made for specific user
- Built-in gateway to view the files while the background process creates the deals for each

This is strictly a deal-making service. It'll pin the files/CIDs but it won't keep the files. Once a deal is made, the CID will be removed from the blockstore. For retrievals, use the retrieval-deal microservice.

## Standalone mode
- Run a single instance of the deal-maker microservice. This will use a local SQLite database and local file system for the blockstore.

## HA mode
- Run multiple instances of the deal-maker microservice all pointing to HA Postgres and Centralize/Shared filesystem.

