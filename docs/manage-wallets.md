# Managing Filecoin Wallets in Delta
Filecoin wallets are essential for making deals with miners. Delta can manage multiple Filecoin wallets. This section will walk you through the steps to manage Filecoin wallets in Delta.

## Register a Filecoin wallet
In order to make deals with miners, you need to register a Filecoin wallet with Delta. You can register a Filecoin wallet by sending a `POST` request to the `/api/v1/wallet/register` endpoint. The `metadata` request is the information required to register the wallet.

To register a wallet to a live Delta node, we can use the `/admin/wallet/register-hex` endpoint. This endpoint is only available on the admin port.
## Request
```
curl --location --request POST 'http://localhost:1414/admin/wallet/register-hex' \
--header 'Authorization: Bearer [API_KEY]' \
--header 'Content-Type: application/json' \
--data-raw '{"hex_key":"<HEX FROM LOTUS / BOOSTD WALLET EXPORT>"}'
```

## Response
The response will contain the `wallet_addr` and `wallet_uuid` of the registered wallet.
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


## List all registered wallets
### Request
```
curl --location --request GET 'http://localhost:1414/admin/wallet/list' \
--header 'Authorization: Bearer [API_KEY]' \
```
### Response
```
{
    "wallets": [
        {
            "ID": 1,
            "uuid": "4d4589d0-c7a2-11ed-b245-9e0bf0c70138",
            "addr": "f1mmb3lx7lnzkwsvhridvpugnuzo4mq2xjmawvnfi",
            "owner": "ESTc904e6ee-8dfe-44b8-864f-37280e1117f9ARY",
            "key_type": "secp256k1",
            "private_key": "<REDACATED>",
            "created_at": "2023-03-21T00:39:01.339102-04:00",
            "updated_at": "2023-03-21T00:39:01.339102-04:00"
        }
    ]
}
```
