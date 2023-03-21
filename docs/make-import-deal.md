# Make an import / offline deal.
Delta is a deal-making service that enables users to make deals with Storage Providers. It is an application that allows users to upload files to the Filecoin network and get them stored by Storage Providers.

In this section, we will walk you through the steps to use a Delta node to make deals.

# Make sure you have a `Delta` node.
If you are looking for a running Delta node, you can use [node.delta.store](https://node.delta.store).

If you want to stand up your own node, you can follow the instructions in [this](./getting-started-run-delta.md) document.

# Prepare the deal `metadata` request.
In order to create a successful deal, Delta requires the following information `metadata` request:
- `Estuary API Key`. This is used to authenticate the request. This is attached to the Authentication Header of the request.
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

# Make an import / offline deal.
Create a deal for a pre-computed piece-commitment by sending a `POST` request to the `/api/v1/deal/piece-commitments` endpoint. The `metadata` request is the information required to make the deal.
## Request
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
}]'
```
A few things to note here:
- The `metadata` request is the information required to make the deal. This is a JSON array. Each element in the array is a JSON object. At minimum, it should contain the `miner`, `piece_commitment.piece_cid`,`piece_commitment.padded_piece_size`, `size` and `connection_mode` fields.
- The `miner` field is the miner that will store the content.
- The `piece_commitment` field is the piece-commitment of the content.
    - This is the pre-computed piece cid, piece size (padded and unpadded) and file size.
- The `connection_mode` field is the connection mode to use to make the deal. This is an `import`.

## Response
The response will look like this:
```
[
    {
        "status": "success",
        "message": "File uploaded and pinned successfully",
        "content_id": 1,
        "request_meta": {
            "cid": "bafybeidylyizmuhqny6dj5vblzokmrmgyq5tocssps3nw3g22dnlty7bhy",
            "miner": "f01963614",
            "duration_in_days": 537,
            "piece_commitment": {
                "piece_cid": "baga6ea4seaqblmkqfesvijszk34r3j6oairnl4fhi2ehamt7f3knn3gwkyylmlq",
                "padded_piece_size": 34359738368
            },
            "connection_mode": "import",
            "size": 18010019221,
            "start_epoch": 2730480,
            "start_epoch_at_days": 3,
            "remove_unsealed_copies": true,
            "skip_ipni_announce": true
        }
    }
]
```
Take note of the `content_id` field. This is the id of the content that was uploaded. This is used to get the status of the deal.

# Using a specific wallet
The wallet is one of the most important aspect of making a filecoin and the wallet holds the FIL and DataCap that's going to be used to transaction with the network.

Registering a wallet to a Delta node means that the wallet owner is TRUSTING the delta node to hold it's keys. This is a very important step and should be done with care.

## Register a wallet
To register a wallet to a live Delta node, we can use the `/admin/wallet/register-hex` endpoint. This endpoint is only available on the admin port.
### Request
```
curl --location --request POST 'http://localhost:1414/admin/wallet/register-hex' \
--header 'Authorization: Bearer [ESTUARY_API_KEY]' \
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

# Next
Now that we can make an import deal, we can move on to the next step
- [Make an e2e deal](make-e2e-deal.md)
- [Check the status of your deal](deal-status.md)
- [Learn how to repair a deal](repair.md)