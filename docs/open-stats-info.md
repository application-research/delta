# Open statistics and information

There are several statistics and information endpoints that are open to the public. These endpoints are useful for monitoring the health of the delta service.

Open endpoints doesn't require any authentication (API_KEY).

## Information
### Get node information
To get `Delta` node information, we can use the `/open/node/info` endpoint.
```
curl --location 'http://localhost:1414/open/node/info'
```
### Get node connected peers
To get `Delta` node connected peers information, we can use the `/open/node/peers` endpoint.
```
curl --location 'http://localhost:1414/open/node/peers'
```
### Get node uuids
To get `Delta` node uuids information, we can use the `/open/node/uuids` endpoint.
```
curl --location 'http://localhost:1414/open/node/uuids'
```
Note: UUIDs are delta node identifiers. We use this to identify a delta node in the network.

## Stats
### Get totals info
To get `Delta` node totals information, we can use the `/open/stats/totals/info` endpoint. This includes total e2e and import deals made by the delta node.
```
curl --location 'http://localhost:1414/open/stats/totals/info'
```
### Get deal by cid
To get `Delta` node deal information by cid, we can use the `/open/stats/deals/by-cid/:cid` endpoint.
```
curl --location 'http://localhost:1414/open/stats/deals/by-cid/<cid>'
```
### Get deal by uuid
To get `Delta` node deal information by deal uuid, we can use the `/open/stats/deals/by-uuid/:uuid` endpoint.
```
curl --location 'http://localhost:1414/open/stats/deals/by-uuid/<uuid>'
```
### Get deal by deal id
To get `Delta` node deal information by deal id, we can use the `/open/stats/deals/by-dealid/:dealid` endpoint.
```
curl --location 'http://localhost:1414/open/stats/deals/by-dealid/<dealid>'
```
