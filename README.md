# Δ Delta
Filecoin Storage Deal Making Service

*Delta is in active development and is not ready for production use.*

![image](https://user-images.githubusercontent.com/4479171/226853191-e19e8fa4-abc1-4652-970f-d3d6cea0df13.png)


For more information, check out the [docs](docs/overview.md)

## Quick set-up: Build and Run Delta

Copy the `.env.example` file to `.env` and update the values as needed.

```
# Node info
NODE_NAME=stg-deal-maker
NODE_DESCRIPTION=Experimental Deal Maker
NODE_TYPE=delta-main

# Database configuration
MODE=standalone
DB_DSN=delta.db
#REPO=/mnt/.whypfs # shared mounted repo

# Frequencies
MAX_CLEANUP_WORKERS=1500
```

Running this the first time will generate a wallet. Make sure to get FIL/DataCap from the [faucet](https://verify.glif.io/) and fund the wallet

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

### Using `make` lang
```
make all
./delta daemon --repo=.whypfs --wallet-dir=<walletdir>
```

### Using `go` lang
```
go build -tags netgo -ldflags '-s -w' -o delta
./delta daemon --repo=.whypfs --wallet-dir=<walletdir>
```

### Using `docker`
```
docker build -t delta .
docker run -it --rm -p 1414:1414 delta --repo=.whypfs --wallet-dir=<walletdir>
```

## Running Delta
```
./delta daemon --mode=standalone
```

## Test the API server
Try the following endpoints to test the API server
```
curl --location --request GET 'http://localhost:1414/open/node/info'
curl --location --request GET 'http://localhost:1414/open/node/peers'
curl --location --request GET 'http://localhost:1414/open/node/host'
```

If it return the following, then the API server is working
```
{"name":"stg-deal-maker","description":"Experimental Deal Maker","type":"delta-main"}
```

# Getting Started with `Delta`
- To get started on running delta, go to the [getting started to run delta](docs/getting-started-run-delta.md)
- To get started on using a live delta instance, go to the [getting started to use delta](docs/getting-started-use-delta.md)
- To learn more about deployment modes, go to the [deployment modes](docs/deployment-modes.md)
- To get estuary api key, go to the [estuary api keys](docs/getting-estuary-api-key.md)
- To manage wallets, go to the [managing wallets](docs/manage-wallets.md)
- To make an end-to-end deal, go to the [make e2e deals](docs/make-e2e-deal.md)
- To make an import deal, go to the [make import deals](docs/make-import-deal.md)
- To learn how to repair a deal, go to the [repairing and retrying deals](docs/repair.md)
- To learn how to access the open statistics and information, go to the [open statistics and information](docs/open-stats-info.md)
- To learn about the content lifecycle and check status of the deals, go to the [content lifecycle and deal status](docs/deal-status.md)

## Author
Protocol Labs Outercore Engineering.
