# Yampa

A real time streaming tool for stocks and cryptocurrencies.

## Getting Started

### Stream trades into Redpanda with the Yampa CLI

1. Create a docker network for the Redpanda containers:

   ```bash
   docker network create --driver bridge --attachable redpanda-net
   ```

2. Start your Redpanda cluster:

   ```bash
   cd redpanda && docker compose up --remove-orphans -d && cd ..
   ```

3. Start streaming trades from Coincap into Redpanda:

   - Copy the contents of `yampa-cli/.env.example` into `yampa-cli/.env`.
   - Start the yampa-cli container:

     ```bash
     cd yampa-cli && docker compose up --build --remove-orphans -d && cd ..
     ```

4. Go to [localhost:8080](http://localhost:8080) to view the trade data from the Redpanda Console UI.

### Compute average prices in real time with Arroyo

1. Start your Arroyo cluster:

   ```bash
   cd arroyo && docker compose up --remove-orphans -d && cd ..
   ```

   Go to [localhost:8000](http://localhost:8000) to view the Arroyo UI.

2. From the Arroyo UI, create a new Kafka source connector called _trades_source_ which consumes from the _trades_ topic in Redpanda.

3. Create a new Kafka sink connector called _prices_sink_ which produces to the _prices_ topic in Redpanda.

4. Start the streaming pipeline using `arroyo/pipelines/prices.sql`.
