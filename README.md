# Simple todo app backend written in golang

Uses postgres as database.

To run:

- specify DB_CONNECTION_STRING key in .env file
- install [migrate](https://github.com/golang-migrate/migrate)
- create tables in your db by running `migrate -database DB_CONNECTION_STRING -path migrations up` (insert DB_CONNECTION_STRING value here)
- run using `go run main.go`

```
DB_CONNECTION_STRING=postgresql://user:password@host:port/dbname
```
