# Make a batch import / offline deal.
Delta is a deal-making service that enables users to make deals with Storage Providers. It is an application that allows users to upload files to the Filecoin network and get them stored by Storage Providers.

In this section, we will walk you through the steps to use a Delta node to make deals.

# Make sure you have a `Delta` node.
If you are looking for a running Delta node, you can use [`node.delta.store`](https://node.delta.store/open/node/info).

If you want to stand up your own node, you can follow the instructions in [this](./getting-started-run-delta.md) document.

# Prepare the deal `metadata` request.
In order to create a successful deal, Delta requires the following information `metadata` request:
- `API_KEY`. This is used to authenticate the request. This is attached to the Authentication Header of the request.
- The content to be stored or the piece-commitment of the content.
    - `data file or cid`: The content to be stored. This can be a file or a directory.
    - `piece-commitment`: The piece-commitment of the content. This is the pre-computed piece cid, piece size (padded and unpadded) and file size.
- The `miner` to store the content.
  - If you don't have a miner, you can use the following:
    - run the SP miner selection cli tool `./delta sp selection --size-in-bytes=34359738368`. Get the `address` field from the response.
    - run the SP miner selection api `curl --location --request GET 'https://simple-sp-selection.onrender.com/api/providers?size_bytes=34359738368'`, get the `address` field from the response.
    - check out `https://data.storage.market/api/providers` to get a list of miners.
    - you can also omit this field and let the daemon select a miner for you.
- The `auto_retry` flag. This is a boolean flag that indicates whether the deal should be retried if it fails. The default value is `false`.
- The `deal_verify_state` to use to make the deal. The default value for an import deal is `verified`.
  - `verified` state is used to make deals with miners that support the `Verified` deal state.
  - `unverified` state is used to make deals with miners that support the `Unverified` deal state.
- The `unverified_max_price` to use to make the deal. The default value for an import deal is `0`.
  - `unverified_max_price` is used to make deals with miners that support the `unverified` deal state. This is the maximum price that the miner can charge for the deal.
  - If the `deal_verify_state` is `verified`, this field is ignored.
  - If the `deal_verify_state` is `unverified`, this field is used to make the deal.
- The `transfer_parameters` is an optional field that can be used to pull the data from remote URL source. 
  - `url`: The url of the content to be stored. This is the information required to transfer the content to the miner. This is an optional field. If this field is not provided, the content is assumed to be available external manual import process.

Here's the complete structure of the `metadata` request.
```
{
    "cid": "bafybeidty2dovweduzsne3kkeeg3tllvxd6nc2ifh6ztexvy4krc5pe7om",
    "miner":"f01963614",
    "wallet": {
        "address":"f1mmb3lx7lnzkwsvhridvpugnuzo4mq2xjmawvnfi"
    },
    "transfer_parameters": {
        "url":"https://ipfs.io/ipfs/bafybeicgdjdvwes3e5aaicqljrlv6hpdfsducknrjvsq66d4gsvepolk6y"
    },
    "piece_commitment": {
        "piece_cid": "baga6ea4seaqhfvwbdypebhffobtxjyp4gunwgwy2ydanlvbe6uizm5hlccxqmeq",
        "padded_piece_size": 4294967296
    },
    "size": 2500366291,
    "remove_unsealed_copy":true, 
    "skip_ipni_announce": true
}
```

# Make an import / offline deal.
Create a deal for a pre-computed piece-commitment by sending a `POST` request to the `/api/v1/deal/batch/imports` endpoint. The `metadata` request is the information required to make the deal.
## Request
```
curl --location --request POST 'http://localhost:1414/api/v1/deal/batch/imports' \
--header 'Authorization: Bearer [API_KEY]' \
--header 'Content-Type: application/json' \
--data-raw '[{
    "cid": "bafybeidty2dovweduzsne3kkeeg3tllvxd6nc2ifh6ztexvy4krc5pe7om",
    "miner":"f01963614",
    "piece_commitment": {
        "piece_cid": "baga6ea4seaqhfvwbdypebhffobtxjyp4gunwgwy2ydanlvbe6uizm5hlccxqmeq",
        "padded_piece_size": 4294967296
    },
    "size": 2500366291,
    "remove_unsealed_copy":true, 
    "auto_retry": false,
    "skip_ipni_announce": true
}]'
```

## Response
The response will look like this:
```
{
    "status": "success",
    "message": "Batch import request received. Please take note of the batch_import_id. You can use the batch_import_id to check the status of the deal.",
    "batch_import_id": 19
}
```
Take note of the `batch_import_id` field. This is the id of the batch import request. You can use this id to check the status of the deals made for this batch.

## Check the status of the batch
```
curl --location 'http://localhost:1414/open/stats/batch/imports/<batch_import_id>' \
--header 'Authorization: Bearer [API_KEY]'
[
    {
        "content": {
            "ID": 4607,
            "name": "bafybeidylyizmuhqny6dj5vblzokmrmgyq5tocssps3nw3g22dnlty7bhy",
            "size": 18010019221,
            "cid": "bafybeidylyizmuhqny6dj5vblzokmrmgyq5tocssps3nw3g22dnlty7bhy",
            "piece_commitment_id": 4540,
            "status": "making-deal-proposal",
            "request_type": "",
            "connection_mode": "import",
            "auto_retry": false,
            "last_message": "",
            "created_at": "2023-05-03T15:28:03.580059-04:00",
            "updated_at": "2023-05-03T15:28:06.744265-04:00"
        },
        "deal_proposal_parameters": [
            {
                "ID": 4564,
                "content": 4607,
                "label": "bafybeidylyizmuhqny6dj5vblzokmrmgyq5tocssps3nw3g22dnlty7bhy",
                "duration": 1480320,
                "start_epoch": 2845680,
                "end_epoch": 4326000,
                "transfer_params": "{}",
                "remove_unsealed_copy": false,
                "skip_ipni_announce": true,
                "verified_deal": true,
                "unverified_deal_max_price": "0",
                "created_at": "2023-05-03T15:28:03.857084-04:00",
                "updated_at": "2023-05-03T15:28:03.857085-04:00"
            }
        ],
        "deal_proposals": null,
        "deals": null,
        "piece_commitments": [
            {
                "ID": 4540,
                "cid": "bafybeidylyizmuhqny6dj5vblzokmrmgyq5tocssps3nw3g22dnlty7bhy",
                "piece": "baga6ea4seaqblmkqfesvijszk34r3j6oairnl4fhi2ehamt7f3knn3gwkyylmlq",
                "size": 18010019221,
                "padded_piece_size": 34359738368,
                "unnpadded_piece_size": 0,
                "status": "committed",
                "last_message": "",
                "created_at": "2023-05-03T15:28:03.493456-04:00",
                "updated_at": "2023-05-03T15:28:03.493457-04:00"
            }
        ]
    },
    {
        "content": {
            "ID": 4608,
            "name": "bafybeidylyizmuhqny6dj5vblzokmrmgyq5tocssps3nw3g22dnlty7bhy",
            "size": 18010019221,
            "cid": "bafybeidylyizmuhqny6dj5vblzokmrmgyq5tocssps3nw3g22dnlty7bhy",
            "piece_commitment_id": 4541,
            "status": "deal-proposal-failed",
            "request_type": "",
            "connection_mode": "import",
            "auto_retry": false,
            "last_message": "deal proposal rejected: deal proposal is identical to deal ea8bec53-45a1-4485-bbe9-a9136f217e96 (proposed at 2023-05-03 12:28:08.892250163 -0700 -0700)",
            "created_at": "2023-05-03T15:28:04.039661-04:00",
            "updated_at": "2023-05-03T15:28:10.707114-04:00"
        },
        "deal_proposal_parameters": [
            {
                "ID": 4565,
                "content": 4608,
                "label": "bafybeidylyizmuhqny6dj5vblzokmrmgyq5tocssps3nw3g22dnlty7bhy",
                "duration": 1480320,
                "start_epoch": 2845680,
                "end_epoch": 4326000,
                "transfer_params": "{}",
                "remove_unsealed_copy": false,
                "skip_ipni_announce": true,
                "verified_deal": true,
                "unverified_deal_max_price": "0",
                "created_at": "2023-05-03T15:28:04.33316-04:00",
                "updated_at": "2023-05-03T15:28:04.33316-04:00"
            }
        ],
        "deal_proposals": [
            {
                "ID": 169,
                "content": 4608,
                "unsigned": "",
                "signed": "bafyreigqq6k6bql5swasyqryqbcuncjcrdxj5qb4zkkmx46y4lqg3lu3oy",
                "meta": "bafyreigqq6k6bql5swasyqryqbcuncjcrdxj5qb4zkkmx46y4lqg3lu3oy",
                "created_at": "2023-05-03T15:28:08.803378-04:00",
                "updated_at": "2023-05-03T15:28:08.803378-04:00"
            }
        ],
        "deals": [
            {
                "ID": 177,
                "content": 4608,
                "propCid": "bafyreigqq6k6bql5swasyqryqbcuncjcrdxj5qb4zkkmx46y4lqg3lu3oy",
                "dealUuid": "31f4e925-7d58-4b2c-b0f9-7fcea91ee42a",
                "miner": "f01963614",
                "dealId": 0,
                "failed": true,
                "verified": true,
                "slashed": false,
                "failedAt": "0001-01-01T00:00:00Z",
                "dtChan": "",
                "transferStarted": "0001-01-01T00:00:00Z",
                "transferFinished": "0001-01-01T00:00:00Z",
                "onChainAt": "0001-01-01T00:00:00Z",
                "sealedAt": "0001-01-01T00:00:00Z",
                "lastMessage": "deal proposal rejected: deal proposal is identical to deal ea8bec53-45a1-4485-bbe9-a9136f217e96 (proposed at 2023-05-03 12:28:08.892250163 -0700 -0700)",
                "deal_protocol_version": "/fil/storage/mk/1.2.0",
                "created_at": "2023-05-03T15:28:08.061493-04:00",
                "updated_at": "2023-05-03T15:28:09.448522-04:00"
            }
        ],
        "piece_commitments": [
            {
                "ID": 4541,
                "cid": "bafybeidylyizmuhqny6dj5vblzokmrmgyq5tocssps3nw3g22dnlty7bhy",
                "piece": "baga6ea4seaqblmkqfesvijszk34r3j6oairnl4fhi2ehamt7f3knn3gwkyylmlq",
                "size": 18010019221,
                "padded_piece_size": 34359738368,
                "unnpadded_piece_size": 0,
                "status": "committed",
                "last_message": "",
                "created_at": "2023-05-03T15:28:03.948518-04:00",
                "updated_at": "2023-05-03T15:28:03.948518-04:00"
            }
        ]
    },
    {
        "content": {
            "ID": 4609,
            "name": "bafybeidylyizmuhqny6dj5vblzokmrmgyq5tocssps3nw3g22dnlty7bhy",
            "size": 18010019221,
            "cid": "bafybeidylyizmuhqny6dj5vblzokmrmgyq5tocssps3nw3g22dnlty7bhy",
            "piece_commitment_id": 4542,
            "status": "deal-proposal-failed",
            "request_type": "",
            "connection_mode": "import",
            "auto_retry": false,
            "last_message": "deal proposal rejected: deal proposal is identical to deal ea8bec53-45a1-4485-bbe9-a9136f217e96 (proposed at 2023-05-03 12:28:08.892250163 -0700 -0700)",
            "created_at": "2023-05-03T15:28:04.503183-04:00",
            "updated_at": "2023-05-03T15:28:09.730607-04:00"
        },
        "deal_proposal_parameters": [
            {
                "ID": 4566,
                "content": 4609,
                "label": "bafybeidylyizmuhqny6dj5vblzokmrmgyq5tocssps3nw3g22dnlty7bhy",
                "duration": 1480320,
                "start_epoch": 2845680,
                "end_epoch": 4326000,
                "transfer_params": "{}",
                "remove_unsealed_copy": false,
                "skip_ipni_announce": true,
                "verified_deal": true,
                "unverified_deal_max_price": "0",
                "created_at": "2023-05-03T15:28:04.764453-04:00",
                "updated_at": "2023-05-03T15:28:04.764454-04:00"
            }
        ],
        "deal_proposals": [
            {
                "ID": 170,
                "content": 4609,
                "unsigned": "",
                "signed": "bafyreigqq6k6bql5swasyqryqbcuncjcrdxj5qb4zkkmx46y4lqg3lu3oy",
                "meta": "bafyreigqq6k6bql5swasyqryqbcuncjcrdxj5qb4zkkmx46y4lqg3lu3oy",
                "created_at": "2023-05-03T15:28:09.121541-04:00",
                "updated_at": "2023-05-03T15:28:09.121541-04:00"
            }
        ],
        "deals": [
            {
                "ID": 178,
                "content": 4609,
                "propCid": "bafyreigqq6k6bql5swasyqryqbcuncjcrdxj5qb4zkkmx46y4lqg3lu3oy",
                "dealUuid": "218a512b-55cc-40ac-a297-9dcfc6976414",
                "miner": "f01963614",
                "dealId": 0,
                "failed": true,
                "verified": true,
                "slashed": false,
                "failedAt": "0001-01-01T00:00:00Z",
                "dtChan": "",
                "transferStarted": "0001-01-01T00:00:00Z",
                "transferFinished": "0001-01-01T00:00:00Z",
                "onChainAt": "0001-01-01T00:00:00Z",
                "sealedAt": "0001-01-01T00:00:00Z",
                "lastMessage": "deal proposal rejected: deal proposal is identical to deal ea8bec53-45a1-4485-bbe9-a9136f217e96 (proposed at 2023-05-03 12:28:08.892250163 -0700 -0700)",
                "deal_protocol_version": "/fil/storage/mk/1.2.0",
                "created_at": "2023-05-03T15:28:08.061475-04:00",
                "updated_at": "2023-05-03T15:28:09.642771-04:00"
            }
        ],
        "piece_commitments": [
            {
                "ID": 4542,
                "cid": "bafybeidylyizmuhqny6dj5vblzokmrmgyq5tocssps3nw3g22dnlty7bhy",
                "piece": "baga6ea4seaqblmkqfesvijszk34r3j6oairnl4fhi2ehamt7f3knn3gwkyylmlq",
                "size": 18010019221,
                "padded_piece_size": 34359738368,
                "unnpadded_piece_size": 0,
                "status": "committed",
                "last_message": "",
                "created_at": "2023-05-03T15:28:04.417978-04:00",
                "updated_at": "2023-05-03T15:28:04.417978-04:00"
            }
        ]
    }
]
```
# Next
Now that we can make an import deal, we can move on to the next step
- [Make an e2e deal](make-e2e-deal.md)
- [Check the status of your deal](content-deal-status.md)
- [Learn how to repair a deal](repair-retry.md)