#!/bin/sh

docker run -rm -d --hostname my-rabbit --name management-broker -e RABBITMQ_DEFAULT_USER= -e RABBITMQ_DEFAULT_PASS= -p 5672:5672 -p 15672:15672 rabbitmq:3-management
