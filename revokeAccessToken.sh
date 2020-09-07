#!/bin/sh
CLIENT_ID=$(cat $HOME/twitch/clientId)
ACCESS_TOKEN=$(cat $TWITCH_HOME/accessToken)

curl -X POST -G https://id.twitch.tv/oauth2/revoke \
  -d client_id=${CLIENT_ID} \
  -d token=${ACCESS_TOKEN}
