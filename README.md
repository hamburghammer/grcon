# grcon

A basic Golang library for the
[RCON Protocol](https://developer.valvesoftware.com/wiki/Source_RCON_Protocol).

## Features

- Max control over the underlying connection.
- Works on a slow internet connection.
- Stable, reliable and tested implementation.
- Small API surface.
- Simple and powerfully error type system.
- Good documented.
- Offering different levels of abstraction.

## Structure

Small overview of the structure and most important files.

### Main

Inside the "main" package you will find the low-level implementation of the
grcon Protocol. The starting point to get into it is the [grcon.go](grcon.go)
File.

### Util

This is the location for helper functions. It is a collection to facilitate the
interaction with `Packet`s and `RemoteConsole`.

### Client

The spot to look into for a higher abstracted API to interact with a
`RemoteConsole`. It simplifies the interaction and hides all complexity.

The currently only existing client is the
[SimpleClient](client/simple_client.go) but I can imagine to add game/server
specific implementations in the future.

## Motivation

Make the best std lib that provides a low-level implementation but also offers
packages/ways to abstract the handling with single packets.

## License

This lib is licensed under the [MIT License](LICENSE) and contains parts of the
implementation from [james4k/rcon](https://github.com/james4k/rcon).

## Contributors

If you should encounter a bug or a missing feature don't hesitate to open an
issue or even submit a pull-request.
