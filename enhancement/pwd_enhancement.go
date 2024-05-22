package enhancement

import (
	"hirsi/message"
	"os"
)

type PwdEnhancement struct{}

func (e *PwdEnhancement) Enhance(message *message.Message) error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	message.Tags["pwd"] = dir

	return nil
}
