# Δ Delta
Generic DealMaking MicroService using whypfs + filclient + estuary_auth

![image](https://user-images.githubusercontent.com/4479171/218267752-9a7af133-4e36-4f4c-95da-16b3c7bd73ae.png)


## Features
- Creates a deal for large files. The recommended size is 1GB. 
- Shows all the deals made for specific user

This is strictly a deal-making service. It'll pin the files/CIDs but it won't keep the files. Once a deal is made, the CID will be removed from the blockstore. For retrievals, use the retrieval-deal microservice.

## Process Flow
- client upload files or specifies a pre-computed piece_commitments
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
MODE=standalone
DB_DSN=stg-deal-maker
#REPO=/mnt/.whypfs # shared mounted repo

# Frequencies
MAX_CLEANUP_WORKERS=1500
```

Running this the first time will generate a wallet. Make sure to get FIL from the [faucet](https://verify.glif.io/) and fund the wallet




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

## REST API Endpoints

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
- miner is required
- connection_mode is optional. Default is `e2e` (formerly known as online deals)
- 
```
curl --location --request POST 'http://localhost:1414/api/v1/deal/content' \
--header 'Authorization: Bearer [ESTUARY_API_KEY]' \
--form 'data=@"/Users/alvinreyes/Downloads/baga6ea4seaqhfvwbdypebhffobtxjyp4gunwgwy2ydanlvbe6uizm5hlccxqmeq.car"' \
--form 'metadata="{\"miner\":\"f01963614\",\"connection_mode\":\"e2e\"}"'
```

### Import mode (formerly known as offline deal)
Use the following endpoint to upload a file with a specific miner, duration, piece size and connection mode.
```
curl --location --request POST 'http://localhost:1414/api/v1/deal/piece-commitment' \
--header 'Authorization: Bearer [ESTUARY_API_KEY]' \
--header 'Content-Type: application/json' \
--data-raw '{
    "cid": "bafybeidty2dovweduzsne3kkeeg3tllvxd6nc2ifh6ztexvy4krc5pe7om",
    "miner":"f01963614",
    "piece_commitment": {
        "piece_cid": "baga6ea4seaqhfvwbdypebhffobtxjyp4gunwgwy2ydanlvbe6uizm5hlccxqmeq",
        "padded_piece_size": 4294967296
    },
    "connection_mode": "import",
    "size": 2500366291,
    "remove_unsealed_copies":true, 
    "skip_ipni_announce": true
}'
```

### Batch Import mode (formerly known as offline deals)
The request body is an array of objects. The following is an example of a batch import request.
```
curl --location --request POST 'http://localhost:1414/api/v1/deal/piece-commitments' \
--header 'Authorization: Bearer [ESTUARY_API_KEY]' \
--header 'Content-Type: application/json' \
--data-raw '[{
    "cid": "bafybeidty2dovweduzsne3kkeeg3tllvxd6nc2ifh6ztexvy4krc5pe7om",
    "miner":"f01963614",
    "piece_commitment": {
        "piece_cid": "baga6ea4seaqhfvwbdypebhffobtxjyp4gunwgwy2ydanlvbe6uizm5hlccxqmeq",
        "padded_piece_size": 4294967296
    },
    "connection_mode": "import",
    "size": 2500366291,
    "remove_unsealed_copies":true, 
    "skip_ipni_announce": true
},
{
    "cid": "bafybeidty2dovweduzsne3kkeeg3tllvxd6nc2ifh6ztexvy4krc5pe7om",
    "miner":"f01963614",
    "piece_commitment": {
        "piece_cid": "baga6ea4seaqhfvwbdypebhffobtxjyp4gunwgwy2ydanlvbe6uizm5hlccxqmeq",
        "padded_piece_size": 4294967296
    },
    "connection_mode": "import",
    "size": 2500366291,
    "remove_unsealed_copies":true, 
    "skip_ipni_announce": true
},
{
    "cid": "bafybeidty2dovweduzsne3kkeeg3tllvxd6nc2ifh6ztexvy4krc5pe7om",
    "miner":"f01963614",
    "piece_commitment": {
        "piece_cid": "baga6ea4seaqhfvwbdypebhffobtxjyp4gunwgwy2ydanlvbe6uizm5hlccxqmeq",
        "padded_piece_size": 4294967296
    },
    "connection_mode": "import",
    "size": 2500366291,
    "remove_unsealed_copies":true, 
    "skip_ipni_announce": true
}]'
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


## CLI
### Get the commp of a file using commp cli
```
./delta commp --file=<>
```

### Get the commp of a CAR file using commp cli
```
./delta commp-car --file=<>
```

