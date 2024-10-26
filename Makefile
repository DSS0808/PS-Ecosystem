POSTGRES_PASSWORD ?= password
POSTGRES_USER ?= user
POSTGRES_DB ?= users
PROTO_PATH ?= proto/user

.PHONY: docker_build docker_stop run_postgres build_client client_list client_done proto swag

swag: ## generate swagger
	swag init -d cmd,internal/server,internal/models

proto: ## generate go code from protofiles
	protoc \
		--go_out=${PROTO_PATH} \
		--go-grpc_out=${PROTO_PATH} \
		--proto_path=${PROTO_PATH} \
		messages.proto \
		service.proto

run: ## run server in local machine
	go run cmd/main.go

docker_build: ## build and run app+DB in docker compose
	docker-compose -f build/docker-compose.yaml up --build

run_postgres: ## run postgres in docker
	POSTGRES_PASSWORD=${POSTGRES_PASSWORD} \
	POSTGRES_USER=${POSTGRES_USER} \
	POSTGRES_DB=${POSTGRES_DB} \
	docker-compose -f build/docker-compose.yaml up -d postgres

docker_stop: ## stop docker compose with app+DB
	docker-compose -f build/docker-compose.yaml down

build_client: ## build client
	mkdir -p bin
	go build -o bin/client  cmd/client/main.go

client_done: ## set status DONE of a task with given id. Need to add id=*needed id*
	go run cmd/client/main.go -c=done -i=${id}

client_list: ## get list of tasks using existion client
	./bin/client/main -c=list

.PHONY: help
help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'