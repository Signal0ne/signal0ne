#!/bin/sh
# healthcheck.sh

SOCKET_FILE="/var/run/d.sock"

if [ ! -S "$SOCKET_FILE" ]; then
  echo "Socket file does not exist"
  exit 1
fi

if ! nc -U -z "$SOCKET_FILE"; then
  echo "Cannot connect to socket"
  exit 1
fi

exit 0