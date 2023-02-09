# Δ Delta
Generic DealMaking MicroService using whypfs + filclient + estuary_auth

![image](https://user-images.githubusercontent.com/4479171/217404677-7fca404c-a89a-48b4-bc83-3f223dd6508d.png)

## Features
- Creates a deal for large files. The recommended size is 1GB. 
- Shows all the deals made for specific user
- Built-in gateway to view the files while the background process creates the deals for each

This is strictly a deal-making service. It'll pin the files/CIDs but it won't keep the files. Once a deal is made, the CID will be removed from the blockstore. For retrievals, use the retrieval-deal microservice.

## Process Flow
- client upload files
- service queues the request for content or commp
- dispatcher runs every N seconds to check the request

## Configuration

Create the .env file in the root directory of the project. The following are the required fields.
```
# Node info
NODE_NAME=stg-deal-maker
NODE_DESCRIPTION=Experimental Deal Maker
NODE_TYPE=delta-main

# Database configuration
MODE=standalone # HA
DB_NAME=stg-deal-maker
#REPO=/mnt/.whypfs # shared mounted repo

# Frequencies
DISPATCH_JOBS_EVERY=10
MAX_DISPATCH_WORKERS=100
MINER_INFO_UPDATE_JOB_FREQ=300
```

Running this the first time will generate a wallet. Make sure to get FIL from the [faucet](https://verify.glif.io/) and fund the wallet

## Standalone mode
- Run a single instance of the deal-maker microservice. This will use a local SQLite database and local file system for the blockstore.
- Enable this option by setting .env `MODE=standalone`

## HA mode
- Run multiple instances of the deal-maker microservice all pointing to HA Postgres and Centralize/Shared filesystem.
- Enable this option by setting .env `MODE=HA` and `DB_NAME` to the name of the database and REPO to the shared filesystem.

![image](https://user-images.githubusercontent.com/4479171/217404957-21fd15be-f0c8-4bd2-a5c6-a2770c5c1db1.png)


## Build and run
```
go build -tags netgo -ldflags '-s -w' -o stg-dealer
./stg-dealer daemon --repo=.whypfs
```

## Endpoints

### Node information
```
curl --location --request GET 'http://localhost:1313/api/v1/open/node/info'
curl --location --request GET 'http://localhost:1313/api/v1/open/node/peers'
curl --location --request GET 'http://localhost:1313/api/v1/open/node/host'
```

### Upload a file
```
curl --location --request POST 'http://localhost:1313/api/v1/content/add' \
--header 'Authorization: Bearer [ESTUARY_API_KEY]' \
--form 'data=@"random_1675815458N.dat"'
```

### Stats (content, commps and deals) 
```
curl --location --request GET 'http://localhost:1313/api/v1/status/deals' \
--header 'Authorization: Bearer [ESTUARY_API_KEY]' \
--form 'data=@"/Users/alvinreyes/Projects/arg/fc-deal-making-service/files/random_1675818258N.dat"'
```