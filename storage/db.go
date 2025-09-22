package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"hirsi/message"
	"hirsi/tracing"
	"os"
	"path"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

var tr = otel.Tracer("storage")

func MigrateDatabase(ctx context.Context, dbPath string) error {

	dir := path.Dir(dbPath)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	createTables := `
CREATE TABLE IF NOT EXISTS log(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	written_at datetime NOT NULL,
	stored_at datetime NOT NULL,
	message TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS tags(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	log_id INTEGER,
	key TEXT,
	value TEXT,

	FOREIGN KEY(log_id) REFERENCES log(id)
);
`

	_, err = db.ExecContext(ctx, createTables)
	return err
}

func StoreMessage(ctx context.Context, dbPath string, m *message.Message) error {
	ctx, span := tr.Start(ctx, "store_message")
	defer span.End()

	span.SetAttributes(attribute.String("db_path", dbPath))

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	tx, err := db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return tracing.Error(span, err)
	}
	defer tx.Rollback()

	result, err := tx.ExecContext(ctx,
		"insert into log (written_at, stored_at, message) values (?,?,?)",
		m.WrittenAt, time.Now(), m.Message,
	)
	if err != nil {
		return tracing.Error(span, err)
	}

	logId, err := result.LastInsertId()
	if err != nil {
		return tracing.Error(span, err)
	}

	for _, tag := range m.Tags {

		_, err := tx.ExecContext(ctx,
			"insert into tags (log_id, key, value) values (?, ?, ?)",
			logId, tag.Key, tag.Value,
		)
		if err != nil {
			return tracing.Error(span, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return tracing.Error(span, err)
	}

	return nil
}

func EachMessage(ctx context.Context, dbPath string, recentCount int, onMessage func(m *message.Message) error) error {
	ctx, span := tr.Start(ctx, "list_messages")
	defer span.End()

	span.SetAttributes(attribute.String("db_path", dbPath))

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	query := `
select 
	log.id,
	log.stored_at,
	log.written_at,
	log.message,
	json_group_array(json_object('key', tags.key, 'value', tags.value)) tags
from log
inner join tags on log.id = tags.log_id
group by log.id, log.stored_at, log.written_at, log.message
order by log.stored_at desc
`

	parameters := []any{}

	if recentCount > 0 {
		query += `limit ?`
		parameters = append(parameters, recentCount)
	}

	//"select * from log sort by stored_at desc limit ?"
	rows, err := db.QueryContext(ctx, query, parameters...)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		id := 0
		m := &message.Message{}

		tag := &TagReader{Target: []message.Tag{}}

		if err := rows.Scan(&id, &m.StoredAt, &m.WrittenAt, &m.Message, &tag); err != nil {
			return err
		}

		m.Tags = tag.Target

		if err := onMessage(m); err != nil {
			return err
		}
	}

	return nil
}

func ListMessages(ctx context.Context, dbPath string, recentCount int) ([]*message.Message, error) {

	messages := make([]*message.Message, 0, recentCount)

	err := EachMessage(ctx, dbPath, recentCount, func(m *message.Message) error {
		messages = append(messages, m)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return messages, nil
}

type TagReader struct {
	Target []message.Tag
}

func (r *TagReader) Scan(src any) error {
	return json.Unmarshal([]byte(fmt.Sprint(src)), &r.Target)
}
