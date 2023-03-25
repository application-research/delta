# Running Delta using Docker

## Install and run `Delta` using `docker`
Make sure you have docker installed on your machine. You check the following link to install docker on your machine: https://docs.docker.com/engine/install/

Once you have docker installed, you're good to proceed to the next steps.

### Build a delta image with a specific wallet address
If you already have a wallet with datacap, you can pass it to the command below. This copies over the wallet directory to the containerized app and sets it as the default wallet.
```
make docker-compose-build WALLET_DIR=<walletdir> DESCRIPTION=MY_OWN_DELTA_WITH_A_SPECIFIC_WALLET TAG
```
You can then run the containerized app using the command below
```
make docker-compose-up
```

### Build and run the current delta clone using docker-compose
Alternatively, you can build and run the current delta clone using docker-compose.
```
make docker-compose-run WALLET_DIR=<walletdir> DESCRIPTION=MY_OWN_DELTA_WITH_A_SPECIFIC_WALLET
```

### Check localhost
You can check the localhost to see if the delta app is running
```
curl --location --request GET 'http://localhost:1414/open/node/info'
```

# Next
Now that you can access a live Delta node, you are now ready to make a deal. You can now go to the following guides:

- [Make an e2e deal](make-e2e-deal.md)
- [Make an import deal](make-import-deal.md)

If you to run your own `Delta` node, go to [getting started running a delta node](getting-started-run-delta.md).