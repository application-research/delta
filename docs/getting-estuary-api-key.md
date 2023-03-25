# Getting API Token/Key

Delta requires an API key to make deals. To make deals, you need to get an API token. You can get an API token from [here](https://estuary.tech/).

Alternatively, you can also get an API token by running the following request:

## Request
```
curl --location --request GET 'https://auth.estuary.tech/register-new-token'
```

## Response
```
{
"expires": "2123-02-03T21:12:15.632368998Z",
"token": "<API_KEY>"
}
```