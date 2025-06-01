package helper

import "fmt"

type Mode int

const (
	Development Mode = iota
	Production
)

var modeMap = map[string]Mode{
	"development": Development,
	"production":  Production,
}

func ParseMode(s string) (Mode, error) {
	val, ok := modeMap[s]
	if ok {
		return val, nil
	}
	return 0, fmt.Errorf("invalid mode: %s\n", s)
}
