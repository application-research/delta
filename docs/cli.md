# Delta CLI tools

Delta CLI tools are used to interact with the Delta node. It provides a set of commands to manage the node, and interact with the Delta API.
## Installation

Delta CLI is packaged with the Delta node. Build and install `Delta` node as described in [getting started running a delta node](getting-started-run-delta.md).

- Car file generation cli
- Piece commitment computation cli
- Storage deal making cli
- Content status check cli
- Wallet cli

## Usage
### Car file generation cli
#### Running `delta car` on a file
This will generate a car file for the file.
```bash
mkdir -p output
./delta car --source=<file> --output-dir=output
```

#### Running `delta car` on a directory
This will generate a car file for each file in the directory.
```bash
mkdir -p output
./delta car --source=<dir> --output-dir=output
```

#### Running `delta car` on a directory with a split size
This will generate a car file for each file in the directory. The car file will be split into chunks of the specified size.
```bash
mkdir -p output
./delta car --source=<dir> --output-dir=output --split-size=1024
```

#### Running `delta car` on a directory and generate piece_commitment for each file
This will generate a car file for each file in the directory. The car file will be split into chunks of the specified size.
```bash
mkdir -p output
./delta car --source=<dir> --output-dir=output --split-size=1024 --include-commp=true
```

This will generate collection of car files with split-size of 1024 bytes for each file in the directory.

**Output along with the CAR files generated on the `output-dir` folder**
```
[
    {
        "payload_cid": "bafybeibgymnw3smet7qzkxdbt6o33rxdzgyz7u5e4tebomgy72tigdb22y",
        "commp": "baga6ea4seaqd5q5y7fgdv3gsot4pxzoojov7cwobm3yy54ygy2n3msksatkaepy",
        "padded_size": 2048,
        "size": 1024,
        "cid_map": {
            "": {
                "is_dir": true,
                "cid": "bafybeibgymnw3smet7qzkxdbt6o33rxdzgyz7u5e4tebomgy72tigdb22y"
            },
            "output": {
                "is_dir": true,
                "cid": "bafybeid3wn5zdwg6focmh2mmul3shq75t6iakwogfx5ur36gxsixa5wugi"
            },
            "output/README.md": {
                "is_dir": true,
                "cid": "bafybeiey4qtix73tuka6s5fjflub63dv4opjs4wgylxgs5xxwcewosg44q"
            },
            "output/README.md/README.md_0000": {
                "is_dir": false,
                "cid": "bafkreibgngawuoylmk4lq2gvltd2c2kksfhpwok7smjtnrsfxgm2efbwqm"
            }
        }
    },
    {
        "payload_cid": "bafybeibt7rhs7qke6xyw4h4jrrjbastlcense443oyh73rhpokkkk4zqy4",
        "commp": "baga6ea4seaqevyznwx7sqrhwuvooftdqaaonlq6ovjtgvescpkixvcav7fzwogi",
        "padded_size": 4096,
        "size": 1024,
        "cid_map": {
            "": {
                "is_dir": true,
                "cid": "bafybeibt7rhs7qke6xyw4h4jrrjbastlcense443oyh73rhpokkkk4zqy4"
            },
            "output": {
                "is_dir": true,
                "cid": "bafybeifmzm2n2stcttfaskpmckyke4gdall7sr7o6npq2cln3l5cmctjha"
            },
            "output/README.md": {
                "is_dir": true,
                "cid": "bafybeidndf3333a7ytd4vk6phawfv7vrmycobl3pjd4hpsdh2yuraran6u"
            },
            "output/README.md/README.md_0000": {
                "is_dir": false,
                "cid": "bafkreibgngawuoylmk4lq2gvltd2c2kksfhpwok7smjtnrsfxgm2efbwqm"
            },
            "output/README.md/README.md_0001": {
                "is_dir": false,
                "cid": "bafkreie5uynugmo6okjcyk4tpmpjacn37rwnkeofqc7qrdfkgdef4m56ni"
            }
        }
    },
]
```


### Piece Commitment computation cli
#### Running `delta commp` on a file

Get the piece commitment of a file
```
./delta commp --file=README.md --include-payload-cid=false --mode=fast 
```
**Output**
```
{
    "file_name": "README.md",
    "size": 4317,
    "cid": "bafybeifjbv6uenj75owkbdhz7aiaribxmgpufelpaeznp2m74x65b5soxq",
    "piece_commitment": {
        "piece_cid": "baga6ea4seaqniq6k5yh3iiguuj4rhsz235n7wnbqrkqkjrdzwbihsctujtlpgcq",
        "padded_piece_size": 8192,
        "unpadded_piece_size": 8128
    }
}

```

#### Running `delta commp` on a directory
Get the piece commitment of all the files in a directory
```
./delta commp --dir=<dir> --include-payload-cid=false --mode=fast 
```

