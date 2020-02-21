package credsutil

import (
	"fmt"

	"github.com/hashicorp/go-multierror"
)

const (
	DefaultPasswordLength = 15
)

type PasswordProducer struct {
	Charset string
	Length  int
}

type PasswordOpt func(*PasswordProducer) error

func PasswordCharset(charset string) PasswordOpt {
	return func(p *PasswordProducer) error {
		p.Charset = charset
		return nil
	}
}

func PasswordLength(length int) PasswordOpt {
	return func(p *PasswordProducer) error {
		if length < 10 {
			return fmt.Errorf("length must be >= 10")
		}
		p.Length = length
		return nil
	}
}

func NewPasswordProducer(opts ...PasswordOpt) (p PasswordProducer, err error) {
	merr := &multierror.Error{}
	for _, opt := range opts {
		merr = multierror.Append(merr, opt(&p))
	}

	return p, merr.ErrorOrNil()
}

func (pg PasswordProducer) GeneratePassword() (password string, err error) {
	charset := pg.Charset
	if charset == "" {
		charset = DefaultCharset
	}

	length := pg.Length
	if length == 0 {
		length = DefaultPasswordLength
	}

	return randomChars(charset, length)
}
