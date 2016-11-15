package token

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"strings"

	"net/url"

	"github.com/mitchellh/go-homedir"
)

// InternalTokenHelper fulfills the TokenHelper interface when no external
// token-helper is configured, and avoids shelling out
type InternalTokenHelper struct {
	defaultTokenPath string
	hashedTokenPath  string
	addr             string
}

func (i *InternalTokenHelper) hashAddress() (string, error) {
	// Ignore the protocol field; just use the address and port
	u, err := url.Parse(i.addr)
	if err != nil {
		return "", err
	}
	hashedAddr := sha256.Sum256([]byte(u.Host))
	asciiHashedAddr := fmt.Sprintf("%x", hashedAddr)
	return asciiHashedAddr, nil
}

// populateTokenPaths figures out the token paths using homedir to get the user's
// home directory
func (i *InternalTokenHelper) populateTokenPaths() error {
	homePath, err := homedir.Dir()
	if err != nil {
		panic(fmt.Errorf("error getting user's home directory: %v", err))
	}
	i.defaultTokenPath = homePath + "/.vault-token"
	if len(i.addr) > 0 {
		asciiHashedAddr, err := i.hashAddress()
		if err != nil {
			return err
		}
		i.hashedTokenPath = i.defaultTokenPath + "-" + asciiHashedAddr[0:8]
	}
	return nil
}

func (i *InternalTokenHelper) Path() string {
	return i.hashedTokenPath
}

// SetVaultAddress sets the Vault address in use.  This is used to determine the token path.
func (i *InternalTokenHelper) SetVaultAddress(addr string) {
	i.addr = addr
}

// Get gets the value of the stored token, if any
func (i *InternalTokenHelper) Get() (string, error) {
	var f *os.File
	var err error
	if err = i.populateTokenPaths(); err != nil {
		return "", err
	}
	if len(i.hashedTokenPath) > 0 {
		f, err = os.Open(i.hashedTokenPath)
		if err != nil && !os.IsNotExist(err) {
			return "", err
		}
	}
	if f == nil { // No token exists at hashed path
		f, err = os.Open(i.defaultTokenPath)
		if os.IsNotExist(err) {
			return "", nil
		}
		if err != nil {
			return "", err
		}
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
	if err := i.populateTokenPaths(); err != nil {
		return err
	}
	var path string
	if len(i.hashedTokenPath) > 0 {
		path = i.hashedTokenPath
	} else {
		path = i.defaultTokenPath
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
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
	i.populateTokenPaths()
	if len(i.hashedTokenPath) > 0 {
		if err := os.Remove(i.hashedTokenPath); err != nil && !os.IsNotExist(err) {
			return err
		}
	}
	if err := os.Remove(i.defaultTokenPath); err != nil && !os.IsNotExist(err) {
		return err
	}

	return nil
}
