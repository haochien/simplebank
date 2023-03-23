#!/bin/sh
# run in bin shell because bash shell is not available in alpine image

# -e: make sure that the script will exit immediately if a command return a non-zero status
set -e 

echo "run db migration"
# for prod, need to set all var in .env to the environment first
source /app/app.env
# call the migrate binary; db source is from env var defined in docker-compose.yaml
/app/migrate -path /app/migration -database "$DB_SOURCE" -verbose up

echo "start the app"
# $@: takes all parameters passed to the script and run it (i.e. command in CMD[] in dockerfile)
exec "$@"