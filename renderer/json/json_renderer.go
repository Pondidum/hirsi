package json

import (
	"encoding/json"
	"hirsi/message"
	"io"
)

type JsonRenderer struct {
	Writer io.Writer
}

func (r *JsonRenderer) Render(message *message.Message) error {
	b, err := json.MarshalIndent(message, "", "  ")
	if err != nil {
		return err
	}

	if _, err := r.Writer.Write(b); err != nil {
		return err
	}

	if _, err := r.Writer.Write([]byte("\n")); err != nil {
		return err
	}

	return nil
}
