# Redpanda Cluster

- This compose file creates a 3-node Redpanda cluster which runs on a single server, together with the Redpanda Console UI.
- The Redpanda UI configuration assumes that you are running a Redpanda Connect cluster using `yampa/redpanda-connect/docker-compose.yml`.
- The Redpanda Connect cluster can be remote. The host addresses are configurable via environment variables.

## Configuration

- To use the default settings, copy the contents of `.env.example` into a new file called `.env`.
- The default environment settings assume that you are running the Redpanda cluster and Redpanda Connect cluster on the same local server.
- To customize the settings, modify `.env` accordingly.

## Setup

- Create a docker network for the Redpanda cluster (if it does not already exist):

  ```bash
  docker network create --driver bridge --attachable redpanda-net
  ```

- Start the Redpanda cluster:

    ```bash
    docker compose up -d
    ```

- The Redpanda Console UI should now be availble at [localhost:8080](http://localhost:8080).
