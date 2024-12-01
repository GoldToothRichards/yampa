# Yampa

A real time streaming tool for stocks and cryptocurrencies.

## 📚 Description

Yampa is a data engineering project that streams real-time cryptocurrency trade data from [Coincap](https://coincap.io) into a Redpanda cluster and processes it using ClickHouse for real-time analytics. The processed data can be used for various analyses, including the statistical distributions of trading patterns.

## 🛠️ Tech Stack

- **Data Collection**: Yampa CLI (Go)
- **Data Streaming**: Redpanda (Kafka API-compatible)
- **Stream Processing**: Redpanda Connect
- **Storage & Analytics**: ClickHouse

## 🔗 Related Projects

[Trade Distributions](https://trade-distributions.vercel.app) - A companion project that visualizes the statistical patterns in cryptocurrency trading data collected by Yampa. The analysis reveals interesting patterns in trade volume distributions that resemble well-known probability distributions.

## 🚀 Getting Started

### Stream trades into Redpanda with the Yampa CLI

1. Create a docker network for the Redpanda containers:

   ```bash
   docker network create --driver bridge --attachable redpanda-net
   ```

2. Start your Redpanda cluster:

   ```bash
   cd redpanda && docker compose up -d && cd ..
   ```

3. Start streaming trades from Coincap into Redpanda:

   - Copy the contents of `yampa-cli/.env.example` into `yampa-cli/.env`
   - Start the yampa-cli container:

     ```bash
     cd yampa-cli && docker compose up --build -d && cd ..
     ```

4. View the trade data in the Redpanda Console UI at [localhost:8080](http://localhost:8080)

### Setup your ClickHouse database for persistant storage

1. Copy the contents of `clickhouse/.env.example` into `clickhouse/.env`
2. Start your ClickHouse container:

   ```bash
   cd clickhouse && docker compose up -d && cd ..
   ```

3. Create the `trades` database and `raw_trades` table in the ClickHouse UI at [localhost:8123](http://localhost:8123) using the `clickhouse/trades/tables/raw_trades.sql` schema.

### Persist the trade data to ClickHouse using Redpanda Connect

1. Copy the contents of `redpanda-connect/.env.example` into `redpanda-connect/.env`
2. Start your Redpanda Connect cluster:

   ```bash
   cd redpanda-connect && docker compose up -d && cd ..
   ```

3. Create the ClickHouse sink connector in the Redpanda Console UI using the `redpanda-connect/connectors/clickhouse-sink.json` configuration.

## 📄 License

This project is open source and available under the [MIT License](LICENSE).
