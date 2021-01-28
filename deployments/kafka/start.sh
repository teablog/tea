#!/usr/bin/env sh
KAFKADIR=/data/kafka
KAFKABIN=$KAFKADIR/kafka_2.11-0.11.0.3/bin/
KAFKACONF=$KAFKADIR/kafka_2.11-0.11.0.3/config

sh $KAFKABIN/zookeeper-server-start.sh --daemon $KAFKACONF/zookeeper.properties