#!/usr/bin/env sh
KAFKADIR=/data/kafka
KAFKABIN=$KAFKADIR/kafka_2.11-0.11.0.3/bin/
KAFKACONF=$KAFKADIR/kafka_2.11-0.11.0.3/config

echo "> start zookeeper"
sh $KAFKABIN/zookeeper-server-start.sh --daemon $KAFKACONF/zookeeper.properties;

n = 3
for (( i = 0; i < n; i++ )); do
    echo "."
    sleep 1;
done

echo "> start kafka";
sh $KAFKABIN/kafka-server-start.sh --daemon $KAFKACONF/server.properties;
