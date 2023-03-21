# Getting Delta Node Information

Delta is a deal-making service that enables users to make deals with Storage Providers. It is an application that allows users to upload files to the Filecoin network and get them stored by Storage Providers.

In this section, we will walk you through the steps to use a Delta node to get the status of a deal.

# Node Information
To get `Delta` node information, we can use the `/open/node/info` endpoint.
```
curl --location --request GET 'http://localhost:1414/open/node/info'
```

# Connected Peers 
To get `Delta` node connected peers information, we can use the `/open/node/peers` endpoint.
```
curl --location --request GET 'http://localhost:1414/open/node/peers'
```

# Node Uuid Information 
To get `Delta` node uuids information, we can use the `/open/node/uuids` endpoint.
```
curl --location --request GET 'http://localhost:1414/open/node/uuids'
```