

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


## message format

```json
{
  "message": ".......",
  "tags": {
    "one": "value",
    "two": "value",
  }
}
```


## data model



### entries

- multiple writers fine
- insert only

```sql
id pk -- uuid7
written_at timestamp -- when the person wrote the line, client timestamp
stored_at timestamp -- server insert timestamp
data json -- { "text": "" }
```

### tags

- probably single writer

```sql
id uuid7 pk
entry_id fk entries.pk
tag text
value json
```

### state

- used to figure out what messages have appeared since an enricher ran

```sql
id uuid7 pk
enricher text
last_stored_seen timestamp
data json
```

