package storage

import (
	"context"
	"database/sql"
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

	for key, value := range m.Tags {

		_, err := tx.ExecContext(ctx,
			"insert into tags (log_id, key, value) values (?, ?, ?)",
			logId, key, value,
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
