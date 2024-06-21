#!/bin/bash

# e-commerce-chat-microservice_https-portal-data以外を削除する
volumes_to_delete=$(docker volume ls -q | grep -v 'e-commerce-chat-microservice_https-portal-data')

for volume in $volumes_to_delete; do
  sudo docker volume rm $volume
done
