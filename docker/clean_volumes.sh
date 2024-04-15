#! /bin/bash

docker volume rm $(docker volume ls | awk '!/NAME/{ print $2 }')
