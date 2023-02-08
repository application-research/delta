# Storage-deal Microservice
Generic Deal Making Service using whypfs + filclient + estuary_auth

## Features
- Creates a deal for large files. The recommended size is 1GB. 
- Shows all the deals made for specific user
- Built-in gateway to view the files while the background process creates the deals for each

This is strictly a deal-making service. It'll pin the files/CIDs but it won't keep the files. Once a deal is made, the CID will be removed from the blockstore. For retrievals, use the retrieval-deal microservice.

## Configuration

Create the .env file in the root directory of the project. The following are the required fields.
```
# Database configuration
MODE=standalone # HA
DB_NAME=stg-deal-maker
REPO=/mnt/.whypfs # shared mounted repo

# Piece commp and deal-making job frequency in seconds
PIECE_COMMP_JOB_FREQ=300 // runs every 5 mins
REPLICATION_JOB_FREQ=600 // runs every 10 mins
```

Running this the first time will generate a wallet. Make sure to get FIL from the [faucet](https://verify.glif.io/) and fund the wallet

## Standalone mode
- Run a single instance of the deal-maker microservice. This will use a local SQLite database and local file system for the blockstore.
- Enable this option by setting .env `MODE=standalone`

## HA mode
- Run multiple instances of the deal-maker microservice all pointing to HA Postgres and Centralize/Shared filesystem.
- Enable this option by setting .env `MODE=HA` and `DB_NAME` to the name of the database and REPO to the shared filesystem.

## Build and run
```
go build -tags netgo -ldflags '-s -w' -o stg-dealer
./stg-dealer daemon --repo=.whypfs
```

## Endpoints
### Upload a file
```
curl --location --request POST 'http://localhost:1313/api/v1/content/add' \
--header 'Authorization: Bearer [ESTUARY_API_KEY]' \
--form 'data=@"random_1675815458N.dat"'
```
