# Redpanda Connect Cluster

## Configuration

- To use the default settings, copy the contents of `.env.example` into a new file called `.env`.
- These settings assume that you have already configued a Redpanda cluster using `yampa/redpanda/docker-compose.yml`.
- To customize the settings, modify `.env` accordingly.

## Setup

- Download the connect plugins:

  - https://www.confluent.io/hub/confluentinc/kafka-connect-avro-converter
  - https://www.confluent.io/hub/clickhouse/clickhouse-kafka-connect
  - https://www.confluent.io/hub/tabular/iceberg-kafka-connect

- Extract the plugins to `/opt/kafka/connect-plugins`:

  ```bash
  unzip *.zip -d /opt/kafka/connect-plugins
  ```

- Create a docker network for the Redpanda containers (if it does not already exist):

  ```bash
  docker network create --driver bridge --attachable redpanda-net
  ```

- Start the cluster:

  ```bash
  docker compose up -d
  ```

- The Connect cluster should now be available from the Redpanda Console UI at [http://localhost:18080/connect-clusters/Primary](http://localhost:18080/connect-clusters/Primary).