**Output**
```
[
    {
        "file_name": "docs/open-stats-info.md",
        "size": 544,
        "cid": "bafybeif7ztnhq65lumvvtr4ekcwd2ifwgm3awq4zfr3srh462rwyinlb4y",
        "piece_commitment": {
            "piece_cid": "baga6ea4seaqecjhgezzr4arsswzyvx5weqv6twh6bth6tuv7ftip5bmjg55q6aa",
            "padded_piece_size": 1024,
            "unpadded_piece_size": 1016
        }
    },
    {
        "file_name": "docs/repair.md",
        "size": 302,
        "cid": "bafybeif7ztnhq65lumvvtr4ekcwd2ifwgm3awq4zfr3srh462rwyinlb4y",
        "piece_commitment": {
            "piece_cid": "baga6ea4seaqmriuea5uxxnq6oqv7akpng5asuiq7f6wvbijr5u24mieob6yzknq",
            "padded_piece_size": 512,
            "unpadded_piece_size": 508
        }
    },
    {
        "file_name": "docs/running-delta-docker.md",
        "size": 1476,
        "cid": "bafybeif7ztnhq65lumvvtr4ekcwd2ifwgm3awq4zfr3srh462rwyinlb4y",
        "piece_commitment": {
            "piece_cid": "baga6ea4seaqddn2heeq3dqytsxmhl6sejddqbhijedmys35yuzjpdifb2rvtqga",
            "padded_piece_size": 2048,
            "unpadded_piece_size": 2032
        }
    },
    {
        "file_name": "docs/storage-deals.md",
        "cid": "bafybeif7ztnhq65lumvvtr4ekcwd2ifwgm3awq4zfr3srh462rwyinlb4y",
        "piece_commitment": {
            "piece_cid": "baga6ea4seaqomqafu276g53zko4k23xzh4h4uecjwicbmvhsuqi7o4bhthhm4aq"
        }
    }
]
```
### Storage deal making cli
The storage deal making cli needs a running Delta daemon to work. 

Run delta node first with a wallet address that has datacap. You can use the following command to run the delta node with a wallet address that has datacap.
```
./delta daemon --wallet-dir=/path/to/wallet
```
Note that you can register multiple wallets with Delta. Once the daemon is up, you can register a wallet using the instructions at [registering a wallet](manage-wallets.md)
You also need an API_KEY to run. Use the following command to get an API_KEY.
```
curl --location --request GET 'https://auth.estuary.tech/register-new-token'

{
"expires": "2123-02-03T21:12:15.632368998Z",
"token": "<API_KEY>"
}
```

#### Make an E2E deal
```
./delta deal --api-key=<API_KEY> --type=e2e --file=<> --metadata='{"miner":"f01963614","connection_mode":"e2e", "skip_ipni_announce":false}'
```

Making a deal with a registered wallet
```
./delta deal --api-key=<API_KEY> --type=e2e --file=<> --metadata='{"miner":"f01963614","connection_mode":"e2e", "skip_ipni_announce":false,"wallet":{"address":<address>}}'
```

**Output**

Take note of the content_id. You'll use this to get the status of the deal.
```
{
    "status": "success",
    "message": "File uploaded and pinned successfully",
    "content_id": 12,
    "deal_request_meta": {
        "cid": "bafybeifjbv6uenj75owkbdhz7aiaribxmgpufelpaeznp2m74x65b5soxq",
        "connection_mode": "e2e",
        "miner": "f01963614",
        "piece_commitment": {},
        "wallet": {}
    },
    "deal_proposal_parameter_request_meta": {
        "ID": 12,
        "content": 12,
        "created_at": "2023-03-25T03:07:56.251200708Z",
        "duration": 1494720,
        "label": "bafybeifjbv6uenj75owkbdhz7aiaribxmgpufelpaeznp2m74x65b5soxq",
        "updated_at": "2023-03-25T03:07:56.251208042Z"
    }
}

```

#### Make an Import deal
```
./delta deal --api-key=<API_KEY> --type=import --metadata='[
    {
        "cid": "bafybeidylyizmuhqny6dj5vblzokmrmgyq5tocssps3nw3g22dnlty7bhy",
        "wallet": {
            "address": "f1mmb3lx7lnzkwsvhridvpugnuzo4mq2xjmawvnfi",
            "dataset_name": "",
            "balance": {
                "balance_filecoin": 0,
                "balance_datacap": 0
            }
        },
        "miner": "f01963614",
        "piece_commitment": {
            "piece_cid": "baga6ea4seaqblmkqfesvijszk34r3j6oairnl4fhi2ehamt7f3knn3gwkyylmlq",
            "padded_piece_size": 34359738368
        },
        "connection_mode": "import",
        "size": 18010019221,
        "remove_unsealed_copies": true,
        "skip_ipni_announce": true,
        "duration_in_days": 537,
        "start_epoch_at_days": 3
    }
]'
```
**Output**

