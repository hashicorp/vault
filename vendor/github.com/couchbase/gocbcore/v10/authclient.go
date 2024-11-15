package gocbcore

import (
	"crypto/sha1" // nolint: gosec
	"crypto/sha256"
	"crypto/sha512"
	"hash"
	"time"

	"github.com/couchbase/gocbcore/v10/memd"

	scram "github.com/couchbase/gocbcore/v10/scram"
)

// AuthMechanism represents a type of auth that can be performed.
type AuthMechanism string

const (
	// PlainAuthMechanism represents that PLAIN auth should be performed.
	PlainAuthMechanism = AuthMechanism("PLAIN")

	// ScramSha1AuthMechanism represents that SCRAM SHA1 auth should be performed.
	ScramSha1AuthMechanism = AuthMechanism("SCRAM-SHA1")

	// ScramSha256AuthMechanism represents that SCRAM SHA256 auth should be performed.
	ScramSha256AuthMechanism = AuthMechanism("SCRAM-SHA256")

	// ScramSha512AuthMechanism represents that SCRAM SHA512 auth should be performed.
	ScramSha512AuthMechanism = AuthMechanism("SCRAM-SHA512")
)

// AuthClient exposes an interface for performing authentication on a
// connected Couchbase K/V client.
type AuthClient interface {
	Address() string
	SupportsFeature(feature memd.HelloFeature) bool

	SaslListMechs(deadline time.Time, cb func(mechs []AuthMechanism, err error)) error
	SaslAuth(k, v []byte, deadline time.Time, cb func(b []byte, err error)) error
	SaslStep(k, v []byte, deadline time.Time, cb func(err error)) error
}

// SaslListMechsCompleted is used to contain the result and/or error from a SaslListMechs operation.
type SaslListMechsCompleted struct {
	Err   error
	Mechs []AuthMechanism
}

// SaslAuthPlain performs PLAIN SASL authentication against an AuthClient.
func SaslAuthPlain(username, password string, client AuthClient, deadline time.Time, cb func(err error)) error {
	// Build PLAIN auth data
	userBuf := []byte(username)
	passBuf := []byte(password)
	authData := make([]byte, 1+len(userBuf)+1+len(passBuf))
	authData[0] = 0
	copy(authData[1:], userBuf)
	authData[1+len(userBuf)] = 0
	copy(authData[1+len(userBuf)+1:], passBuf)

	// Execute PLAIN authentication
	err := client.SaslAuth([]byte(PlainAuthMechanism), authData, deadline, func(b []byte, err error) {
		if err != nil {
			cb(err)
			return
		}
		cb(nil)
	})
	if err != nil {
		return err
	}

	return nil
}

func saslAuthScram(saslName []byte, newHash func() hash.Hash, username, password string, client AuthClient,
	deadline time.Time, continueCb func(), completedCb func(err error)) error {
	scramMgr := scram.NewClient(newHash, username, password)

	// Perform the initial SASL step
	scramMgr.Step(nil)
	err := client.SaslAuth(saslName, scramMgr.Out(), deadline, func(b []byte, err error) {
		if err != nil && !isErrorStatus(err, memd.StatusAuthContinue) {
			completedCb(err)
			return
		}

		if !scramMgr.Step(b) {
			err = scramMgr.Err()
			if err != nil {
				completedCb(err)
				return
			}

			logErrorf("Local auth client finished before server accepted auth")
			completedCb(nil)
			return
		}

		err = client.SaslStep(saslName, scramMgr.Out(), deadline, completedCb)
		if err != nil {
			completedCb(err)
			return
		}

		continueCb()
	})
	if err != nil {
		return err
	}

	return nil
}

// SaslAuthScramSha1 performs SCRAM-SHA1 SASL authentication against an AuthClient.
func SaslAuthScramSha1(username, password string, client AuthClient, deadline time.Time, continueCb func(), completedCb func(err error)) error {
	return saslAuthScram([]byte("SCRAM-SHA1"), sha1.New, username, password, client, deadline, continueCb, completedCb)
}

// SaslAuthScramSha256 performs SCRAM-SHA256 SASL authentication against an AuthClient.
func SaslAuthScramSha256(username, password string, client AuthClient, deadline time.Time, continueCb func(), completedCb func(err error)) error {
	return saslAuthScram([]byte("SCRAM-SHA256"), sha256.New, username, password, client, deadline, continueCb, completedCb)
}

// SaslAuthScramSha512 performs SCRAM-SHA512 SASL authentication against an AuthClient.
func SaslAuthScramSha512(username, password string, client AuthClient, deadline time.Time, continueCb func(), completedCb func(err error)) error {
	return saslAuthScram([]byte("SCRAM-SHA512"), sha512.New, username, password, client, deadline, continueCb, completedCb)
}

func saslMethod(method AuthMechanism, username, password string, client AuthClient, deadline time.Time, continueCb func(), completedCb func(err error)) error {
	switch method {
	case PlainAuthMechanism:
		return SaslAuthPlain(username, password, client, deadline, completedCb)
	case ScramSha1AuthMechanism:
		return SaslAuthScramSha1(username, password, client, deadline, continueCb, completedCb)
	case ScramSha256AuthMechanism:
		return SaslAuthScramSha256(username, password, client, deadline, continueCb, completedCb)
	case ScramSha512AuthMechanism:
		return SaslAuthScramSha512(username, password, client, deadline, continueCb, completedCb)
	default:
		return errNoSupportedMechanisms
	}
}
