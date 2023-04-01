# Repair and Retry storage deals

Delta has built-in repair and retry functionality. This is useful for when a storage deal fails for some reason. The repair and retry functionality is built into the daemon and can be accessed thru HTTP API.

## Auto Retry
The daemon will automatically retry a storage deal if it fails

This can be set in the `metadata` request parameter when making a storage deal. The `metadata` parameter is a JSON object that can contain the following fields:

- `auto_retry` - boolean value that indicates whether the daemon should automatically retry a storage deal if it fails. Default is `false`.
  When an auto_retry field is set to true, the deal will run retries on several miners using https://sp-select.delta.store/api/providers until a miner accepts the deal.

- if miner is set and the auto-retry is set to true, delta will use the given miner. If the miner rejects or faults the deal, Delta will re-try the deal with other miners.
- if the miner is not set and the auto-retry is set to true, delta will look into https://sp-select.delta.store/api/providers to check miners who can accept the deal. If the miner rejects or faults the deal, Delta will re-try the deal with other miners.
- if the miner is set but auto-retry is set to false, delta will use the given miner. It will not retry the deal if the miner rejects the proposal.
- if the miner is not set and auto-retry is set to false, delta will look into https://sp-select.delta.store/api/providers to check miners who can accept the deal. It will not retry the deal if the miner rejects the proposal.

The [status check](content-deal-status.md) API will return the list of deals that have been retried.

## Manual Repair and Retry
### Retry a deal for a content
```
curl --location --request GET 'http://localhost:1414/api/v1/retry/deal/:contentId' \
--header 'Authorization: Bearer [API_KEY]' \
--header 'Content-Type: application/json' 
```

### Repair or retry a deal for a content
```
curl --location --request GET 'http://localhost:1414/api/v1/repair/deal/content/:contentId' \
--header 'Authorization: Bearer [API_KEY]' \
--header 'Content-Type: application/json' 
```

### Repair or retry a deal for a piece-commitment-id
```
curl --location --request GET 'http://localhost:1414/api/v1/repair/deal/piece-commitment-id/:pieceCommitmentId' \
--header --header 'Authorization: Bearer [API_KEY]' \
--header 'Content-Type: application/json' 
```

