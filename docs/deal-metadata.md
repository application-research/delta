# Metadata Request

The metadata request is used to create a deal. It contains the information required to create a deal.

# Properties of a metadata request
## cid
The `cid` field is the cid of the content to be stored. This is only required if the `connection_mode` is `import`. If the `connection_mode` is `e2e`, then the `cid` field is not required.
# miner
The `miner` field is the miner to store the content. This is required.
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
# connection_mode
The `connection_mode` field is the connection mode to use to make the deal. This is either `e2e` or `import`. This is required.
# size
The `size` field is the size of the content to be stored. This is only required if the `connection_mode` is `import`. If the `connection_mode` is `e2e`, then the `size` field is not required.
# remove_unsealed_copies
The `remove_unsealed_copies` field is a boolean field that indicates whether to remove unsealed copies of the content after the deal is made. This is optional. 
# skip_ipni_announce
The `skip_ipni_announce` field is a boolean field that indicates whether to skip announcing the deal to interplanetary indexer. This is optional. 
# duration_in_days or duration
- The `duration_in_days` field is the duration of the deal in days. This is optional. 
- The `duration` field is the duration of the deal in epochs. This is optional. 
# start_epoch_at_days or start_epoch
- The `start_epoch_at_days` field is the epoch to start the deal. This is optional.
- The `start_epoch` field is the epoch to start the deal. This is optional.
# deal_verify_state
The `deal_verify_state` field is the state of the deal verification. This is to indicate if the deal is from verified FIL or not. This is optional.
# label
The `label` field is a label for the deal. It has a limit of less than 100 characters. This is optional.

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
    },
    "connection_mode": "import",
    "size": 2500366291,
    "remove_unsealed_copies":true, 
    "skip_ipni_announce": true,
    "duration_in_days": 540, 
    // OR "duration": "1555200" // duration in epochs (30 seconds)
    "start_epoch_at_days": 14, // days to delay before the deal starts
    // OR "start_epoch": 2729999 // epoch at which to start the deal - see https://www.epochclock.io/
}
```

# Next
Now that you know how to create a metadata request, you can:
- [Make an import deal](./make-import-deal.md)
- [Check the status of your deal](content-deal-status.md)
- [Learn how to repair a deal](repair.md)
