# ZD

This is the Zendesk Service that is responsible for producing User Events.

## Running ZD

### Service Dependencies

This service depends on 2 other's, one would be a service that provides user and event data as this is needed to generate random user events. The second required service is an event queue, this is used as a way for consumers of this service to receive user events for further processing.

When it comes to the service that provides ZD user and event data, we would use the service called [pd-users-api](https://github.com/TSE-Coders/pd-users-api). For the event queue we make use of [RabbitMQ](https://www.rabbitmq.com/). 

### Single Host/Local Environment

To run this service you will need to have Golang version 1.21. When it comes to RabbitMQ you can either install it directly or run it with Docker. For your convenience I created a script [rabbitmq.sh](rabbitmq.sh) that you can use to start a RabbitMQ instance in a container. 

Once your RabbitMQ instance is running, start up the [pd-users-api](https://github.com/TSE-Coders/pd-users-api) service by following it's documentation. 

Now before starting up this service, you will need to set some environment variables, the list of the required variables can be found in the file [.env.example](./.env.example). 

Now that the required services and the environment variables are set, you can run this service. 

To start the service go to the app directory:

``` bash
cd cmd/app
```

Run the application:

``` bash
go run main.go
```
