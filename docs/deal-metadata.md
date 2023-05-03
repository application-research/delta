# Metadata Request

The metadata request is used to create a deal. It contains the information required to create a deal.

# Properties of a metadata request
## cid
The `cid` field is the cid of the content to be stored. This is only required if the `connection_mode` is `import`. If the `connection_mode` is `e2e`, then the `cid` field is not required.
# miner
The `miner` field is the miner to store the content. This is required.
# source
This is the source mulitaddress on where to retrieve the given cid. This is only required when the "cid" field is given.
# unverified_deal_max_price
The `unverified_deal_max_price` field is the maximum price to pay for an unverified deal. This is required only if the `deal_verify_state` is "unverified". If no `unverified_deal_max_price` is specified, Delta will use the default value of 0.000000000000000000 FIL.
# wallet
The `wallet` field is the wallet to use to make the deal. This is optional. If no wallet is specified, Delta will use the default wallet that it generated when it was started.
## address
The `address` field is the address of the wallet to use to make the deal. This is only required if the `wallet` field is specified.
# piece_commitment
The `piece_commitment` field is the piece commitment of the content to be stored. This is only required if the `connection_mode` is `import`. If the `connection_mode` is `e2e`, then the `piece_commitment` field is not required.
## piece_cid 
The `piece_cid` field is the piece cid of the content to be stored. This is only required if the `connection_mode` is `import`. If the `connection_mode` is `e2e`, then the `piece_cid` field is not required.
## padded_piece_size
The `padded_piece_size` field is the padded piece size of the content to be stored. This is only required if the `connection_mode` is `import`. If the `connection_mode` is `e2e`, then the `padded_piece_size` field is not required.
## unpadded_piece_size 
The `unpadded_piece_size` field is the unpadded piece size of the content to be stored. This is only required if the `connection_mode` is `import`. If the `connection_mode` is `e2e`, then the `unpadded_piece_size` field is not required.
# transfer_parameters
The `transfer_parameters` field is the transfer parameters set when preparing a deal. It allows deal clients to pull data using the transfer parameter specified. This is only required if the `connection_mode` is `import`. If the `connection_mode` is `e2e`, then the `transfer_parameters` field is not required.
## url 
The `url` field is the url of the content to be stored. This will allow deal clients to pull data from a remote url source. This is only required if the `connection_mode` is `import`. If the `connection_mode` is `e2e`, then the `url` field is not required.
# connection_mode
The `connection_mode` field is the connection mode to use to make the deal. This is either `e2e` or `import`. This is required.
# size
The `size` field is the size of the content to be stored. This is only required if the `connection_mode` is `import`. If the `connection_mode` is `e2e`, then the `size` field is not required.
# remove_unsealed_copy
The `remove_unsealed_copy` field is a boolean field that indicates whether to remove unsealed copies of the content after the deal is made. This is optional. 
# skip_ipni_announce
The `skip_ipni_announce` field is a boolean field that indicates whether to skip announcing the deal to interplanetary indexer. This is optional. 
# duration_in_days
The `duration_in_days` field is the duration of the deal in days. This is optional.
# start_epoch_in_days
The `start_epoch_at_days` field is the epoch to start the deal. This is optional.
# deal_verify_state
The `deal_verify_state` field is the state of the deal verification. This is to indicate if the deal is from verified FIL or not. This is optional.
valid values are: `verified`, `unverified`. Default: `verified`.
# label
The `label` field is a label for the deal. It has a limit of less than 100 characters. This is optional.
# auto_retry
The `auto_retry` field is a boolean field that indicates whether to automatically retry the deal if it fails. This is optional.

When set to true, the deal will be retried if the failure falls under the "acceptable" failures. Note that we only have a list of acceptable failures at the moment.

The deal will be retried with a randomly selected miner based on the file and location of delta instance.

# Here's the complete structure of the `metadata` request.
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
        "unpadded_piece_size": 2500366291
    },
    "connection_mode": "import",
    "size": 2500366291,
    "label": "my deal",
    "deal_verify_state": "verified",
    "remove_unsealed_copy":true, 
    "skip_ipni_announce": true,
    "duration_in_days": 540, 
    "auto_retry": false,
    "start_epoch_in_days": 14, // days to delay before the deal starts
}
```

# Next
Now that you know how to create a metadata request, you can:
- [Make an import deal](./make-import-deal.md)
- [Check the status of your deal](content-deal-status.md)
- [Learn how to repair a deal](repair-retry.md)
