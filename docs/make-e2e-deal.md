# Make an e2e / online deal.

Create an online deal for a content by sending a `POST` request to the `/api/v1/deal/end-to-end` endpoint. The `data` request is the content to be stored. The `metadata` request is the information required to make the deal.

# Make sure you have a `Delta` node.
- If you are looking for a running Delta node, you can use [node.delta.store](https://node.delta.store).
- If you want to stand up your own node, you can follow the instructions in [this](./getting-started-run-delta.md) document.

# Prepare the deal `metadata` request.
In order to create a successful deal, Delta requires the following information `metadata` request:
- `API_KEY`. This is used to authenticate the request. This is attached to the Authentication Header of the request.
- The content to be stored or the piece-commitment of the content.
    - `data file or cid`: The content to be stored. This can be a file or a directory.
    - `piece-commitment`: The piece-commitment of the content. This is the pre-computed piece cid, piece size (padded and unpadded) and file size.
- The `miner` to store the content.
- The connection mode to use to make the deal. This is either `e2e` or `import`.
    - `e2e` mode (online deal) is used to make deals with miners that support the `e2e` connection mode.
    - `import` mode (offline) is used to make deals with miners that support the `import` connection mode.

Here's the complete structure of the `metadata` request.
```
{
    "cid": "bafybeidty2dovweduzsne3kkeeg3tllvxd6nc2ifh6ztexvy4krc5pe7om",
    "miner":"f01963614",
    "wallet": {
        "address":"f1mmb3lx7lnzkwsvhridvpugnuzo4mq2xjmawvnfi"
    },
    "piece_commitment": {
        "piece_cid": "baga6ea4seaqhfvwbdypebhffobtxjyp4gunwgwy2ydanlvbe6uizm5hlccxqmeq",
        "padded_piece_size": 4294967296
    },
    "connection_mode": "import",
    "size": 2500366291,
    "remove_unsealed_copies":true, 
    "skip_ipni_announce": true
}
```

## Request
```
curl --location --request POST 'http://localhost:1414/api/v1/deal/end-to-end' \
--header 'Authorization: Bearer [API_KEY]' \
--form 'data=@"my-file"' \
--form 'metadata="{\"miner\":\"f01963614\",\"connection_mode\":\"e2e\"}"'
```

A few things to note here:
- The `data` request is the content to be stored. This can be a file or a directory.
- The `metadata` request is the information required to make the deal. This is a JSON object. At minimum, it should contain the `miner` and `connection_mode` fields. The `miner` field is the miner to store the content. The `connection_mode` field is the connection mode to use to make the deal. This is either `e2e` or `import`.
- If no `wallet` is specified, Delta will use the default wallet that it generated when it was started.

## Response
The response will look like this:
```
{
    "status": "success",
    "message": "Deal request received. Please take note of the content_id. You can use the content_id to check the status of the deal.",
    "content_id": 1,
    "deal_request_meta": {
        "cid": "bafybeib6l6odanq5zrspbw4c7fys4jspshgwzuuhotnpljsivhdythw6xu",
        "miner": "f02031042",
        "wallet": {
            "address": "f1mmb3lx7lnzkwsvhridvpugnuzo4mq2xjmawvnfi"
        },
        "piece_commitment": {},
        "connection_mode": "e2e"
    }
}
```
Take note of the `content_id` field. This is the id of the content that was uploaded. This is used to get the status of the deal.

# Get the status of the deal.
To get the status of the deal, we can use the `/api/v1/stats/content/:content_id` endpoint.
## Request
```
curl --location --request GET 'http://localhost:1414/api/v1/stats/content/:content_id' \
--header 'Authorization: Bearer [API_KEY]'
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

# Using a specific wallet
The wallet is one of the most important aspect of making a filecoin and the wallet holds the FIL and DataCap that's going to be used to transaction with the network.

Registering a wallet to a Delta node means that the wallet owner is TRUSTING the delta node to hold it's keys. This is a very important step and should be done with care.

## Register a wallet
To register a wallet to a live Delta node, we can use the `/admin/wallet/register-hex` endpoint. This endpoint is only available on the admin port.

### Request
```
curl --location --request POST 'http://localhost:1414/admin/wallet/register-hex' \
--header 'Authorization: Bearer [API_KEY]' \
--header 'Content-Type: application/json' \
--data-raw '{"hex_key":"<HEX FROM LOTUS / BOOSTD WALLET EXPORT>"}'
```


### Response
The response will look like this:
```
{
    "message": "Successfully imported a wallet address. Please take note of the following information.",
    "wallet_addr": "f1mmb3lx7lnzkwsvhridvpugnuzo4mq2xjmawvnfi",
    "wallet_uuid": "4d4589d0-c7a2-11ed-b245-9e0bf0c70138"
}
```

We can now use the `wallet_addr` value to make a deal.

## Use wallet to prepare the `metadata` request
Once a wallet is registered, we can add a `wallet` field to the `metadata` request to make a deal using that wallet.
```
curl --location --request POST 'http://localhost:1414/api/v1/deal/end-to-end' \
--header 'Authorization: Bearer [API_KEY]' \
--form 'data=@"my-file"' \
--form 'metadata="{\"miner\":\"f02031042\",\"connection_mode\":\"e2e\", \"wallet\":{\"address\":\"f1mmb3lx7lnzkwsvhridvpugnuzo4mq2xjmawvnfi\"}}"'
```

# Next
Now that we can make an e2e deal, let's look at how to make an import deal.
- [Make an import deal](./make-import-deal.md)
- [Check the status of your deal](content-deal-status.md)
- [Learn how to repair a deal](repair.md)
