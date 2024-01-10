#!/bin/sh


# Exit on first failure
set -e

# Command Line Arguments
if [ "$1" = "-h" ] || [ "$1" = "--help" ]; then
  echo "This is a script that simplifies running RabbitMQ in a docker container."
  echo "The customization options you have is providing the username and password for the default management account.\n"
  echo "By default the username and password will be 'guest' but it is recommended to provide your own."
  echo "Example usage:\n\trabbitmq.sh admin supersecretpassword"
  exit 0
fi

USER=$1
PASSWORD=$2

# RabbitMQ Defaults
DEFAULT_USER="guest"
DEFAULT_PASSWORD="guest"

if [ "$USER" = "" ]; then
  echo "You have not set a user, \nWill fallback to the default user '$DEFAULT_USER'.\n"
  USER="$DEFAULT_USER"
fi

if [ "$PASSWORD" = "" ]; then
  echo "You have not set a password, \nWill fallback to the default password '$DEFAULT_PASSWORD'.\n"
  PASSWORD="$DEFAULT_PASSWORD"
fi

CONTAINER_ID=$(docker run --rm -d --hostname my-rabbit --name management-broker -e RABBITMQ_DEFAULT_USER=$USER -e RABBITMQ_DEFAULT_PASS=$PASSWORD -p 5672:5672 -p 15672:15672 rabbitmq:3-management)

echo "The RabbitMQ container has been started, the container id is: '$CONTAINER_ID'."
echo "You can login to the admin panel at http://localhost:15672 with the credientials below:"
echo "Username: $USER\nPassword: $PASSWORD\n"

