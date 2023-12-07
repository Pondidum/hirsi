package send

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/spf13/pflag"

	_ "github.com/mattn/go-sqlite3"
)

type SendCommand struct{}

func NewSendCommand() *SendCommand {
	return &SendCommand{}
}

func (c *SendCommand) Synopsis() string {
	return "sends a message"
}

func (c *SendCommand) Flags() *pflag.FlagSet {

	flags := pflag.NewFlagSet("send", pflag.ContinueOnError)
	return flags
}

var createTable = `
CREATE TABLE IF NOT EXISTS log(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	json TEXT NOT NULL
);
`

func (c *SendCommand) Execute(args []string) error {

	db, err := sql.Open("sqlite3", "sqlite.db")
	if err != nil {
		return err
	}

	defer db.Close()

	if _, err := db.Exec(createTable); err != nil {
		return err
	}

	json := fmt.Sprintf(`{"text": "%s"}`, strings.Join(args, " "))

	if _, err := db.Exec(`insert into log(json) values(?)`, json); err != nil {
		return err
	}

	rows, err := db.Query(`select * from log order by id`)
	if err != nil {
		return err
	}

	defer rows.Close()

	for rows.Next() {
		var e entry

		if err := rows.Scan(&e.id, &e.json); err != nil {
			return err
		}

		fmt.Printf("%v: %v\n", e.id, e.json)
	}

	return nil
}

type entry struct {
	id   int
	json string
}
