# Running Containerized Delta

## Install and run `Delta` using `docker`
Make sure you have docker installed on your machine.

### Build a delta image with a specific wallet address
If you already have a wallet with datacap, you can pass it to the command below. This copies over the wallet directory to the containerized app and sets it as the default wallet.
```
make docker-compose-build WALLET_DIR=<walletdir> DESCRIPTION=MY_OWN_DELTA_WITH_A_SPECIFIC_WALLET TAG
```
You can then run the containerized app using the command below
```
make docker-compose-up
```

### Build and Run the current delta clone using docker-compose
Alternatively, you can build and run the current delta clone using docker-compose.
```
make docker-compose-run WALLET_DIR=<walletdir> DESCRIPTION=MY_OWN_DELTA_WITH_A_SPECIFIC_WALLET
```
