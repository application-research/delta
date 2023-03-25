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
- To get started on running delta, go to the [getting started to run delta](getting-started-run-delta.md)
- To get started on using a live delta instance, go to the [getting started to use delta](getting-started-use-delta.md)
- To learn more about `delta cli` go to the [delta cli](cli.md)
- To learn more about running delta using docker, go to the [run delta using docker](running-delta-docker.md)
- To learn more about deployment modes, go to the [deployment modes](deployment-modes.md)
- To get estuary api key, go to the [estuary api keys](getting-estuary-api-key.md)
- To manage wallets, go to the [managing wallets](manage-wallets.md)
- To make an end-to-end deal, go to the [make e2e deals](make-e2e-deal.md)
- To make an import deal, go to the [make import deals](make-import-deal.md)
- To learn how to repair a deal, go to the [repairing and retrying deals](repair.md)
- To learn how to access the open statistics and information, go to the [open statistics and information](open-stats-info.md)
- To learn about the content lifecycle and check status of the deals, go to the [content lifecycle and deal status](content-deal-status.md)

# Author
Protocol Labs Outercore Engineering.