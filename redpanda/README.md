# Redpanda Cluster

## Setup

- Download the connect plugins:
  - https://www.confluent.io/hub/confluentinc/kafka-connect-avro-converter
  - https://www.confluent.io/hub/clickhouse/clickhouse-kafka-connect
  - https://www.confluent.io/hub/tabular/iceberg-kafka-connect

- Extract the plugins to `/opt/kafka/connect-plugins`:

    ```bash
    unzip *.zip -d /opt/kafka/connect-plugins
    ```

- Start the cluster:

    ```bash
    docker compose up -d
    ```

- The Redpanda Console UI is availble at [localhost:18080](http://localhost:18080).
- The Redpanda Broker URL's are `http://redpanda-1:29092,http://redpanda-2:29093,http://redpanda-3:29094`.
- The Redpanda Schema Registry URL is `http://redpanda-1:8081`.
