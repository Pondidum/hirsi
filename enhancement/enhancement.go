package enhancement

import "hirsi/message"

type Enhancement interface {
	Enhance(message *message.Message) error
}

type EnhancementBuilder interface {
	Build() Enhancement
}

func GetBuilder(configType string) (EnhancementBuilder, bool) {
	switch configType {
	case "pwd":
		return &PwdConfig{}, true
	case "tag-add":
		return &TagAddConfig{}, true
	}

	return nil, false
}
