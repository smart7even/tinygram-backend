# Simple todo app backend written in golang

- Provides REST and GRPC interfaces
- Uses MySQL as database

## API

REST schema is declared in `api.yaml` and GRPC go client is defined in `client` folder.

To regenerate golang GRPC code use the following command (You need protoc installed):

```
protoc --go_out=internal/transport/grpc_handler --go_opt=paths=source_relative --go-grpc_out=internal/transport/grpc_handler --go-grpc_opt=paths=source_relative todo.proto
```

## Installation

### Run with docker compose

specify following variables in .env:

- DB_USER (MySQL user)
- DB_PASS (MySQL password)
- HTTP_ADRESS (adress on which REST API will be available)
- GRPC_ADRESS (adress on which GRPC API will be available)

.env example:

```
DB_USER=user
DB_PASS=password
HTTP_ADRESS=127.0.0.1:8080
GRPC_ADRESS=127.0.0.1:8081
```

Then run with `docker compose up`

Well done. Now try to use the REST api using swagger and `api.yaml` schema. Or try to use GRPC client (example of usage in client folder)

### Run without docker compose

- specify DB_CONNECTION_STRING key in .env file
- specify HTTP_ADRESS key in .env file (IP adress and port on which your HTTP service will be available)
- specify GRPC_ADRESS key in .env file (IP adress and port on which your GRPC service will be available)
- install [migrate](https://github.com/golang-migrate/migrate)
- create tables in your db by running `migrate -database mysql://{DB_CONNECTION_STRING} -path migrations up` (insert DB_CONNECTION_STRING value here)
- run using `go run main.go`

.env example:

```
HTTP_ADRESS=127.0.0.1:8080
GRPC_ADRESS=127.0.0.1:8081
DB_CONNECTION_STRING=user:password@tcp(host:port)/dbname
```
