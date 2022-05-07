# Simple todo app backend written in golang

Uses postgres as database.

To run:

- specify DB_CONNECTION_STRING key in .env file
- specify ADRESS key in .env file (IP adress and port on which your service will be available)
- install [migrate](https://github.com/golang-migrate/migrate)
- create tables in your db by running `migrate -database DB_CONNECTION_STRING -path migrations up` (insert DB_CONNECTION_STRING value here)
- run using `go run main.go`

.env example

```
ADRESS=127.0.0.1:8080
DB_CONNECTION_STRING=postgresql://user:password@host:port/dbname
```
