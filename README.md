# Metal Gear Online 1 - Discord Alert Bot

A simple discord bot that listens to events coming from the MGO1 server and executes a discord webhook on creation/deletion of a game.

---

The bot is intend to work on the /stream/events websocket made available by the REST API provided by the [TX-55 MGO1 repository](https://github.com/underscore-zi/tx55).

For the "offical" revived server run by SaveMGO.com this is exposed at:
 - `wss://api.mgo1.savemgo.com/api/v1/stream/events`

## Setup

The bot can take three arguments, each can either be specified as a command list argument, or as an environment variable. If both are provided the environment variable takes precedence.

1. The Discord Webhook url. Get this from discord, configure it how you want, it executes it with a message only. [`-discord` or `DISCORD_WEBHOOK`]
2. The Websocket URL to listen to for events. [`-socket` or `WEBSOCKET_URL`]
3. A Format string with a single %d in it to replace with the game id to generate a URL to the game. [`-url` or `URL_FORMAT_STR`]
