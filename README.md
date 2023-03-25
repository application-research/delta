# Δ Delta
Generic DealMaking MicroService using whypfs + filclient + estuary_auth

![image](https://user-images.githubusercontent.com/4479171/218267752-9a7af133-4e36-4f4c-95da-16b3c7bd73ae.png)

## Features
- Make e2e / online and import / offline storage deals.
- Compute piece_commitments using variety of methods
  - boost piece commitment computation
  - parallel piece commitment computation
  - filclient piece commitment computation
- Assign deals to specific miners
- Assign deals to specific wallets
- Shows all the deals made for specific user
- Extensive status information about the deals made
- Support for multiple wallets
- Lightweight. Uses a light ipfs node `whypfs-core` to stage the files.
- Cleans up itself after a deal is made.
- Monitors deal progress
- Monitors transfer progress
- Containerized deployment

## Getting Started
To get started, go to our docs [here](docs).
# Author
Protocol Labs Outercore Engineering.