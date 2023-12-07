

## Commands

- `send` write a message
- `watch` see what messages are coming in
- `server` recieve and do things with messages


## Architecture

- client sends messages to `server`
- `server`
  - recieves http messages and writes them to a store
  - by default, after writing to store, will then run any processors
  - litestream replicates datastore to s3
- processor
  - litestream pulling "read" copies of the datastore
  - uses the `watch` functionality to process messages its interested in
  - writes are sent back to `server` (single writer model)