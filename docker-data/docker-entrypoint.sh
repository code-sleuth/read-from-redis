#!/bin/sh

if [ "$1" = 'redis-cluster' ]; then
    sleep 10
    echo "yes" | redis-cli --cluster create 173.17.0.2:7002 173.17.0.3:7003 173.17.0.4:7004 --cluster-replicas 0
    echo "DONE"
else
  exec "$@"
fi
