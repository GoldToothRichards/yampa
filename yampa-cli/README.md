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
  "Timestamp": "string"
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

- To use the default settings, copy the contents of `.env.example` into a new file called `.env`.
- The default configuration assumes that you are running a Redpanda cluster using `yampa/redpanda/docker-compose.yml`.
- All of the environment settings are optional, but most of the upstream/downstream options require some of them to be set.
- To customize the settings, modify `.env` accordingly.

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
--env-file .env \
--network redpanda-net \
--volume /var/lib/yampa:/var/lib/yampa \
yampa-cli [COMMAND]
```

- **Note:** The yampa-cli container must run on the same `redpanda-net` docker network as the Redpanda cluster when using the `redpanda/docker-compose.yml` setup with the Kafka sink. The container requires a volume mount in order to persist the data files generated when using the `--file` option.
