#!/bin/sh

VAULT_PATH="secret/twitchToken"

docker-compose up -d redis vault
sleep 2
export VAULT_TOKEN=$(docker-compose logs vault | grep "Root Token" | cut -d ":" -f2 | cut -d " " -f2)

stty -echo
read -p "Twitch Token: " TWITCH_TOKEN
stty echo
vault kv put $VAULT_PATH token=${TWITCH_TOKEN}
