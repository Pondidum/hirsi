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

	pwdTag := message.Tag{Key: PwdTag, Value: dir}

	// there can be only one pwd tag.
	for i, tag := range m.Tags {
		if tag.Key != PwdTag {
			continue
		}

		m.Tags[i] = pwdTag
		return nil
	}

	m.Tags = append(m.Tags, pwdTag)

	return nil
}