if you want to get the commp of a CAR file for offline deal, use the following command
```
./delta commp-car --file=<> --for-import
```
The output will be as follows
```
{
    "cid": "bafybeidty2dovweduzsne3kkeeg3tllvxd6nc2ifh6ztexvy4krc5pe7om",
    "wallet": {},
    "commp": {
        "piece": "baga6ea4seaqhfvwbdypebhffobtxjyp4gunwgwy2ydanlvbe6uizm5hlccxqmeq",
        "padded_piece_size": 4294967296
    },
    "connection_mode": "import",
    "size": 2500366291
}
```

### Get the commp of a CAR file using commp cli and pass to the delta api to make an offline deal
```
./delta commp-car --file=baga6ea4seaqhfvwbdypebhffobtxjyp4gunwgwy2ydanlvbe6uizm5hlccxqmeq.car --for-offline --delta-api-url=http://localhost:1414 --delta-api-key=[ESTUARY_API_KEY]
```

Output
```
{
   "status":"success",
   "message":"File uploaded and pinned successfully",
   "content_id":208,
   "piece_commitment_id":172,
   "meta":{
      "cid":"bafybeidty2dovweduzsne3kkeeg3tllvxd6nc2ifh6ztexvy4krc5pe7om",
      "wallet":{
         
      },
      "commp":{
         "piece":"baga6ea4seaqhfvwbdypebhffobtxjyp4gunwgwy2ydanlvbe6uizm5hlccxqmeq",
         "padded_piece_size":4294967296,
         "unpadded_piece_size":4261412864
      },
      "connection_mode":"import",
      "size":2500366291
   }
}
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

# Kubernetes
## Prerequisites
This repository provides all kubernetes deployment artifacts required for delta. It has been developed and tested in WSL, Git Bash (on Windows) and linux. All environments require the following:

- git
- make
- [docker](https://www.docker.com/products/docker-desktop)

## Development / Deployment Environments
There is a kubernetes environment that can be quickly run from WSL or Git Bash (and probably less quickly from powershell). It is a docker-wrapped kubernetes cluster using the popular Kind library. This configuration produces a nicely emulated Kind "cluster", using containers instead of physical nodes which is really handy for testing and developing different cluster node configurations without having to leave your local machine.

### Quick Start

Calling the `cluster.setup` target will install all of the necessary kubernetes tools for developing against a local development [Kind](https://kind.sigs.k8s.io/) cluster:

`make cluster.setup`

Instantiate the local Kind development cluster

`make cluster.up`

To uninstall all kubernetes manifests and create a fresh cluster (danger: you will lose all data)

`make cluster.down`

### Kubernetes YAML

There is also a directory with the raw kuberenetes yaml files for development purposes. It follows the same workflow as the helm quickstart except it uses the `k8s.` prefixed make targets:

`make k8s.up`

`make k8s.down`

Expose the delta service port address
`kubectl port-forward --address 0.0.0.0 svc/delta 1414:1414`

Optionally expose the Portainer management tool service port address
`kubectl port-forward --address 0.0.0.0 svc/delta 9000:9000`

## Kubernetes Installation

Install the persistent volume, deployment, service and hpa artifacts onto an external kubernetes cluster via:

`kubectl apply -f k8s/delta`

Expose the delta service port address
`kubectl port-forward --address 0.0.0.0 svc/delta 1414:1414`

## Make commands
For convienience, a lot of the commands to manage the deployment have been bundled in a project Makefile. For development, the most common targets will be:

- Start Cluster

  `make cluster.start`

- Generate TLS keys, deploy to kubernetes, wait for the pods/services/deployments to start and map the frontend to ports 80 and 443

  `make k8s.up`

- Delete the deployment from kubernetes

  `make k8s.down`

- Generate TLS keys, deploy to kubernetes, wait for the pods/services/deployments to start and map the frontend to ports 80 and 443

  `make k8s.all`

- Generate TLS keys, deploy to kubernetes

  `make k8s.install`

- Delete the deployment from kubernetes

  `make k8s.uninstall`

- Delete the generated keys and remove the deployment from kubernetes

  `make k8s.clean`

- Delete the generated keys

  `make k8s.clean-local`

- Run in a clean development container, mounting the current directory to /workspace

  `make k8s.dev`

- Open the k8s dashboard web app

  `make dash.k8s`

- Open the portainer dashboard web app

  `make k8s.portainer`

- Stop the local cluster instance

  `cluster.stop`

- Clean up the entire cluster from the system

  `make cluster.delete`

- generate tls keys and deploy them as a k8s tls secret for use in the nginx pod 

  `make generate.keys`

- Apply the kubernetes artifacts to the cluster

  `make k8s.deploy`

- Remove all components from kubernetes

  `make k8s.delete`

- Map the local port 8080 to cluster service port

  `make k8s.start`

- Start a local port-forward service mapping the k8s service to external port 8080

  `make k8s.startd`

- Stop the local service

  `k8s.stopd`

- Wait time for all services to start up

  `k8s.wait`

## Author
Outercore Engineering Team.
