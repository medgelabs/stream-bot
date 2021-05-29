# Stream Bot

Twitch ChatBot built on stream! Utilizes Twitch IRC to interact with Twitch Chat.

* Subscriptions, bit donations, and raid notifications are extracted from IRC messages
* Commands are also extracted from IRC messages
* All messages sent _FROM_ the bot _TO_ chat are done via IRC PRIVMSG commands

## Build & Run

The bot should compile to Windows/Mac/Linux native binaries with `go build`.
Linux and Mac have been tested. (Sorry Windows)

## CLI

## Secrets

A `Secret Store` provides secrets to the app (`/secrets` package). Currently supported options
are:

* Environment variables: `TWITCH_TOKEN`

## config.yaml

Configuration for the Bot can be defined in a `config.yaml` file found in the same directory
as the binary / root of the project. Config is by channel, like so:

```
CHANNEL_NAME:
  feature:
    config: ....
```

Each feature's configuration can be found in the following sections.

### MessageFormat and Templates

The `text/template` package is used for any messageFormat config keys. These configs are used
to define what shape of message is sent to chat on an event. For example: what is sent
to a user when they subscribe if the subscription thanks feature is turned on.

The `Event` struct in `bot/event.go` is what is used to bind to the template. As such, exported
fields can be used in config.yaml to form a messageFormat string. For example:

```
medgelabs:
  subs:
    messageFormat: "{{.Sender}} renewed their Lab Assistant role!"
```

`{{.Sender}}` is a placeholder that will be filled in with the user that just subscribed. Take
a look at the `Event` struct to see all available options.

## Greeter

Auto-Greeter will greet viewers on their first chat message. It does NOT greet
lurkers (i.e no greetings on JOIN / PART)!

The message sent is a Go fmt-style string stored in `config.yaml`. Different messages
can be set for different channels:

```
CHANNEL_NAME:
  greeter:
    messageFormat: "Welcome @{{.Sender}}!"
```

The only variable injected is the username. Any other substitutions are ignored.

## Commands

There are two sets of commands:

* `config.yaml` - defines simple response commands and command aliases
* `bot/commandHandler.go` - defines commands that require more complex logic

## Subscribers

Subscribers can be sent a message on a new subscription. The message sent is
found in `config.yaml` and can be set for different channels:

```
CHANNEL_NAME:
  subs:
    messageFormat: "Thank you for the subscription, @{{.Sender}}!"
  giftsubs:
    messageFormat: "Thank you for gifting a sub to @{{.Recipient}}, @{{.Sender}}!"
```

## Bits

Bits Donations can be automatically thanked on donation. The message sent is
found in `config.yaml` and can be set for different channels:

```
CHANNEL_NAME:
  bits:
    messageFormat: "Thank you for the {{.Amount}} bits, @{{.Sender}}!"
```

## TODO - Followers

Followers API doesn't appear to be in IRC or PubSub.
EventSub appears to have an event for it

## TODO - Emote Stats

If enabled, and by setting the following `config.yaml` entry:

```
CHANNEL_NAME:
  emotePrefix: (YOUR EMOTE PREFIX HERE)
```

The bot will analyze messages for usages of emotes with the given prefix. A single
emote is only counted once per message (i.e sending the same emote 4 times in one
message will only count as 1 usage).

TODO - how are these accessed?

## TODO - QoL

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

Mass Gift Subs:

Note: This message shows first, but also generates 5 individual Sub USERNOTICE messages where msg-id=subgift

```
map[badge-info:subscriber/3 badges:broadcaster/1,subscriber/0 color:#2E8B57 display-name:medgelabs emotes: flags: id:059d0e94-59d5-4ded-95d5-c5d6bb58ba0b login:medgelabs mod:0 msg-id:submysterygift msg-param-mass-gift-count:5 msg-param-origin-id:07\s23\s61\s03\sa4\sa1\s5a\sb0\s72\s4a\s72\s2d\s08\sdd\sda\s57\sc8\s24\sf6\s17 msg-param-sender-count:0 msg-param-sub-plan:1000 room-id:62232210 subscriber:1 system-msg:medgelabs\sis\sgifting\s5\sTier\s1\sSubs\sto\smedgelabs's\scommunity! tmi-sent-ts:1607696526324 user-id:62232210 user-type:] tmi.twitch.tv USERNOTICE #medgelabs
```

Raid:

```
map[badge-info:founder/4 badges:vip/1,founder/0,bits-leader/3 color:#FF0000 display-name:RAIDER emotes: flags: id:05f87ef5-793e-424f-aa52-d10e96073925 login:RAIDER mod:0 msg-id:raid msg-param-displayName:RAIDER msg-param-login:RAIDER msg-param-profileImageURL:https://URL msg-param-viewerCount:4 room-id:66490190 subscriber:1 system-msg:4\sraiders\sfrom\sRAIDER\shave\sjoined! tmi-sent-ts:1607544051081 user-id:RAIDER_ID user-type:]
<<< tmi.twitch.tv USERNOTICE #CHANNEL
```

Host:
TODO - this was initiated by BTTV. Does Twitch do this?

```
jtv!jtv@jtv.tmi.twitch.tv PRIVMSG medgelabs :HOSTER is now hosting you.
```
