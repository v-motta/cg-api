#!/bin/bash

IMAGE_NAME="cost-guardian-api-i"

#docker pull $IMAGE_NAME

CONTAINER_NAME="cost-guardian-api-c"
CONTAINER_PORT="5000"

docker stop $CONTAINER_NAME
docker rm $CONTAINER_NAME

DB_HOST="172.17.0.3"
DB_PORT="5432"
DB_NAME="postgres"
DB_USER="mottinha"
DB_PASSWORD="123456"
DB_SSLMODE="disable"

docker run -dt \
    -p $CONTAINER_PORT:$CONTAINER_PORT \
    -e DB_HOST=$DB_HOST \
    -e DB_PORT=$DB_PORT \
    -e DB_NAME=$DB_NAME \
    -e DB_USER=$DB_USER \
    -e DB_PASSWORD=$DB_PASSWORD \
    -e DB_SSLMODE=$DB_SSLMODE \
    -h $CONTAINER_NAME \
    --name $CONTAINER_NAME \
    $IMAGE_NAME

if [ $? -eq 0 ]; then
    echo "Container $CONTAINER_NAME is running on port $CONTAINER_PORT"
else
    echo "Failed to run container $CONTAINER_NAME"
fi
