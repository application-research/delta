# Δ Delta
Generic DealMaking MicroService using whypfs + filclient + estuary_auth

![image](https://user-images.githubusercontent.com/4479171/218267752-9a7af133-4e36-4f4c-95da-16b3c7bd73ae.png)

## Features
- Make e2e / online and import / offline storage deals.
- Compute piece_commitments using variety of methods
  - boost commp
  - parallel commp
  - filclient commp
- Assign deals to specific miners
- Assign deals to specific wallets
- Shows all the deals made for specific user
- Extensive status information about the deals made
- Support for multiple wallets
- Lightweight. Uses a light ipfs node `whypfs-core` to stage the files.
- Cleans up itself after a deal is made.
- Monitors transfer progress

This is strictly a deal-making service. It'll pin the files/CIDs but it won't keep the files. Once a deal is made, the CID will be removed from the blockstore. For retrievals, use the retrieval-deal microservice.

To Get Started, check out the [docs](https://delta.estuary.tech/docs/overview)