# Yampa CLI

A CLI utility for streaming real-time trading data.

## Description

The purpose of this tool is to stream trading data from a variety of upstreams sources, convert them to a consistent schema, and output them to variety of downstream sinks. All trade streams produced by the yampa-cli contain messages with the following schema:

```json
{
    "Source": "string",
    "Base": "string",
    "Quote": "string",
    "Exchange": "string",
    "VolumeBase": "float64",
    "VolumeQuote": "float64",
    "Price": "float64",
    "Timestamp": "string",
}
```

The available upstream options are:

- [Coincap](https://coincap.io/)
- [Alpaca](https://alpaca.markets/docs/)

The available downstream options are:

- Stdout
- Local file
- Kafka
- [Pusher](https://pusher.com/)

The available output formats are:

- JSON
- Avro

## Configuration

The yampa-cli can read configuration information directly from the environment or from an env file. To get started, copy the contents of __.env.example__ into __yampa.env__. Update the values as necessary. The default broker URL's and schema registry URL can be used with the Redpanda cluster created by `redpanda/docker-compose.yml`.

 The following environment variables are supported:

- KAFKA_BROKER_URLS

- SCHEMA_REGISTRY_URL

- KAFKA_MESSAGE_FORMAT

- KAFKA_MESSAGE_INCLUDE_KEY

- ALPACA_KEY

- ALPACA_SECRET

- PUSHER_APP_ID

- PUSHER_KEY

- PUSHER_SECRET

- PUSHER_CLUSTER

- PUSHER_CHANNEL

- PUSHER_EVENT_NAME

- PUSHER_BATCH_SIZE

## Usage

Stream to the terminal:

```bash
yampa-cli stream [coincap/alpaca]
```

Stream to a local json log file:

```bash
yampa-cli stream [coincap/alpaca] --file $PATH_TO_FILE
```

Stream to kafka:

```bash
yampa-cli stream [coincap/alpaca] --kafka --topic $TOPIC_NAME
```

Silence terminal output:

```bash
yampa-cli stream [coincap/alpaca] --quiet
```

Stream to multiple destinations:

```bash
yampa-cli stream [coincap/alpaca] --kafka --pusher --file $PATH_TO_FILE 
```

## Docker Usage

Build the image:

```bash
docker build . -t yampa-cli
```

Run a command:

```bash
docker run --rm -it \
--env-file yampa.env \
--network redpanda-net \
--volume /var/lib/yampa:/var/lib/yampa \
yampa-cli [COMMAND]
```

- __Note:__ The yampa-cli container must run on the same `redpanda-net` docker network as the Redpanda cluster when using the `redpanda/docker-compose.yml` setup with the Kafka sink. The container requires a volume mount in order to persist the data files generated when using the `--file` option.
