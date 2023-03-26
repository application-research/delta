# Open statistics and information

There are several statistics and information endpoints that are open to the public. These endpoints are useful for monitoring the health of the delta service.

Open endpoints doesn't require any authentication (API_KEY).

## Information
### Get node information
```
curl --location 'http://localhost:1414/open/node/info'
```
### Get node connected peers
```
curl --location 'http://localhost:1414/open/node/peers'
```
### Get node uuids
```
curl --location 'http://localhost:1414/open/node/uuids'
```

## Stats
### Get totals info
```
curl --location 'http://localhost:1414/open/stats/totals/info'
```

### Get deal by cid
```
curl --location 'http://localhost:1414/open/stats/deals/by-cid/<cid>'
```
### Get deal by uuid
```
curl --location 'http://localhost:1414/open/stats/deals/by-uuid/<uuid>'
```
### Get deal by deal id
```
curl --location 'http://localhost:1414/open/stats/deals/by-dealid/<dealid>'
```


