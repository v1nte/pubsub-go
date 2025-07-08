# PubSub Patern in Go

Implementation of Publisher and Suscribers patter in Go

## Usage

### Connect

You'll need a websocket client, either build one or use one like [wscat](https://github.com/websockets/wscat) as I do.
Connect to `ws://localhost:9876/ws`

All comunnication must be in json format

### First Message

First message have to be a `name` otherwise it'll close the connection:
e.g.:

```json
{ "name": "Jhon" }
```

### Subscribe

Then you can subscribe to a topic with the `SUB` command to receive messages from a topic:

```json
{ "Command": "SUB", "Topic": "news" }
```

### Publish

Or, send message to a topic with `PUB` command. You can send a message without being subscribed to it. Feature, not a bug.

```json
{
  "Command": "PUB",
  "Topic": "Sports",
  "Message": "Germany 7-1 Brazil"
}
```

## Features

- [x] Subscribe
- [x] Unsubscribe
- [x] Unsubscribe all
- [x] Send msg via Broker
- [x] Handle multiples clients
- [x] Include name in messages
- [x] Use channels
- [x] Dockerize
- [ ] Logs
- [x] DB connection
- [ ] Auth
