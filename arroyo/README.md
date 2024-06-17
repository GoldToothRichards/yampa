# Arroyo Cluster

- This compose file creates a single node Arroyo cluster.
- The default configuration assumes that you are running a Redpanda cluster using `yampa/redpanda/docker-compose.yml`.
- The Redpanda brokers can be remote. The broker urls are configurable via environment variables.

## Configuration

- To use the default settings, copy the contents of `.env.example` into a new file called `.env`.
- To customize the settings, modify `.env` accordingly.

## Setup

- Create a docker network for the Arroyo cluster (if it does not already exist):

  ```bash
  docker network create --driver bridge --attachable redpanda-net
  ```

- Start the Arroyo cluster:

  ```bash
  docker compose up -d
  ```

- The Arroyo UI should now be available at [http://localhost:8000](http://localhost:8000).
