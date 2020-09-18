# TCP and UDP echo app – client and server in one

This is a simple TCP and UDP echo server and client, designed for testing
container lifetime and reachability.

In server mode, it listens on both TCP and UDP on the port specified by the
`PORT` environment variable. It reads data sent by clients, and echos back
the data along with some additional context that's useful for debugging.

In client mode, you specify the protocol (TCP or UDP), the target host and an
interval, and it will create a new connection for every echo request.

## Features

- Client and server in one binary (controlled by env vars)
- TCP and UDP mode for both client and server
- JSON logs with lots of detail
- SIGINT/SIGTERM handler with extra logging
- The AWS ECS task ARN is included in the logs, if available

## Configuration

All configuration happens through environment variables:

- `APP_MODE` (default: `client`) – `client` or `server`

- `PORT` (default: `2222`) – the TCP/UDP port to use

- `CLIENT_PROTO` (default `tcp`) – `tcp` or `udp`

- `CLIENT_TARGET_HOST` (default: `127.0.0.1`) – hostname or IP to target in
  client mode.

- `CLIENT_CONN_INTERVAL_MS` (default: `1000`) – connection interval in milliseconds

- `HANDLE_SIGNALS` (default: unset) – if set to any value, SIGINT and SIGTERM
  are caught and timing info is logged every 200 milliseconds for as long as
  possible

## Non-configuration

Some things are not yet configurable:

- Client TCP/UDP dial timeout is 800 milliseconds

- Client/server TCP/UDP read/write timeout is 800 milliseconds

- Client/server read buffer is 4 kB

## Known issues

- Go issues spurious SIGURG signals that can be a bit noisy.
