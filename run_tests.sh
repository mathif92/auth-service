#!/bin/sh

# Function to check if the MySQL database is up and running
check_mysql() {
  mysqladmin_ping_command="mysqladmin ping -h localhost -P 3310 -u root"
  echo "Executing: $mysqladmin_ping_command"
  if $mysqladmin_ping_command; then
    echo "MySQL is up and running."
  else
    echo "MySQL is not yet available. Waiting for 1 second..."
    sleep 1
    check_mysql # Retry the check
  fi
}

# Wait for the MySQL service to be available
check_mysql

# Run the migrations
migrate -path ./db/migration/ -database "mysql://auth:auth@(localhost:3310)/auth" up


# Run your tests
go test -v ./... -covermode=atomic -coverprofile=coverage.out -coverpkg=./... -count=1  -race -timeout=30m