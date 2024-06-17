# Yampa
A real time streaming tool for stocks and cryptocurrencies. 


## Getting Started

1) Create a docker network for the Redpanda containers:

    ```bash
    docker network create --driver bridge --attachable redpanda-net
    ```

4) Start your Redpanda cluster:

    ```bash
    cd redpanda && docker compose up --remove-orphans -d && cd ..
    ```


3) Start streaming trades from Coincap into Redpanda:

    - Copy the contents of `yampa-cli/.env.example` into `yampa-cli/yampa.env`.
    - Start the yampa-cli container:

        ```bash
        cd yampa-cli && docker compose up --build --remove-orphans -d && cd ..
        ```

5) Go to [localhost:18080](http://localhost:18080) to view the trade data from the Redpanda Console UI.
