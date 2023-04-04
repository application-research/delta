# Repair and Retry storage deals

Delta has built-in repair and retry functionality. This is useful for when a storage deal fails for some reason. The repair and retry functionality is built into the daemon and can be accessed thru HTTP API.

Definition of terms:
- retry - retry is a process of dispatching the same content to the deal processor. The content will go thru the same process ie, piece commitment computation to deal making for e2e deal and preparing the deal proposal to sending the deal proposal for import deals.
- repair - repair is a process of redefining some of the metadata of a deal and re-dispatching the deal to the deal processor. 

## Manual Repair and Retry
Users can also manually repair and retry deals via HTTP API.

### Retry a deal for an e2e content
```
curl --location --request GET 'http://localhost:1414/api/v1/retry/deal/endo:contentId' \
--header 'Authorization: Bearer [API_KEY]' \
--header 'Content-Type: application/json' 
```

### To retry an e2e deal
```
curl --location --request GET 'http://localhost:1414/api/v1/retry/deal/end-to-end/:contentId' \
--header 'Authorization: Bearer [API_KEY]' \
--header 'Content-Type: application/json' 
```

### To retry an import deal
```
curl --location --request GET 'http://localhost:1414/api/v1/retry/deal/import/:contentId' \
--header 'Authorization: Bearer [API_KEY]' \
--header 'Content-Type: application/json' 
```


## Auto Retry
Users who wants to retry deals can set this up via metadata request.

The `metadata` parameter is a JSON object that can contain the following fields:

- `auto_retry` - boolean value that indicates whether the daemon should automatically retry a storage deal if it fails. Default is `false`.
  When an auto_retry field is set to true, the deal will run retries on several miners using https://sp-select.delta.store/api/providers until a miner accepts the deal.

- if miner is set and the auto-retry is set to true, delta will use the given miner. If the miner rejects or faults the deal, Delta will re-try the deal with other miners.
- if the miner is not set and the auto-retry is set to true, delta will look into https://sp-select.delta.store/api/providers to check miners who can accept the deal. If the miner rejects or faults the deal, Delta will re-try the deal with other miners.
- if the miner is set but auto-retry is set to false, delta will use the given miner. It will not retry the deal if the miner rejects the proposal.
- if the miner is not set and auto-retry is set to false, delta will look into https://sp-select.delta.store/api/providers to check miners who can accept the deal. It will not retry the deal if the miner rejects the proposal.

The [status check](content-deal-status.md) API will return the list of deals that have been retried.
