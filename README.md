[![CodeQL](https://img.shields.io/github/actions/workflow/status/application-research/delta/codeql.yml?label=CodeQL&style=for-the-badge)](https://github.com/application-research/delta/actions/workflows/codeql.yml)
[![Release](https://img.shields.io/github/v/release/application-research/delta?display_name=release&style=for-the-badge)](https://github.com/application-research/delta/releases)

# Δ Delta
Filecoin Storage Deal Making Service

Delta is a deal-making service that enables users to make deals with Storage Providers. It is an application that allows users to upload files to the Filecoin network and get them stored by Storage Providers.

*Delta is in active development and is not ready for production use.*

<div align="center">

## Quick stats

[![](https://img.shields.io/badge/dynamic/json?style=for-the-badge&label=Total%20Content%20processed&color=brightgreen&query=total_content_consumed&url=https%3A%2F%2Fnode.delta.store%2Fopen%2Fstats%2Ftotals%2Finfo)]()
[![](https://img.shields.io/badge/dynamic/json?style=for-the-badge&label=Total%20end-to-end%20deals&query=total_e2e_deals&color=brightgreen&url=https%3A%2F%2Fnode.delta.store%2Fopen%2Fstats%2Ftotals%2Finfo)]()
[![](https://img.shields.io/badge/dynamic/json?style=for-the-badge&label=Total%20commp%20made&color=brightgreen&query=total_piece_commitment_made&url=https%3A%2F%2Fnode.delta.store%2Fopen%2Fstats%2Ftotals%2Finfo)]()
[![](https://img.shields.io/badge/dynamic/json?style=for-the-badge&label=Total%20import%20deals&color=brightgreen&query=total_import_deals&url=https%3A%2F%2Fnode.delta.store%2Fopen%2Fstats%2Ftotals%2Finfo)]()
[![](https://img.shields.io/badge/dynamic/json?style=for-the-badge&label=Total%20Content%20processed%20in-bytes%20&color=brightgreen&query=total_storage_allocated&url=https%3A%2F%2Fnode.delta.store%2Fopen%2Fstats%2Ftotals%2Finfo)]()
[![](https://img.shields.io/badge/dynamic/json?style=for-the-badge&label=Total%20e2e%20in-bytes&color=brightgreen&query=total_e2e_deals_in_bytes&url=https%3A%2F%2Fnode.delta.store%2Fopen%2Fstats%2Ftotals%2Finfo)]()
[![](https://img.shields.io/badge/dynamic/json?style=for-the-badge&label=Total%20import%20in-bytes&color=brightgreen&query=total_import_deals_in_bytes&url=https%3A%2F%2Fnode.delta.store%2Fopen%2Fstats%2Ftotals%2Finfo)]()
[![](https://img.shields.io/badge/dynamic/json?style=for-the-badge&label=Proud%20to%20work%20with%20SPs&color=brightgreen&query=total_miners&url=https%3A%2F%2Fnode.delta.store%2Fopen%2Fstats%2Ftotals%2Finfo)]()


</div>

---

![image](https://user-images.githubusercontent.com/4479171/226853191-e19e8fa4-abc1-4652-970f-d3d6cea0df13.png)

For more information, check out the [docs](docs)

---

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

## Install and run `Delta` using `docker`

Running this the first time will generate a wallet. Make sure to get FIL/DataCap from the [faucet](https://verify.glif.io/) and fund the wallet

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
- go 1.19
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
