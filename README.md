# Description

The goal is to build a simple crypto wallet screening HTTP service. For how-to run please see [README.md](./cmd/serverd/README.md).

## Service providers

### Blockmate

Blockmate is used as risk data provider in the application. In order to successfully run `wallet-scanner` token must be provided. Follow steps below in order to acquire a token:

* Project token for authorising requests must be created in [portal](portal.blockmate.io), after creating a new project. 
* To acquire JWT token please see [docs](https://docs.blockmate.io/reference/userapi-authenticateproject).

### Immudb

Immudb is used as a tamper-proof database to store history of address risk categories for audit history purposes.

## Functional description

Service at this point has two endpoints:

### POST /wallet/{address}/categories
Returns risk categories list for given address. Additionally, returned list of categories will be stored into immudb for audit history purposes. 

Accepts URL query parameter `address` which represents Ethereum network wallet.
Send a request to the running service instance ( presuming its running on port 80 ):

```bash
curl -X POST 'http://localhost/wallet/0xe9e9afac38e64728f1afbb2b65dec7be7c704c05/categories' -v
```

### GET /wallet/{address}/categories
Retrieves a list of historical risk categories for given address.

Accepts URL query parameter `address` which represents Ethereum network wallet.
Send a request to the running service instance ( presuming its running on port 80 ):

```bash
curl 'http://localhost/wallet/0xe9e9afac38e64728f1afbb2b65dec7be7c704c05/categories' -v
```

## Implementation rationale

Solution was implemented having following presumptions in mind:

* Project will grow in the future, - design architecture that can scale by creating right abstractions and separations of concerns.
* Having one day to the implement same solution would result in way much simpler application and structure.
* There are `TODO's` here and there and not everything was covered with tests. Idea was to demonstrate design and architecture allowing to scale and implement mentioned `TODO's` without refactoring much of the code resulting into complete solution.
* To have some fun, - just notice how neatly handlers layer is separated from business logic allowing us to test and stay focused only on what's relevant in particular layer! :)

Possible improvements:

* circuit breaker guarding our service in case of provider, network outage.
* logs redirect in test mode
* TODOs
* ...
