# Stream Bot

Twitch ChatBot built on stream!

https://twitch.tv/medgelabs

## Greeter

Auto-Greeter will greet viewers on their first chat message. It does NOT greet
lurkers!

The message is a Go fmt-style string stored in `config.yaml`. Different messages
can be set for different channels:

```
greeter:
  CHANNEL_NAME:
    messageFormat: "Welcome @%s!"
```

The only variable injected is the username. Any other substitutions are ignored.

## Commands

Commands are currently set in code in the `bot/commandHandler.go` file.

## TODO - Followers / Subscribers

Followers / Subscribers can be sent a message automatically. The message sent is
found in `config.yaml` and can be set for different channels:

```
followers:
  CHANNEL_NAME:
    messageFormat: "Thanks for the follow, @%s!"
subscribers:
  CHANNEL_NAME:
    messageFormat: "Thanks for the subscription, @%s!"
```

The only variable injected is the username. Any other substitutions are ignored.

## TODO - Bits

## TODO - Emote Stats
