package enhancement

import (
	"hirsi/message"
	"os"
)

const PwdTag = "pwd"

type PwdConfig struct{}

func (ec *PwdConfig) Build() Enhancement {
	return &PwdEnhancement{}
}

type PwdEnhancement struct{}

func (e *PwdEnhancement) Enhance(message *message.Message) error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	message.Tags[PwdTag] = dir

	return nil
}
