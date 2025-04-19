#!/bin/sh

chown -R appuser:appgroup /app
exec su -s /bin/sh appuser -c "$@"