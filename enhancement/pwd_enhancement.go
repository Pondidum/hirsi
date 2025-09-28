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

func (e *PwdEnhancement) Enhance(m *message.Message) error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	m.Tags.Set(PwdTag, dir)

	return nil
}
