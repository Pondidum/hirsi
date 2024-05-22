package ls

import (
	"context"
	"database/sql"
	"fmt"
	"hirsi/config"
	"strings"

	"github.com/ryanuber/columnize"
	"github.com/spf13/pflag"

	_ "github.com/mattn/go-sqlite3"
)

type LsCommand struct {
	config *config.Config
}

func NewLsCommand(config *config.Config) *LsCommand {
	return &LsCommand{config: config}
}

func (c *LsCommand) Synopsis() string {
	return "view the logs"
}

func (c *LsCommand) Flags() *pflag.FlagSet {
	flags := pflag.NewFlagSet("ls", pflag.ContinueOnError)
	return flags
}

func (c *LsCommand) Execute(ctx context.Context, args []string) error {

	db, err := sql.Open("sqlite3", c.config.DbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	query := `
select log.id, log.stored_at, log.written_at, log.message, json_group_object(tags.key, tags.value) tags
from log inner join tags on log.id = tags.log_id
group by log.id, log.stored_at, log.written_at, log.message
`

	rows, err := db.Query(query)
	if err != nil {
		return err
	}

	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return err
	}

	output := []string{
		strings.Join(cols, " | "),
	}

	for rows.Next() {

		values := make([]any, len(cols))
		for i := range values {
			values[i] = &RowWriter{}
		}

		if err := rows.Scan(values...); err != nil {
			return err
		}

		outputRow := make([]string, len(cols))
		for i, v := range values {
			outputRow[i] = fmt.Sprint(v)
		}

		output = append(output, strings.Join(outputRow, " | "))
	}

	fmt.Println(tableOutput(output))
	return nil
}

type RowWriter struct {
	value any
}

func (r *RowWriter) Scan(src any) error {
	r.value = src
	return nil
}

func (r *RowWriter) String() string {
	return fmt.Sprint(r.value)
}

func tableOutput(list []string) string {
	if len(list) == 0 {
		return ""
	}

	delim := "|"
	underline := ""
	headers := strings.Split(list[0], delim)
	for i, h := range headers {
		h = strings.TrimSpace(h)
		u := strings.Repeat("-", len(h))

		underline = underline + u
		if i != len(headers)-1 {
			underline = underline + delim
		}
	}

	list = append(list, "")
	copy(list[2:], list[1:])
	list[1] = underline

	return columnize.Format(list, &columnize.Config{
		Glue: "    ",
	})
}
