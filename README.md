# Stream Bot

Twitch ChatBot built on stream! Utilizes Twitch IRC and Twitch PubSub APIs to
interact with Twitch Chat.

* All messages sent _FROM_ the bot _TO_ chat are done via IRC PRIVMSG commands
* PubSub is used for Subscriptions, Point Redemptions, and Bits since IRC does not send these

## Build & Run

The bot should compile to Windows/Mac/Linux native binaries with `go build`.
Linux and Mac have been tested. (Sorry Windows)

## CLI

The Bot is a CLI-based program and accepts CLI flags to configure behavior. Run the
binary with the `-h` flag to see available options.

## Secrets

Secrets required by the app:

* Twitch OAuth token

A `Secret Store` provides these to the app (`/secrets` package). Currently supported options
are:

* Hashicorp Vault (available via `docker-compose up vault`)
  * Secret is expected under `secrets/twitchToken`
* Environment variable: `TWITCH_TOKEN`

https://twitch.tv/medgelabs


## Greeter

Auto-Greeter will greet viewers on their first chat message. It does NOT greet
lurkers (i.e no greetings on JOIN / PART)!

The message sent is a Go fmt-style string stored in `config.yaml`. Different messages
can be set for different channels:

```
greeter:
  CHANNEL_NAME:
    messageFormat: "Welcome @%s!"
```

The only variable injected is the username. Any other substitutions are ignored.

## Commands

Commands are currently set in code in the `bot/commandHandler.go` file.

## TODO - Subscribers

Subscribers can be sent a message on a new subscription. The message sent is
found in `config.yaml` and can be set for different channels:

```
subscribers:
  CHANNEL_NAME:
    messageFormat: "Thank you for the subscription, @%s!"
```

The only variable injected is the username. Any other substitutions are ignored.

TODO resub?

## TODO - Followers

Followers API doesn't appear to be in IRC or PubSub. Needs investigation

## TODO - Bits

## TODO - Emote Stats

If enabled, and by setting the following `config.yaml` entry:

```
emotes:
  CHANNEL_NAME:
    prefix: medgel (YOU'RE EMOTE PREFIX HERE)
```

The bot will analyze messages for usages of emotes with the given prefix. A single
emote is only counted once per message (i.e sending the same emote 4 times in one
message will only count as 1 usage).

TODO - how are these accessed?

## TODO - QoL

* Use Mustache for message strings for more options?
* Simplify config.yaml to be channel-top to simplify for multiple features. i.e:
    ```
    CHANNEL_NAME:
      greeter:
        messageFormat: "..."
    ```
* Config without YAML file for more portability / quick-run?
* Automated OAuth token

## IRC Messages

Bits:

```
map[badge-info: badges: bits:100 color: display-name:CHEERER emotes: flags: id:592bcd9d-649f-4a63-bf49-70250d37cd95 mod:0 room-id:ROOM_ID subscriber:0 tmi-sent-ts:1607510776445 turbo:0 user-id:CHEERER_ID user-type:]
> USER: Cheer100
```

Whisper:

```
map[badges: color:#FF7F50 display-name:SOURCE emotes: message-id:82 thread-id:THREAD_ID turbo:0 user-id:USER_ID user-type:]
<<< SOURCE WHISPER target :MESSAGE HERE
```

Timeout:

```
tmi.twitch.tv CLEARCHAT #CHANNEL :USERNAME
```

First Sub:

```
map[badge-info: badges: color:#FFAFC8 display-name:SUBSCRIBER emotes: flags: id:433dc84a-6dae-4f45-ad1a-1d2a27c086f2 login:SUBSCRIBER mod:0 msg-id:sub msg-param-cumulative-months:1 msg-param-months:0 msg-param-multimonth-duration:1 msg-param-multimonth-tenure:0 msg-param-should-share-streak:0 msg-param-sub-plan:1000 msg-param-sub-plan-name:Tier1 msg-param-was-gifted:false room-id:ROOM_ID subscriber:1 system-msg:SUBSCRIBER\ssubscribed\sat\sTier\s1. tmi-sent-ts:1607629647151 user-id:267451666 user-type:] tmi.twitch.tv USERNOTICE
```

Gift Sub:

```
map[badge-info:subscriber/6 badges:moderator/1,subscriber/6,bits/100000 color:#2E8B57 display-name:medgelabs emotes: flags: id:dcd7b498-c669-400c-8218-236570110258 login:medgelabs mod:1 msg-id:subgift msg-param-gift-months:1 msg-param-months:2 msg-param-origin-id:ORIGIN_ID msg-param-recipient-display-name:GIFT_RECIPIENT msg-param-recipient-id:GIFT_RECIPIENT_ID msg-param-recipient-user-name:GIFT_RECIPIENT msg-param-sender-count:38 msg-param-sub-plan:1000 msg-param-sub-plan-name:Channel\sSubscription\s(CHANNEL) room-id:119432399 subscriber:1 system-msg:medgelabs\sgifted\sa\sTier\s1\ssub\sto\sGIFT_RECIPIENT!\sThey\shave\sgiven\s38\sGift\sSubs\sin\sthe\schannel! tmi-sent-ts:1607514324648 user-id:SENDER_ID user-type:mod]
<<< tmi.twitch.tv USERNOTICE #CHANNEL
```

Raid:

```
map[badge-info:founder/4 badges:vip/1,founder/0,bits-leader/3 color:#FF0000 display-name:RAIDER emotes: flags: id:05f87ef5-793e-424f-aa52-d10e96073925 login:RAIDER mod:0 msg-id:raid msg-param-displayName:RAIDER msg-param-login:RAIDER msg-param-profileImageURL:https://URL msg-param-viewerCount:4 room-id:66490190 subscriber:1 system-msg:4\sraiders\sfrom\sRAIDER\shave\sjoined! tmi-sent-ts:1607544051081 user-id:RAIDER_ID user-type:]
<<< tmi.twitch.tv USERNOTICE #CHANNEL
```
