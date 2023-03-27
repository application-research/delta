# Getting Deal Status Information

Delta is a deal-making service that enables users to make deals with Storage Providers. It is an application that allows users to upload files to the Filecoin network and get them stored by Storage Providers.

In this section, we will walk you through the steps to use a Delta node to get the status of a deal.

![image](https://user-images.githubusercontent.com/4479171/226716467-e8f7b2b5-d5e8-4d19-8b53-31889422ea5c.png)

# Content Lifecycle
When `Delta` accepts a deal request, it will go through the following steps:
1. `Delta` will pin the content to the light node.
2. `Delta` will create a content record in its database.
3. `Delta` will create a miner assignment to a content record in its database.
4. `Delta` will create a wallet assignment to a content record in its database.
5. `Delta` will create a piece commitment record in its database.
6. `Delta` will compute the piece commitment of the content.
7. `Delta` will create a deal proposal parameters record in its database.
8. `Delta` will create a deal proposal record in its database.
9. `Delta` will send the deal proposal to the Storage Provider.
10. If it's an `import` deal, there will be no data transfer and the final status will be `deal-proposal-sent`.
11. If it's an `e2e` deal, `Delta` will send the data to the Storage Provider.
    1. `Delta` will wait for the Storage Provider to finish the deal.
    2. `Delta` will update the deal record in its database.

Over the course of the deal, `Delta` will update the status of the deal in its database. We can use this information to get the status of the deal.

## In-Progress / Successful status
The status of the content will be based on the following:
- `pinned` - The content is pinned to the light node.
- `piece-computing` - The content is being computed into a piece.
- `piece-computed` - The content has been computed into a piece.
- `piece-assigned` - The piece has been assigned to a deal.
- `making-deal-proposal` - The deal proposal is being made.
- `sending-deal-proposal` - The deal proposal is being sent to the Storage Provider.
- `deal-proposal-sent` - The deal proposal has been sent to the Storage Provider.
- `transfer-started` - The data transfer has started.
- `transfer-finished` - The data transfer has finished.

## Completed status
- `deal-proposal-sent` - The deal proposal has been sent to the Storage Provider.
- `transfer-finished` - The data transfer has finished.

## Failed status
during the deal-making process, contents can run into some problems and `Delta` marks as failed. The status of the content will be based on the following:
- `failed-to-pin` - The content failed to be pinned to the light node.
- `failed-to-process` - The content failed to be processed.
- `piece-compute-failed` - The content failed to be computed into a piece.
- `deal-propose-failed` - The deal proposal failed to be made.
- `transfer-failed` - The data transfer failed.

# Get the status of the deal.
Whenever a deal is made, Delta stores the deal information in its database. A content_id is generated for each request. We can use this content_id to get the status of the deal.

To get the status of the deal, we can use the `/api/v1/stats/content/:content_id` or `/open/stats/content/:content_id` endpoint.
## Request
```
curl --location --request GET 'http://localhost:1414/api/v1/stats/content/:content_id' \
--header 'Authorization: Bearer [API_KEY]'
```

Alternatively, you can view the status of the deal using the `/open/stats/content` endpoint.
Note: This endpoint does not require an API key.
```
curl --location --request GET 'http://localhost:1414/open/stats/content/:content_id'
```

Alternatively, you can view the status of the deal using the `/open/stats/content` endpoint.
Note: This endpoint does not require an API key.
```
curl --location --request GET 'http://localhost:1414/open/stats/content/:content_id'
```

## Response
```
{
    "content": {
        "ID": 941,
        "name": "random_1679373273451488879.dat",
        "size": 8000000000,
        "cid": "bafybeidwffy4qs36ybibpzixfm3ut5hcyv2ijwmo7y6voumu4ncsom2t3q",
        "piece_commitment_id": 941,
        "status": "transfer-finished",
        "request_type": "",
        "connection_mode": "e2e",
        "last_message": "transfer-finished",
        "created_at": "2023-03-21T04:36:37.445518451Z",
        "updated_at": "2023-03-21T05:37:33.67016768Z"
    },
    "deal_proposal_parameters": [
        {
            "ID": 941,
            "content": 941,
            "label": "bafybeidwffy4qs36ybibpzixfm3ut5hcyv2ijwmo7y6voumu4ncsom2t3q",
            "duration": 1494720,
            "created_at": "2023-03-21T04:36:37.568568002Z",
            "updated_at": "2023-03-21T04:36:37.568568167Z"
        }
    ],
    "deal_proposals": [
        {
            "ID": 941,
            "content": 941,
            "unsigned": "",
            "signed": "bafyreihh56dud7j2537eoab2kfbjhyiu7cdwva2cxzxntg3kcqxlqv45fu",
            "meta": "bafyreihh56dud7j2537eoab2kfbjhyiu7cdwva2cxzxntg3kcqxlqv45fu",
            "created_at": "2023-03-21T04:57:46.922432926Z",
            "updated_at": "2023-03-21T04:57:46.922433016Z"
        }
    ],
    "deals": [
        {
            "ID": 941,
            "content": 941,
            "propCid": "bafyreihh56dud7j2537eoab2kfbjhyiu7cdwva2cxzxntg3kcqxlqv45fu",
            "dealUuid": "d5289c3b-bf83-4e06-b684-82eb92930256",
            "miner": "f01929568",
            "dealId": 941,
            "failed": false,
            "verified": true,
            "slashed": false,
            "failedAt": "0001-01-01T00:00:00Z",
            "dtChan": "12D3KooWNBXquoZV7SKzVRrNRPUc77CoB5AvTsByQdfZbavtt5Va-12D3KooWCqDPmD8PDJjtjo2j5rDzyJRkaNspwufhaKMBQRXkhWDS-1679372417035255337",
            "transferStarted": "2023-03-21T04:58:22.718340135Z",
            "transferFinished": "2023-03-21T05:37:32.112988888Z",
            "onChainAt": "2023-03-21T05:37:32.112989437Z",
            "sealedAt": "2023-03-21T05:37:32.112989109Z",
            "lastMessage": "transfer-finished",
            "deal_protocol_version": "/fil/storage/mk/1.2.0",
            "created_at": "2023-03-21T04:57:29.030339541Z",
            "updated_at": "2023-03-21T05:37:32.113238357Z"
        }
    ],
    "piece_commitments": [
        {
            "ID": 941,
            "cid": "bafybeidwffy4qs36ybibpzixfm3ut5hcyv2ijwmo7y6voumu4ncsom2t3q",
            "piece": "baga6ea4seaqezfdhxc4qvhu3geeornt2mwykzoxqtc6fxherhu3xrr7o3vfmkey",
            "size": 8000790398,
            "padded_piece_size": 8589934592,
            "unnpadded_piece_size": 8522825728,
            "status": "committed",
            "last_message": "",
            "created_at": "2023-03-21T04:57:07.261890437Z",
            "updated_at": "2023-03-21T04:58:02.625657927Z"
        }
    ]
}
```