Take note of all the content_ids. 
```
[
    {
        "status": "success",
        "message": "Request received",
        "content_id": 10,
        "deal_request_meta": {
            "cid": "bafybeidylyizmuhqny6dj5vblzokmrmgyq5tocssps3nw3g22dnlty7bhy",
            "connection_mode": "import",
            "duration_in_days": 537,
            "miner": "f01963614",
            "piece_commitment": {
                "padded_piece_size": 34359738368,
                "piece_cid": "baga6ea4seaqblmkqfesvijszk34r3j6oairnl4fhi2ehamt7f3knn3gwkyylmlq"
            },
            "remove_unsealed_copies": true,
            "size": 18010019221,
            "skip_ipni_announce": true,
            "start_epoch": 2742000,
            "start_epoch_at_days": 3,
            "wallet": {
                "address": "f1mmb3lx7lnzkwsvhridvpugnuzo4mq2xjmawvnfi"
            }
        },
        "deal_proposal_parameter_request_meta": {
            "ID": 10,
            "content": 10,
            "created_at": "2023-03-25T03:06:20.937783303Z",
            "duration": 1526400,
            "label": "bafybeidylyizmuhqny6dj5vblzokmrmgyq5tocssps3nw3g22dnlty7bhy",
            "remove_unsealed_copy": true,
            "updated_at": "2023-03-25T03:06:20.937783428Z"
        }
    }
]

```


### Status of a content 
Once you get a deal request made, you can get the status of a content.

The status is based on the lifecycle of a content. View the lifecycle of a content [here](content-deal-status.md).

You can get the status of a content by running the following command.
```
./delta status --type=content --id=<content_id>
```

**Output**

```
{
    "content": {
        "ID": 4,
        "name": "README.md",
        "size": 4317,
        "cid": "bafybeifjbv6uenj75owkbdhz7aiaribxmgpufelpaeznp2m74x65b5soxq",
        "piece_commitment_id": 4,
        "status": "deal-proposal-failed",
        "request_type": "",
        "connection_mode": "e2e",
        "last_message": "connecting to f01963614: miner connection failed: failed to dial 12D3KooWRFzN8SoRVayNw3ho8PArwVpZDsRzcQG5W4DguE4euTS9:\n  * [/ip4/208.185.75.116/tcp/11003] dial tcp4 0.0.0.0:6745-\u003e208.185.75.116:11003: i/o timeout",
        "created_at": "2023-03-25T02:35:30.031246585Z",
        "updated_at": "2023-03-25T02:35:36.288990422Z"
    },
    "deal_proposal_parameters": [
        {
            "ID": 4,
            "content": 4,
            "label": "bafybeifjbv6uenj75owkbdhz7aiaribxmgpufelpaeznp2m74x65b5soxq",
            "duration": 1494720,
            "created_at": "2023-03-25T02:35:30.399267585Z",
            "updated_at": "2023-03-25T02:35:30.399267752Z"
        }
    ],
    "deal_proposals": null,
    "deals": null,
    "piece_commitments": [
        {
            "ID": 4,
            "cid": "bafybeifjbv6uenj75owkbdhz7aiaribxmgpufelpaeznp2m74x65b5soxq",
            "piece": "baga6ea4seaqdyszr56jr3lf7acxj22rj52h2dnunnxjuhn7rtaval7f2oo662bi",
            "size": 4425,
            "padded_piece_size": 8192,
            "unnpadded_piece_size": 8128,
            "status": "open",
            "last_message": "",
            "created_at": "2023-03-25T02:35:30.618410752Z",
            "updated_at": "2023-03-25T02:35:30.618410836Z"
        }
    ]
}

```

### Wallet CLI

#### Register a wallet
To register a wallet, you need to export the wallet from lotus/boostd and use the hex value to register the wallet.
```
./delta wallet register --delta-host="https://cake.delta.store" --hex="<HEX FROM LOTUS/BOOSTD EXPORT" --api-key=<API_KEY>
{
    "message": "Successfully imported a wallet address. Please take note of the following information.",
    "wallet_addr": "f1mmb3lx7lnzkwsvhridvpugnuzo4mq2xjmawvnfi",
    "wallet_uuid": "18a21dd5-cb7c-11ed-a090-3cecef773e44"
}

```
#### List all wallets
To list all wallets associated to an API KEY, run the following command.
```
./delta wallet list --delta-host="https://cake.delta.store" --api-key=<API_KEY>
{
    "wallets": [
        {
            "ID": 1,
            "uuid": "84a029e2-c93c-11ed-98be-3cecef773e44",
            "addr": "f1mmb3lx7lnzkwsvhridvpugnuzo4mq2xjmawvnfi",
            "owner": <API_KEY>,
            "key_type": "secp256k1",
            "private_key": "<KEY>",
            "created_at": "2023-03-23T05:35:27.866873092Z",
            "updated_at": "2023-03-23T05:35:27.866873092Z"
        },
        {
            "ID": 3,
            "uuid": "ad7c03ad-cb77-11ed-a090-3cecef773e44",
            "addr": "f1mmb3lx7lnzkwsvhridvpugnuzo4mq2xjmawvnfi",
            "owner": <API_KEY>,
            "key_type": "secp256k1",
            "private_key": "<KEY>",
            "created_at": "2023-03-26T01:43:59.049178076Z",
            "updated_at": "2023-03-26T01:43:59.049178076Z"
        },
    ]
} 
```

