# Δ Delta
Generic DealMaking MicroService using whypfs + filclient + estuary_auth

![image](https://user-images.githubusercontent.com/4479171/218267752-9a7af133-4e36-4f4c-95da-16b3c7bd73ae.png)


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
DB_DSN=stg-deal-maker
#REPO=/mnt/.whypfs # shared mounted repo

# Frequencies
DISPATCH_JOBS_EVERY=10
MAX_DISPATCH_WORKERS=5000
MAX_CLEANUP_WORKERS=1500
```

Running this the first time will generate a wallet. Make sure to get FIL from the [faucet](https://verify.glif.io/) and fund the wallet

## Standalone mode
- Run a single instance of the deal-maker microservice. This will use a local SQLite database and local file system for the blockstore.
- Enable this option by setting .env `MODE=standalone`

## HA mode
- Run multiple instances of the deal-maker microservice all pointing to HA Postgres and Centralize/Shared filesystem.
- Enable this option by setting .env `MODE=HA` and `DB_NAME` to the name of the database and REPO to the shared filesystem.

![image](https://user-images.githubusercontent.com/4479171/217404957-21fd15be-f0c8-4bd2-a5c6-a2770c5c1db1.png)


## Install the following pre-req
- go 1.18
- [jq](https://stedolan.github.io/jq/)
- [hwloc](https://www.open-mpi.org/projects/hwloc/)
- opencl
- [rustup](https://rustup.rs/)
- postgresql

Alternatively, if using Ubuntu, you can run the following commands to install the pre-reqs
```
apt-get update && \
apt-get install -y wget jq hwloc ocl-icd-opencl-dev git libhwloc-dev pkg-config make && \
apt-get install -y cargo
```

## Build and run

### Using `go` lang
```
go build -tags netgo -ldflags '-s -w' -o stg-dealer
./stg-dealer daemon --repo=.whypfs
```

### Using `docker`
```
docker build -t stg-dealer .
docker run -it --rm -p 1414:1414 stg-dealer
```

## Endpoints

### Node information
To get the node information, use the following endpoints
```
curl --location --request GET 'http://localhost:1414/open/node/info'
curl --location --request GET 'http://localhost:1414/open/node/peers'
curl --location --request GET 'http://localhost:1414/open/node/host'
```

### Upload a file
Use the following endpoint to upload a file. The process will automatically compute the piece size and initiate the deal proposal
and transfer
```
curl --location --request POST 'http://localhost:1414/api/v1/content/add' \
--header 'Authorization: Bearer [ESTUARY_API_KEY]' \
--form 'data=@"random_1675815458N.dat"'
```

### Upload a file with a specific pad piece size, duration, miner and connection mode
Use the following endpoint to upload a file with a specific miner, duration, piece size and connection mode.
```
curl --location --request POST 'http://localhost:1414/api/v1/content/add' \
--header 'Authorization: Bearer [ESTUARY_API_KEY]' \
--form 'data=@"random_1675988961N.dat"' \
--form 'miner="f01963614"' \
--form 'commp="{\"piece\":\"bafybeie6sk45lkml45a6s2ftygkt4h2dh5pcpq3pa7dgl64sl7ycxouczq\",\"size\":2000000,\"duration\":20000000}"' \
--form 'connection_mode="online"' 
```

### Upload a CID with a specific pad piece size, duration, miner and connection mode
```
curl --location --request POST 'http://localhost:1414/api/v1/content/add-cid' \
--header 'Authorization: Bearer [ESTUARY_API_KEY]' \
--form 'cid="bafybeie6sk45lkml45a6s2ftygkt4h2dh5pcpq3pa7dgl64sl7ycxouczq"' \
--form 'miner="f01963614"' \
--form 'duration="1494720"' \
--form 'commp="{\"piece\":\"baga6ea4seaqfevzln75frz5bm74wtrddfmr2akcdww4nu45jvewg74xyzva4udi\",\"padded_piece_size\":268435456}"' \
--form 'connection_mode="offline"' \
--form 'size="200019978"'
```

### Stats (content, commps and deals) 
```
curl --location --request GET 'http://localhost:1414/api/v1/stats' \
--header 'Authorization: Bearer [ESTUARY_API_KEY]' \
```

### Stats of a specific content
When you upload, it returns a content id, use that to get the stats of a specific content
```
curl --location --request GET 'http://localhost:1414/api/v1/stats/content/1' \
--header 'Authorization: Bearer [ESTUARY_API_KEY]'
```