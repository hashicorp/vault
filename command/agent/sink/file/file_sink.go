package file

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/hashicorp/errwrap"
	hclog "github.com/hashicorp/go-hclog"
	uuid "github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/command/agent/sink"
	"github.com/hashicorp/vault/helper/dhutil"
	"github.com/hashicorp/vault/helper/jsonutil"
)

// fileSink is a Sink implementation that writes a token to a file
type fileSink struct {
	path         string
	logger       hclog.Logger
	wrapTTL      time.Duration
	dhType       string
	dhPath       string
	aad          string
	cachedPubKey []byte
}

// NewFileSink creates a new file sink with the given configuration
func NewFileSink(conf *sink.SinkConfig) (sink.Sink, error) {
	if conf.Logger == nil {
		return nil, errors.New("nil logger provided")
	}

	conf.Logger.Info("creating file sink")

	f := &fileSink{
		logger:  conf.Logger,
		wrapTTL: conf.WrapTTL,
		dhType:  conf.DHType,
		dhPath:  conf.DHPath,
		aad:     conf.AAD,
	}

	pathRaw, ok := conf.Config["path"]
	if !ok {
		return nil, errors.New("'path' not specified for file sink")
	}
	path, ok := pathRaw.(string)
	if !ok {
		return nil, errors.New("could not parse 'path' as string")
	}

	f.path = path

	if err := f.WriteToken(""); err != nil {
		return nil, errwrap.Wrapf("error during write check: {{err}}", err)
	}

	f.logger.Info("file sink configured", "path", f.path)

	return f, nil
}

// WriteToken implements the Server interface and writes the token to a path on
// disk. It writes into the path's directory into a temp file and does an
// atomic rename to ensure consistency. If a blank token is passed in, it
// performs a write check but does not write a blank value to the final
// location.
func (f *fileSink) WriteToken(token string) error {
	f.logger.Trace("enter write_token", "path", f.path)
	defer f.logger.Trace("exit write_token", "path", f.path)

	u, err := uuid.GenerateUUID()
	if err != nil {
		return errwrap.Wrapf("error generating a uuid during write check: {{err}}", err)
	}

	if token != "" && f.dhType != "" {
		var aesKey []byte
		resp := new(dhutil.Envelope)
		switch f.dhType {
		case "curve25519":
			_, err = os.Lstat(f.dhPath)
			if err != nil {
				if !os.IsNotExist(err) {
					return errwrap.Wrapf("error stat-ing dh parameters file: {{err}}", err)
				}
				if len(f.cachedPubKey) > 0 {
					break
				}
				return errors.New("no dh parameters file found, and no cached pub key")
			}
			fileBytes, err := ioutil.ReadFile(f.dhPath)
			if err != nil {
				return errwrap.Wrapf("error reading file for dh parameters: {{err}}", err)
			}
			theirPubKey := new(dhutil.PublicKeyInfo)
			if err := jsonutil.DecodeJSON(fileBytes, theirPubKey); err != nil {
				// Might just be a token, so ignore if we have a cached key
				if len(f.cachedPubKey) > 0 {
					f.logger.Debug("error decoding public key, may be a token value", "error", err)
					break
				}
				return errwrap.Wrapf("error decoding public key: {{err}}", err)
			}
			if len(theirPubKey.Curve25519PublicKey) == 0 {
				return errors.New("public key is nil")
			}
			f.cachedPubKey = theirPubKey.Curve25519PublicKey
			pub, pri, err := dhutil.GeneratePublicPrivateKey()
			if err != nil {
				return errwrap.Wrapf("error generating pub/pri curve25519 keys: {{err}}", err)
			}
			aesKey, err = dhutil.GenerateSharedKey(pri, theirPubKey.Curve25519PublicKey)
			if err != nil {
				return errwrap.Wrapf("error deriving shared key: {{err}}", err)
			}
			resp.Curve25519PublicKey = pub
		}
		if len(aesKey) == 0 {
			return errors.New("derived AES key is empty")
		}
		resp.EncryptedPayload, resp.Nonce, err = dhutil.EncryptAES(aesKey, []byte(token), []byte(f.aad))
		if err != nil {
			return errwrap.Wrapf("error encrypting with shared key: {{err}}", err)
		}
		m, err := jsonutil.EncodeJSON(resp)
		if err != nil {
			return errwrap.Wrapf("error encoding encrypted payload: {{err}}", err)
		}
		token = string(m)
	}

	targetDir := filepath.Dir(f.path)
	fileName := filepath.Base(f.path)
	tmpSuffix := strings.Split(u, "-")[0]

	tmpFile, err := os.OpenFile(filepath.Join(targetDir, fmt.Sprintf("%s.tmp.%s", fileName, tmpSuffix)), os.O_WRONLY|os.O_CREATE, 0640)
	if err != nil {
		return errwrap.Wrapf(fmt.Sprintf("error opening temp file in dir %s for writing: {{err}}", targetDir), err)
	}

	valToWrite := token
	if token == "" {
		valToWrite = u
	}

	_, err = tmpFile.WriteString(valToWrite)
	if err != nil {
		// Attempt closing and deleting but ignore any error
		tmpFile.Close()
		os.Remove(tmpFile.Name())
		return errwrap.Wrapf(fmt.Sprintf("error writing to %s: {{err}}", tmpFile.Name()), err)
	}

	err = tmpFile.Close()
	if err != nil {
		return errwrap.Wrapf(fmt.Sprintf("error closing %s: {{err}}", tmpFile.Name()), err)
	}

	// Now, if we were just doing a write check (blank token), remove the file
	// and exit; otherwise, atomically rename and chmod it
	if token == "" {
		err = os.Remove(tmpFile.Name())
		if err != nil {
			return errwrap.Wrapf(fmt.Sprintf("error removing temp file %s during write check: {{err}}", tmpFile.Name()), err)
		}
		return nil
	}

	err = os.Rename(tmpFile.Name(), f.path)
	if err != nil {
		return errwrap.Wrapf(fmt.Sprintf("error renaming temp file %s to target file %s: {{err}}", tmpFile.Name(), f.path), err)
	}

	f.logger.Info("token written", "path", f.path)
	return nil
}

func (f *fileSink) WrapTTL() time.Duration {
	return f.wrapTTL
}
