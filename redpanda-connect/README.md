# Redpanda Connect Cluster

- This compose file creates a 2-node Redpanda Connect cluster which runs on a single server.
- The default configuration assumes that you are running a Redpanda cluster using `yampa/redpanda/docker-compose.yml`.
- The Redpanda brokers can be remote. The broker urls are configurable via environment variables.

## Configuration

- To use the default settings, copy the contents of `.env.example` into a new file called `.env`.
- The AWS settings are optional but are required to configure an S3 bucket for the Apache Iceberg Sink.
- To customize the settings, modify `.env` accordingly.

## Setup

- Download the connect plugins:

  - https://www.confluent.io/hub/confluentinc/kafka-connect-avro-converter
  - https://www.confluent.io/hub/clickhouse/clickhouse-kafka-connect
  - https://www.confluent.io/hub/tabular/iceberg-kafka-connect

- Extract the plugins to `/opt/kafka/connect-plugins`:

  ```bash
  mkdir -p /opt/kafka/connect-plugins && unzip *.zip -d /opt/kafka/connect-plugins
  ```

- Create a docker network for the Redpanda Connect cluster (if it does not already exist):

  ```bash
  docker network create --driver bridge --attachable redpanda-net
  ```

- Start the Redpanda Connect cluster:

  ```bash
  docker compose up -d
  ```

- The Connect cluster should now be visible from the Redpanda Console UI at [http://localhost:8080/connect-clusters/Primary](http://localhost:8080/connect-clusters/Primary).
