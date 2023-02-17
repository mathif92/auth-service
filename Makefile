PACKAGES_PATH = $(shell go list -f '{{ .Dir }}' ./...)

before-push: tidy fmt up test

tidy:
	@echo "=> Executing tidy & vendor"
	@go mod tidy
	@go mod vendor

fmt:
	@echo "=> Executing goimports"
	@goimports -w $(PACKAGES_PATH)

# ==============================================================================
# Running tests within the local computer

test: up
	@echo "=> Executing go test"
	@go test -v ./... -covermode=atomic -coverprofile=coverage.out -coverpkg=./... -count=1  -race -timeout=30m

coverage: test
	@echo "=> Running coverage report"
	@go tool cover -html coverage.out

db-up:
	@echo "=> Starting mysql docker compose"
	@docker compose -f mysql-compose.yaml up --detach --remove-orphans

db-down:
	@echo "=> Destroying mysql docker compose"
	@docker compose -f mysql-compose.yaml down --remove-orphans

up:
	@echo "=> Running auth-service"
	@go run cmd/app/main.go

k8s-up:
	@echo "=> Running auth-service in local Minikube k8s cluster"
	@kubectl apply -f deploy/k8s/local/config.yaml
	@kubectl apply -f deploy/k8s/local/service.yaml
	@kubectl apply -f deploy/k8s/local/mysql-service.yaml
	@kubectl create configmap proxysql-configmap --from-file=deploy/k8s/local/proxysql.cnf
	@kubectl apply -f deploy/k8s/base/deployment.yaml

k8s-clean:
	@kubectl delete -f deploy/k8s/base/deployment.yaml
	@kubectl delete -f deploy/k8s/local/service.yaml
	@kubectl delete -f deploy/k8s/local/mysql-service.yaml
	@kubectl delete -f deploy/k8s/local/config.yaml
	@kubectl delete configmap proxysql-configmap


# Administration

init-db:
	@migrate create -ext sql -dir db/migration -seq init_schema

migrate name=1:
	@migrate create -ext sql -dir db/migration -seq $(name)

migrateup:
	@migrate -path db/migration -database "mysql://auth:auth@tcp(localhost:3306)/auth?tls=false" -verbose up

migratedown:
	@migrate -path db/migration -database "mysql://auth:auth@tcp(localhost:3306)/auth?tls=false" -verbose down

