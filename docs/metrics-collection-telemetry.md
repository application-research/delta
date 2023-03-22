# Delta metrics collection

Delta collects metric on all nodes. It does this in two ways:
- open telemetry api using opencensus
- a message queue reporting system using [delta-events-consumer](https://github.com/application-research/delta-events-consumer)

## How does it work?
Every request that goes thru Delta API is logged. This logged message is then sent to a `nsq` where a background consumer 
process collects it and persist is on a timescaleDB/Postgres.

### Delta Events Consumer and OTEL API

Delta uses [delta-events-consumer](https://github.com/application-research/delta-events-consumer/blob/main/README.md) to collect metrics. It is a background process that listens to the `nsq` and persist the data on a timescaleDB/Postgres.
![image](https://user-images.githubusercontent.com/4479171/226726850-59828c4a-dba8-4082-877a-12efd9474641.png)

Delta also uses OpenTelemetry API to collect metrics. This is done by using the [opencensus](https://opencensus.io/) library.

## Data collected
- `api` - all api calls without the request body
- `api_error` - all api calls that return an error
- `content` - all content related events
- `deal` - all deal related events
- `miner` - all miner related events
- `wallet` - all wallet related events without the wallet seed
- `wallet_error` - all wallet related events that return an error
- `deal proposal` - all deal proposal related events
- `deal proposal error` - all deal proposal related events that return an error
- `deal proposal parameters` - all deal proposal parameters related events
- `deal proposal parameters error` - all deal proposal parameters related events that return an error
- `piece commitment` - all piece commitment related events
- `piece commitment error` - all piece commitment related events that return an error
- `miner assignment` - all miner assignment related events
- `miner assignment error` - all miner assignment related events that return an error
- `wallet assignment` - all wallet assignment related events
- `wallet assignment error` - all wallet assignment related events that return an error
- `transfer` - all transfer related events
- `transfer error` - all transfer related events that return an error
- `status` - all status related events

Note: We don't collect any request body or wallet seed. The collected data is only used for debugging and monitoring purposes only to improve the product.  
Grafana dashboard can be found [here](https://protocollabs.grafana.net/d/xCXVv8-4k/global-delta-dashboard?orgId=1&refresh=10s)