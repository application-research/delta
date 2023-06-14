# Getting Started
Delta is a deal-making service that enables users to make deals with Storage Providers. It is an application that allows users to upload files to the Filecoin network and get them stored by Storage Providers.

# Minimum hardware requirements
- Processor: Any multi-core processor with a clock speed of at least 1.5 GHz. AMD based processors recommended (for commp).
- RAM: RAM requirement will vary based on usage but At least 1 GB of RAM to run and process small datasets.
- Hard Drive: Storage requirement will vary based on the dataset that Delta will process. At least 50 GB of storage space.
- Operating System: Tested on MacOS, Linux. Reportedly works as expected on Windows. Most of the Production Deltas are running on Ubuntu latest version.
- Networking: A standard network interface card (NIC) with a speed of at least 100 Mbps. May vary depending on how large the data will be since online/e2e deals can use considerable amount of egress.
  - Ingress: Expose port 1414 (this is the port for the REST API endpoints).
  - Ingress/Egress: Expose port 6745. This is the announce address port to allow the Delta peer to connect to different IPFS peer nodes.


**Note that these requirements may vary depending on the specific needs of your application.**

# Running a delta node.
A delta node is a light ipfs node that can run on set of machines and architectures. It's written in go and can be run on any machine that can run go.

To run a delta node, you need to have the following pre-reqs installed on your machine:
- go 1.18
- [jq](https://stedolan.github.io/jq/)
- [hwloc](https://www.open-mpi.org/projects/hwloc/)
- opencl
- [rustup](https://rustup.rs/)

Alternatively, if using Ubuntu, you can run the following commands to install the pre-reqs
```
apt-get update && \
apt-get install -y wget jq hwloc ocl-icd-opencl-dev git libhwloc-dev pkg-config make && \
apt-get install -y cargo
```

# Clone the `Delta` repo
Clone the repo to your machine
```
git clone github.com/application-research/delta
```

# Prepare the .env file.
Copy the `.env.example` file to `.env` and update the values as needed.
*Note of the LOTUS_API variable. It needs to point to calibration network*

```
# Node info
NODE_NAME=delta-node
NODE_DESCRIPTION=Experimental Deal Maker
NODE_TYPE=delta-main

# Database configuration
MODE=standalone
DB_DSN=delta.db
DELTA_AUTH=[NODE_API_KEY_HERE]

# Frequencies
MAX_CLEANUP_WORKERS=1500

# APIs
LOTUS_API=http://api.calibration.node.glif.io

```
Here are the fields in the `.env` file:

- `NODE_NAME` is the name of the node. 
- `NODE_DESCRIPTION` is the description of the node.
- `NODE_TYPE` is the type of the node.
- `MODE` can be `standalone` or `cluster`. If `standalone`, the node will run as a single node. If `cluster`, the node will run as a cluster node.
  - `standalone` mode is primarily for those who want to run delta in an isolated environment. This mode will create a local database and a local static API key for all requests. To get a static key, run `https://auth.estuary.tech/register-new-token`. Copy the generated key and paste it in the `DELTA_AUTH` field.
  - `cluster` mode is for those who want to run delta as a cluster. This mode will connect to a remote database. In this mode, you don't need to specify the `DELTA_AUTH` field. The node will use the API key provided by `Estuary`.
- `DB_DSN` is the database connection string. If `standalone`, the node will create a local database. If `cluster`, the node will connect to a remote database.
- `DELTA_AUTH` is the API key used to authenticate requests to the node. 
- `MAX_CLEANUP_WORKERS` is the maximum number of workers that can be used to clean up the blockstore. This is an optional field. If not specified, the default value is `1500`.

Put the `.env` file in the same location as the binary/executable.

# Build and run
## Using `make` lang
```
make all
./delta daemon --network=test
```

## Run the node with a custom blockstore location
```
./delta daemon --repo=/path/to/blockstore --network=test
``` 

## Run the node with a custom wallet
A delta node being a deal-making service, needs a wallet to make deals. You can specify a custom wallet location using the `--wallet-dir` flag. 
Note that this is a directory and not a file. The wallet file(s) is expected to be in this directory. 
```
./delta daemon --wallet-dir=/path/to/wallet --network=test
```
Note: You can register a new wallet later using the `/admin/wallet/register` endpoint.

# Build and run using docker
You can also build and run the delta node using docker. To do so, you need to have docker installed on your machine.
```
docker build -t delta .
docker run -it --rm -p 1414:1414 delta --repo=/path/to/blockstore --wallet-dir=<walletdir> --network=test
```

# Running the node the first time will do the following:
- Create a new wallet if one doesn't exist.
- Create a new database if one doesn't exist.
- Start in `cluster` mode if mode is not specified.
- Start the API server.
- Connect to the telemetry server.
- Compute the host operating system specifications and store them in the database.


# Console output 
The console output will look something like this:
```
OS: darwin
Architecture: arm64
Hostname: Alvins-MBP
Starting Delta daemon...
Setting up the whypfs node... 
repo:  .whypfs
walletDir:  ./wallet_estuary
mode:  cluster
enableWebsocket:  false
statsCollection:  true
network:  test
lotus api:  http://api.calibration.node.glif.io
ip:  142.189.91.167
Wallet address is:  t1aj6k36cyhndhscw7yy67f5ttbezzpsl6l7zu6iy
2023/06/13 13:17:36 INF    2 (145.40.77.207:4150) connecting to nsqd
Setting up the whypfs node... DONE
Computing the OS resources to use
Total memory: 108663640 bytes
Total system memory: 256330040 bytes
Total heap memory: 208633856 bytes
Heap in use: 177881088 bytes
Stack in use: 26247168 bytes
Total number of CPUs: 10
Number of CPUs that this Delta will use: 10
Note: Delta instance proactively recalculate resources to use based on the current load.
Computing the OS resources to use... DONE
Running pre-start clean up
Number of rows cleaned up: 0
Running pre-start clean up... DONE
Subscribing the event listeners
Subscribing to transfer channel states...
Subscribing to transfer channel events...
Subscribing the event listeners... DONE
Running the atomic cron jobs
Scheduling dispatchers and scanners...
Running the atomic cron jobs... DONE
Starting Delta.


     %%%%%%%%/          %%%%%%%%%%%%%%% %%%%%     %%%%%%%%%%%%%%%%%     %%%%%%  
    @@@@@@@@@@@@@@@     @@@@@@@@@@@@@@ @@@@@      @@@@@@@@@@@@@@@@@   @@@@@@@@  
    @@@@@     @@@@@@@  @@@@@@          @@@@@           @@@@@         @@@@@@@@@@ 
   @@@@@@       @@@@@  @@@@@          @@@@@            @@@@@       @@@@@  @@@@@ 
   @@@@@        @@@@@ @@@@@@@@@@@@@@ (@@@@@           @@@@@       @@@@@   @@@@@ 
  @@@@@@       @@@@@@ @@@@@@@@@@@@@  @@@@@           /@@@@@      @@@@@    #@@@@,
  @@@@@       @@@@@@ @@@@@*         @@@@@@           @@@@@     @@@@@@@@@@@@@@@@@
 @@@@@@@@@@@@@@@@@   @@@@@@@@@@@@@@ @@@@@@@@@@@@@@  @@@@@@    @@@@@        @@@@@
 @@@@@@@@@@@@@@     @@@@@@@@@@@@@@ @@@@@@@@@@@@@@@  @@@@@    @@@@@         @@@@@

(ᵔᴥᵔ)(ᵔᴥᵔ)(ᵔᴥᵔ)(ᵔᴥᵔ)(ᵔᴥᵔ)(ᵔᴥᵔ)(ᵔᴥᵔ)(ᵔᴥᵔ)(ᵔᴥᵔ)(ᵔᴥᵔ)(ᵔᴥᵔ)(ᵔᴥᵔ)(ᵔᴥᵔ)(ᵔᴥᵔ)(ᵔᴥᵔ)(ᵔᴥᵔ)

By: Protocol Labs - Outercore Engineering
version: v1.0.6-37-g9524cb4-dirty
Reporting Delta startup logs
2023/06/13 13:17:36 INF    1 (145.40.77.207:4150) connecting to nsqd
Reporting Delta startup logs... DONE
----------------------------------
Welcome! Delta daemon is running...
You can check the documentation at: https://github.com/application-research/delta/tree/main/docs
----------------------------------

   ____    __
  / __/___/ /  ___
 / _// __/ _ \/ _ \
/___/\__/_//_/\___/ v4.9.0
High performance, minimalist Go web framework
https://echo.labstack.com
____________________________________O/_______
                                    O\
⇨ http server started on [::]:1414
```

## Get Fil+ using Calibnet faucet
In order to make verified storage deals, we need to add FIL+ to the wallet.
- First, take note of the Wallet address above (`t1aj6k36cyhndhscw7yy67f5ttbezzpsl6l7zu6iy`)
- Go to [https://faucet.calibration.fildev.network/](https://faucet.calibration.fildev.network/) 
- Click on "Grant Datacap"
- Paste the wallet address in the input field and click on "Grant Datacap"
- It'll take a few minutes for the FIL+ to show up in the wallet. You can check the balance using the `http://localhost:1414/open/info/wallet/balance/<address>` endpoint.
- You'll see the verified balance (`verified_client_balance`) set to a non-zero value. 

# Next
Now that you have the node running, you can start making deals. 

- [Make an e2e deal](make-e2e-deal.md)
- [Make an import deal](make-import-deal.md)

If you to use an existing live `Delta` node, go to [getting started using a live delta node](getting-started-use-delta.md).
