package send

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/spf13/pflag"

	_ "github.com/mattn/go-sqlite3"
)

type SendCommand struct {
	addr string
}

func NewSendCommand() *SendCommand {
	return &SendCommand{
		addr: "http://localhost:5757",
	}
}

func (c *SendCommand) Synopsis() string {
	return "sends a message"
}

func (c *SendCommand) Flags() *pflag.FlagSet {

	flags := pflag.NewFlagSet("send", pflag.ContinueOnError)
	return flags
}

func (c *SendCommand) Execute(args []string) error {

	values := url.Values{}
	values.Add("message", strings.Join(args, " "))

	res, err := http.PostForm(c.addr+"/api/messages", values)
	if err != nil {
		return err
	}

	if res.StatusCode < 200 && res.StatusCode > 299 {
		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}

		return fmt.Errorf("error from server:\n " + string(body))
	}

	return nil
}
