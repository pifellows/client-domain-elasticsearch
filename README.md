# client-domain-elasticsearch
A simple clientAPI and DomainService in order to facilitate indexing a large JSON "map" file into elasticsearch

The system comprises of 3 containers and a 4th container (kibana) for development purposes. 

## ClientAPI
The client API hosts a simple REST API that allows the user to initiate loading the data from a file as well as allowing the user to retrieve a document from it's id. It listens on port 5000 by default.

### Routes
*/getport* returns a JSON response and takes a single query string parameter of *portid* which returns the corresponding document to the user. If the document does not exist, a blank record is returned. If an unexpected error occurs then that error is also returned to the user.

*/start* commands the API to begin processing the JSON file specified at startup. Ideally, this endpoint would be replaced by one that would take an update file and begin processing.

### Commandline arguements
The ClientAPI has the following commandline arguements:
- *--domainservice* - This is used to specify the connection string that the ClientAPI uses to connect to the PortsDomainService using gRPC. This requires the hostname and the port number e.g. localhost:5001
- *--file* - This is the filepath of the JSON file to be parsed and indexed once the *start* endpoint is called

## PortsDomainService
The PortsDomainService is responsible for indexing the documents that it has been sent over gRPC and index them into Elasticsearch. When indexing documents, they are batched together in groups of 100 in order to lower the number of HTTP requests made to Elasticsearch.
In addition, the Service also allows a caller to retrieve a document from Elasticsearch using it's document ID.

### Commandline arguements
The PortsDomainService has the following commandline arguements:
- *--es* - This is used to specify the connection string of the Elasticsearch instance to be used for document storage. This requires the protocol, hostname and port number. For example: http://localhost:9200
- *--indexname* - This is used to specify the name of the index that documents should be stored in and retrieved from

The batch size that is used for indexing documents is not exposed to the user

## Elasticsearch
Elasticsearch acts as the datastore in this system. Whilst it is not necessarily the best choice for this scenario, it is a store that I am familiar working with. Elasticsearch servces user traffic on port 9200

## Kibana
Kibana is a browser based app that can be used to interact with Elasticsearch instances. Whilst it is not a part of the system, it can be used to examine any index on the Elasticsearch instance.

### Missing
The solution is missing several features that should be considered best practice. These include:
- sufficient logging
- unit tests for features, especially the JsonChunkReader which is isolated enough to test fully
- healthcheck endpoints so that they can be hosted in various environments
- docker compose does not work as expected. Elastcsearch takes a certain amount of time to load and as a result, the DomainService fails to connect due to there being no "retry" logic. Changing the elasticsearch client to a "SimpleClient" may improve this but it is currently untested.
