#!/bin/sh
set -e

echo "Setting permissions..."
id -u appuser &>/dev/null || adduser -S -D -H -h /app appuser
chown -R appuser /app
cd /app

exec su appuser -m -s /bin/sh -c "$@"