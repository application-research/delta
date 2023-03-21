# Delta metrics collection

Delta collects metric on all nodes. It does this in two ways:
- open telemetry api using opencensus
- a message queue reporting system using [delta-events-consumer](https://github.com/application-research/delta-events-consumer)

## How does it work?
Every request that goes thru Delta API is logged. This logged message is then sent to a `nsq` where a background consumer 
process collects it and persist is on a timescaleDB.


![image](https://user-images.githubusercontent.com/4479171/226726850-59828c4a-dba8-4082-877a-12efd9474641.png)
