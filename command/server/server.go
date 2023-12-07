package server

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/spf13/pflag"

	"github.com/gofiber/fiber/v2"

	_ "github.com/mattn/go-sqlite3"
)

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

func (c *ServerCommand) Execute(args []string) error {

	if err := c.initDB(); err != nil {
		return err
	}

	app := fiber.New(fiber.Config{})

	app.Use(func(c *fiber.Ctx) error {
		err := c.Next()
		fmt.Printf("%v %s %s \n", c.Response().StatusCode(), c.Method(), c.Path())
		return err
	})

	app.Post("/api/messages", func(c *fiber.Ctx) error {

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

		if _, err := db.Exec(`insert into log(timestamp, json) values(?, ?)`, time.Now(), string(content)); err != nil {
			return err
		}

		c.Status(http.StatusAccepted)
		return nil
	})

	return app.Listen(c.addr)
}

func (c *ServerCommand) initDB() error {

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

	_, err = db.Exec(createTable)
	return err
}
