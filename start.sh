#!/bin/sh

VAULT_PATH="secret/twitchToken"

docker-compose down
docker-compose up -d redis vault
sleep 2
export VAULT_TOKEN=$(docker-compose logs vault | grep "Root Token" | cut -d ":" -f2 | cut -d " " -f2)

TWITCH_TOKEN=$(cat $TWITCH_HOME/accessToken | cut -d "=" -f2)
vault kv put $VAULT_PATH token=${TWITCH_TOKEN}
docker-compose up --build bot
