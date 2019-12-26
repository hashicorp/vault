package token

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
)

var _ TokenHelper = (*InternalTokenHelper)(nil)

// InternalTokenHelper fulfills the TokenHelper interface when no external
// token-helper is configured, and avoids shelling out
type InternalTokenHelper struct {
	tokenPath string
}

// populateTokenPath figures out the token path using homedir to get the user's
// home directory
func (i *InternalTokenHelper) populateTokenPath() {
	homePath, err := homedir.Dir()
	if err != nil {
		panic(fmt.Sprintf("error getting user's home directory: %v", err))
	}
	i.tokenPath = filepath.Join(homePath, ".vault-token")
}

func (i *InternalTokenHelper) Path() string {
	return i.tokenPath
}

// Get gets the value of the stored token, if any
func (i *InternalTokenHelper) Get() (string, error) {
	i.populateTokenPath()
	f, err := os.Open(i.tokenPath)
	if os.IsNotExist(err) {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	defer f.Close()

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, f); err != nil {
		return "", err
	}

	return strings.TrimSpace(buf.String()), nil
}

// Store stores the value of the token to the file
func (i *InternalTokenHelper) Store(input string) error {
	i.populateTokenPath()
	f, err := os.OpenFile(i.tokenPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	buf := bytes.NewBufferString(input)
	if _, err := io.Copy(f, buf); err != nil {
		return err
	}

	return nil
}

// Erase erases the value of the token
func (i *InternalTokenHelper) Erase() error {
	i.populateTokenPath()
	if err := os.Remove(i.tokenPath); err != nil && !os.IsNotExist(err) {
		return err
	}

	return nil
}
