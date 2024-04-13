#!/bin/bash
sleep 10
echo "yes" | redis-cli --cluster create --cluster-replicas 1 $(getent hosts redis-node-1 | awk '{print $1}'):6379 $(getent hosts redis-node-2 | awk '{print $1}'):6380 $(getent hosts redis-node-3 | awk '{print $1}'):6381 $(getent hosts redis-node-1-replica | awk '{print $1}'):6382 $(getent hosts redis-node-2-replica | awk '{print $1}'):6383 $(getent hosts redis-node-3-replica | awk '{print $1}'):6384