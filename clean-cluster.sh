#! /bin/bash

docker volume rm $(docker volume ls | awk '/redis/{ print $2 }')
