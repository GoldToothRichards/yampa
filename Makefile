start:
	cd redpanda && docker compose up --remove-orphans -d && cd ../yampa-cli && docker compose up --remove-orphans -d && cd ../redpanda-connect && docker compose up --remove-orphans -d

stop:
	cd redpanda && docker compose down && cd ../yampa-cli && docker compose down && cd ../redpanda-connect && docker compose down
