package credsutil

import (
	"fmt"

	"github.com/hashicorp/go-multierror"
)

const (
	MinPasswordLength     = 10
	DefaultPasswordLength = 15

	// DefaultPasswordPrefix to prepend to the password. This is here to adhere to minimum password requirements
	// that the password may be used in.
	DefaultPasswordPrefix = "A1a-"
)

type PasswordProducer struct {
	Charset string
	Length  int
	Prefix  string
}

type PasswordOpt func(*PasswordProducer) error

// PasswordCharset specifies the charset to be used when generating the random password.
func PasswordCharset(charset string) PasswordOpt {
	return func(p *PasswordProducer) error {
		p.Charset = charset
		return nil
	}
}

// PasswordLength to be generated. This includes the length of the prefix, if specified.
func PasswordLength(length int) PasswordOpt {
	return func(p *PasswordProducer) error {
		if length < MinPasswordLength {
			return fmt.Errorf("length must be >= %d", MinPasswordLength)
		}
		p.Length = length
		return nil
	}
}

// PasswordPrefix specifies the prefix to prepend to the generated password.
func PasswordPrefix(prefix string) PasswordOpt {
	return func(p *PasswordProducer) error {
		p.Prefix = prefix
		return nil
	}
}

// NewPasswordProducer for generating passwords. This has reasonable defaults so it can be used without any arguments.
func NewPasswordProducer(opts ...PasswordOpt) (p PasswordProducer, err error) {
	p = PasswordProducer{
		Charset: DefaultCharset,
		Length:  DefaultPasswordLength,
		Prefix:  DefaultPasswordPrefix,
	}
	merr := &multierror.Error{}
	for _, opt := range opts {
		merr = multierror.Append(merr, opt(&p))
	}

	return p, merr.ErrorOrNil()
}

// GeneratePassword using the values specified. This adheres to the CredentialsProducer interface.
func (pg PasswordProducer) GeneratePassword() (password string, err error) {
	charset := pg.Charset
	if charset == "" {
		charset = DefaultCharset
	}

	length := pg.Length
	if length == 0 {
		length = DefaultPasswordLength
	}
	if length < MinPasswordLength {
		return "", fmt.Errorf("length must be >= %d", MinPasswordLength)
	}

	length = length - len(pg.Prefix)

	chars, err := randomChars(charset, length)
	if err != nil {
		return "", fmt.Errorf("unable to generate password: %w", err)
	}
	return pg.Prefix + chars, nil
}
