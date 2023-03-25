# Δ Delta
Filecoin Storage Deal Making Service

*Delta is in active development and is not ready for production use.*

![image](https://user-images.githubusercontent.com/4479171/226853191-e19e8fa4-abc1-4652-970f-d3d6cea0df13.png)


For more information, check out the [docs](docs)

# Build and Run Delta

## Prepare to `build`
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

## Install and run `Delta` using `docker`
Make sure you have docker installed on your machine.

### Run the current delta clone using docker-compose
If you already have a wallet with datacap, you can pass it to the command below. This copies over the wallet directory to the containerized app and sets it as the default wallet.
```
make docker-compose-run WALLET_DIR=<walletdir>
```

**Check localhost**

You can check the localhost to see if the delta app is running
```
curl --location --request GET 'http://localhost:1414/open/node/info'
```

**Next**

Now that you can access a live Delta node, you are now ready to make a deal. You can now go to the following guides:

- [Make an e2e deal](docs/make-e2e-deal.md)
- [Make an import deal](docs/make-import-deal.md)


### Run a specific docker image tag from the docker hub artifactory
We have a few pre-built images on the docker hub artifactory. You can run the following command to run a specific image tag
```
cd docker
./run.sh <TAG> 
```
Note: no tag means it'll just get the latest image.

## Install, build from source and run `Delta`
### Install the following pre-req
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

### Using `make`
```
make all
./delta daemon --repo=.whypfs --wallet-dir=<walletdir>
```

### Using `go build`
```
go build -tags netgo -ldflags '-s -w' -o delta
./delta daemon --repo=.whypfs --wallet-dir=<walletdir>
```

## Run `Delta`
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
To get started with `Delta`, you can follow the following guides [here](docs)

## Author
Protocol Labs Outercore Engineering.
