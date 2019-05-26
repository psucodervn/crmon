Container Registry Monitor
[![Build Status](https://travis-ci.org/psucodervn/crmon.svg?branch=master)](https://travis-ci.org/psucodervn/crmon)

----------

## Description

Monitoring for new image updates in Google Container Registry.
Can be use as a standalone app or import as a library.
You can write your own subscriber that implements `Subscriber` interface to deal with new events.

```golang
type Event struct {
  Action string `json:"action"`
  Tag    string `json:"tag"`
  Digest string `json:"digest"`
}

type Subscriber interface {
  Name() string
  Init() error
  Cleanup() error
  OnReceive(event Event) error
}
```

## Features:
  - [x] [Print to console](#print-to-console)
  - [x] [Send message to Slack](#send-message-to-slack)
  - [x] Run a shell command
  - [ ] Send message to Mattermost
  - [ ] Pull image and update docker container
  - [ ] ...

## Examples

### Print to console
![Print to console](/docs/console.png)

### Send message to Slack
![Send message to Slack](/docs/slack.png)
