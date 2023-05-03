# auth-service
This repo contains the auth-service application code written in Golang.
It implements a Restful API that provides JWT authentication, roles & actions.

## Running the app locally
For running the app locally it's needed to execute:

```
make db-up
make migrate-up
go run ./cmd/app
```
