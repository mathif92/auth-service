#!/bin/bash
set -e
echo "HELLO"

host="$1"
port="$2"
shift 2
cmd="$@"

until nc -z -v -w30 "$host" "$port"; do
  echo "Waiting for MySQL to be available at $host:$port..."
  sleep 1
done

>&2 echo "MySQL is up and running - executing command"
exec $cmd