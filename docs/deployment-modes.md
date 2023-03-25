# Deployment modes
By default, delta will run on cluster mode but users can standup a standalonemode. standalone mode primarily for those who want to run delta in an isolated environment. This mode will create a local database and a local static API key for all requests.

## Prepare the .env file.
Copy the `.env.example` file to `.env` and update the values as needed.
```
# Node info
NODE_NAME=delta-node
NODE_DESCRIPTION=Experimental Deal Maker
NODE_TYPE=delta-main

# Database configuration
MODE=standalone
DB_DSN=delta.db
DELTA_AUTH=[NODE_ESTUARY_API_KEY_HERE]

# Frequencies
MAX_CLEANUP_WORKERS=1500
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


## MODE=Standalone
With `.env` file in place, run the following command to start the node in standalone mode.
```
./delta daemon --mode=standalone
```

## MODE=Cluster
By default, delta will run on cluster mode. This mode will create a local database that can be reconfigured to use a remote HA database and estuary-auth as the authentication and authorizatio component.

When running in cluster mode, users need to register for an ESTUARY_API_KEY using the following command.

## Request
```
curl --location --request GET 'https://auth.estuary.tech/register-new-token'
```

## Response
```
{
    "expires": "2123-02-03T21:12:15.632368998Z",
    "token": "<ESTUARY_API_KEY>"
}
```

Place the ESTUARY_API_KEY in the .env file.

### Run in cluster mode
With `.env` file in place, run the following command to start the node in standalone mode.
```
./delta daemon --mode=cluster
```
