# ZD

This is the Zendesk Service that is responsible for producing User Events.

## Running ZD

### Service Dependencies

This service depends on 3 other's, one would be a service that provides user and event data as this is needed to generate random user events, we refer to it as the User Service. The second required service is an event queue, this is used as a way for consumers of this service to receive user events for further processing. The last service that is required is a caching service that will cache the data from the User 

When it comes to the service that provides ZD user and event data, we would use the service called [pd-users-api](https://github.com/TSE-Coders/pd-users-api), for the event queue we make use of [RabbitMQ](https://www.rabbitmq.com/) and for caching we are using [Redis](https://redis.io/).

### Single Host/Local Environment

To run this service you will need to have Golang version 1.21. When it comes to RabbitMQ and Redis, you can either install it directly or run it with Docker. For your convenience I created a compose file [docker/dependencies.yaml](./docker/dependencies.yaml) that can be used to spin up an instance of RabbitMQ, Redis and Redis Insight (Redis UI). To spin up these instances run the command below:

``` bash
docker compose -f docker/dependencies.yaml up -d
```

Once the RabbitMQ and Redis instances are running, start up the [pd-users-api](https://github.com/TSE-Coders/pd-users-api) service by following it's documentation. 

The environment variables that this service and the dependencies will use, can be found in [docker/env](./docker/env).

Now that the required services and the environment variables are set, you can run this service. 

To start, go to the [zd](./cmd/zd) directory:

``` bash
cd cmd/zd
```

Build the Go binary:

``` bash
go build -o build/zd -buildvcs=false
```

Run the application:

``` bash
./build/zd
```
