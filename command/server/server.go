package server

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/spf13/pflag"
	"go.opentelemetry.io/otel"

	"github.com/gofiber/fiber/v2"

	// fiberOtel "github.com/psmarcin/fiber-opentelemetry/pkg/fiber-otel"
	"github.com/gofiber/contrib/otelfiber/v2"

	_ "github.com/mattn/go-sqlite3"
)

var tr = otel.Tracer("server")

type ServerCommand struct {
	addr string
}

func NewServerCommand() *ServerCommand {
	return &ServerCommand{
		addr: "localhost:5757",
	}
}

func (c *ServerCommand) Synopsis() string {
	return "run the server"
}

func (c *ServerCommand) Flags() *pflag.FlagSet {

	flags := pflag.NewFlagSet("server", pflag.ContinueOnError)
	return flags
}

func (c *ServerCommand) Execute(ctx context.Context, args []string) error {

	if err := c.initDB(ctx); err != nil {
		return err
	}

	app := fiber.New(fiber.Config{})

	app.Use(otelfiber.Middleware())

	app.Post("/api/messages", func(c *fiber.Ctx) error {

		ctx, span := tr.Start(c.UserContext(), "messages")
		defer span.End()

		db, err := sql.Open("sqlite3", "hirsi.db")
		if err != nil {
			return err
		}
		defer db.Close()

		dto := struct {
			Message string `json:"message"`
		}{}

		if err := c.BodyParser(&dto); err != nil {
			return err
		}

		content, err := json.Marshal(dto)
		if err != nil {
			return err
		}

		if _, err := db.ExecContext(ctx, `insert into log(timestamp, json) values(?, ?)`, time.Now(), string(content)); err != nil {
			return err
		}

		c.Status(http.StatusAccepted)
		return nil
	})

	return app.Listen(c.addr)
}

func (c *ServerCommand) initDB(ctx context.Context) error {

	db, err := sql.Open("sqlite3", "hirsi.db")
	if err != nil {
		return err
	}
	defer db.Close()

	createTable := `
CREATE TABLE IF NOT EXISTS log(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	timestamp datetime NOT NULL,
	json TEXT NOT NULL
);
`

	_, err = db.ExecContext(ctx, createTable)
	return err
}
