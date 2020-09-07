#!/bin/sh

if [ -z "$TWITCH_HOME" ]; then
  echo "ERROR: TWITCH_HOME not set"
  exit 1
fi

CLIENT_ID=$(cat $TWITCH_HOME/clientId)
CLIENT_SECRET=$(cat $TWITCH_HOME/clientSecret)
ACCESS_TOKEN_PATH=$TWITCH_HOME/accessToken

curl -X POST -G 'https://id.twitch.tv/oauth2/token' \
  -d client_id=${CLIENT_ID} \
  -d client_secret=${CLIENT_SECRET} \
  -d grant_type=client_credentials \
  -d scope=channel:read:subscriptions+chat:read+chat:edit+whispers:read+whispers:edit \
  | jq -r '"TWITCH_TOKEN=oauth:" + .access_token' > ${ACCESS_TOKEN_PATH}
