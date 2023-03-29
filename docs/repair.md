# Repair and Retry storage deals

Delta has built-in repair and retry functionality. This is useful for when a storage deal fails for some reason. The repair and retry functionality is built into the daemon and can be accessed thru HTTP API.

## Repair/retry
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

