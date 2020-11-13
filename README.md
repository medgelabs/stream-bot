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
