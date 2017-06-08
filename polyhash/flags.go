package polyhash

import (
	"errors"
	"fmt"
	"strings"
)

type PolyhashPasswordsFlag []string

func (p *PolyhashPasswordsFlag) String() string {
	return fmt.Sprint(*p)
}

func (p *PolyhashPasswordsFlag) Set(value string) error {
	if len(*p) > 0 {
		return errors.New("polyhash can only be specified once")
	}

	splitValues := strings.Split(value, ",")
	for _, password := range splitValues {
		*p = append(*p, password)
	}

	return nil
}
