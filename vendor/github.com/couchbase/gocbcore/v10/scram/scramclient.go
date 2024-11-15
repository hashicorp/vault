// Copyright (c) 2014 - Gustavo Niemeyer <gustavo@niemeyer.net>
// Copyright (c) 2017 - Couchbase Inc.
//
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
// 1. Redistributions of source code must retain the above copyright notice, this
//    list of conditions and the following disclaimer.
// 2. Redistributions in binary form must reproduce the above copyright notice,
//    this list of conditions and the following disclaimer in the documentation
//    and/or other materials provided with the distribution.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
// ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
// WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE LIABLE FOR
// ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
// (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
// LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
// ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
// SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package gocbcore

import (
	"bytes"
	"crypto/hmac"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"hash"
	"strconv"
	"strings"
)

// Client implements a SCRAM-{SHA-1,etc} client per RFC5802.
// http://tools.ietf.org/html/rfc5802
type Client struct {
	newHash func() hash.Hash

	user string
	pass string
	step int
	out  bytes.Buffer
	err  error

	clientNonce []byte
	serverNonce []byte
	saltedPass  []byte
	authMsg     bytes.Buffer
}

// NewClient returns a new instance of the SCRAM client.
func NewClient(newHash func() hash.Hash, user, pass string) *Client {
	c := &Client{
		newHash: newHash,
		user:    user,
		pass:    pass,
	}
	c.out.Grow(256)
	c.authMsg.Grow(256)
	return c
}

// Out returns the data to be sent to the server in the current step.
func (c *Client) Out() []byte {
	if c.out.Len() == 0 {
		return nil
	}
	return c.out.Bytes()
}

// Err returns the error that occurred, or nil if there were no errors.
func (c *Client) Err() error {
	return c.err
}

// SetNonce sets the client nonce to the provided value.
// If not set, the nonce is generated automatically out of crypto/rand on the first step.
func (c *Client) SetNonce(nonce []byte) {
	c.clientNonce = nonce
}

var escaper = strings.NewReplacer("=", "=3D", ",", "=2C")

// Step processes the incoming data from the server and makes the
// next round of data for the server available via Client.Out.
// Step returns false if there are no errors and more data is
// still expected.
func (c *Client) Step(in []byte) bool {
	c.out.Reset()
	if c.step > 2 || c.err != nil {
		return false
	}
	c.step++
	switch c.step {
	case 1:
		c.err = c.step1(in)
	case 2:
		c.err = c.step2(in)
	case 3:
		c.err = c.step3(in)
	}

	return !(c.step > 2 || c.err != nil)
}

func (c *Client) step1(in []byte) error {
	if len(c.clientNonce) == 0 {
		const nonceLen = 6
		buf := make([]byte, nonceLen+b64.EncodedLen(nonceLen))
		if _, err := rand.Read(buf[:nonceLen]); err != nil {
			return fmt.Errorf("cannot read random SCRAM-SHA-1 nonce from operating system: %v", err)
		}
		c.clientNonce = buf[nonceLen:]
		b64.Encode(c.clientNonce, buf[:nonceLen])
	}
	c.authMsg.WriteString("n=")
	if _, err := escaper.WriteString(&c.authMsg, c.user); err != nil {
		return err
	}
	c.authMsg.WriteString(",r=")
	c.authMsg.Write(c.clientNonce)

	c.out.WriteString("n,,")
	c.out.Write(c.authMsg.Bytes())
	return nil
}

var b64 = base64.StdEncoding

func (c *Client) step2(in []byte) error {
	c.authMsg.WriteByte(',')
	c.authMsg.Write(in)

	fields := bytes.Split(in, []byte(","))
	if len(fields) != 3 {
		return fmt.Errorf("expected 3 fields in first SCRAM-SHA-1 server message, got %d: %q", len(fields), in)
	}
	if !bytes.HasPrefix(fields[0], []byte("r=")) {
		return fmt.Errorf("server sent an invalid SCRAM-SHA-1 nonce: %q", fields[0])
	}
	if !bytes.HasPrefix(fields[1], []byte("s=")) {
		return fmt.Errorf("server sent an invalid SCRAM-SHA-1 salt: %q", fields[1])
	}
	if !bytes.HasPrefix(fields[2], []byte("i=")) {
		return fmt.Errorf("server sent an invalid SCRAM-SHA-1 iteration count: %q", fields[2])
	}

	c.serverNonce = fields[0][2:]
	if !bytes.HasPrefix(c.serverNonce, c.clientNonce) {
		return fmt.Errorf("server SCRAM-SHA-1 nonce is not prefixed by client nonce: got %q, want %q+\"...\"", c.serverNonce, c.clientNonce)
	}

	salt := make([]byte, b64.DecodedLen(len(fields[1][2:])))
	n, err := b64.Decode(salt, fields[1][2:])
	if err != nil {
		return fmt.Errorf("cannot decode SCRAM-SHA-1 salt sent by server: %q", fields[1])
	}
	salt = salt[:n]
	iterCount, err := strconv.Atoi(string(fields[2][2:]))
	if err != nil {
		return fmt.Errorf("server sent an invalid SCRAM-SHA-1 iteration count: %q", fields[2])
	}
	if err := c.saltPassword(salt, iterCount); err != nil {
		return err
	}

	c.authMsg.WriteString(",c=biws,r=")
	c.authMsg.Write(c.serverNonce)

	c.out.WriteString("c=biws,r=")
	c.out.Write(c.serverNonce)
	c.out.WriteString(",p=")
	proof, err := c.clientProof()
	if err != nil {
		return err
	}
	c.out.Write(proof)
	return nil
}

func (c *Client) step3(in []byte) error {
	var isv, ise bool
	var fields = bytes.Split(in, []byte(","))
	if len(fields) == 1 {
		isv = bytes.HasPrefix(fields[0], []byte("v="))
		ise = bytes.HasPrefix(fields[0], []byte("e="))
	}
	if ise {
		return fmt.Errorf("SCRAM-SHA-1 authentication error: %s", fields[0][2:])
	} else if !isv {
		return fmt.Errorf("unsupported SCRAM-SHA-1 final message from server: %q", in)
	}

	sig, err := c.serverSignature()
	if err != nil {
		return err
	}

	if !bytes.Equal(sig, fields[0][2:]) {
		return fmt.Errorf("cannot authenticate SCRAM-SHA-1 server signature: %q", fields[0][2:])
	}
	return nil
}

func (c *Client) saltPassword(salt []byte, iterCount int) error {
	mac := hmac.New(c.newHash, []byte(c.pass))
	if _, err := mac.Write(salt); err != nil {
		return err
	}
	if _, err := mac.Write([]byte{0, 0, 0, 1}); err != nil {
		return err
	}
	ui := mac.Sum(nil)
	hi := make([]byte, len(ui))
	copy(hi, ui)
	for i := 1; i < iterCount; i++ {
		mac.Reset()
		if _, err := mac.Write(ui); err != nil {
			return err
		}
		mac.Sum(ui[:0])
		for j, b := range ui {
			hi[j] ^= b
		}
	}
	c.saltedPass = hi
	return nil
}

func (c *Client) clientProof() ([]byte, error) {
	mac := hmac.New(c.newHash, c.saltedPass)
	if _, err := mac.Write([]byte("Client Key")); err != nil {
		return nil, err
	}
	clientKey := mac.Sum(nil)
	hash := c.newHash()
	if _, err := hash.Write(clientKey); err != nil {
		return nil, err
	}
	storedKey := hash.Sum(nil)
	mac = hmac.New(c.newHash, storedKey)
	if _, err := mac.Write(c.authMsg.Bytes()); err != nil {
		return nil, err
	}
	clientProof := mac.Sum(nil)
	for i, b := range clientKey {
		clientProof[i] ^= b
	}
	clientProof64 := make([]byte, b64.EncodedLen(len(clientProof)))
	b64.Encode(clientProof64, clientProof)
	return clientProof64, nil
}

func (c *Client) serverSignature() ([]byte, error) {
	mac := hmac.New(c.newHash, c.saltedPass)
	if _, err := mac.Write([]byte("Server Key")); err != nil {
		return nil, err
	}
	serverKey := mac.Sum(nil)

	mac = hmac.New(c.newHash, serverKey)
	if _, err := mac.Write(c.authMsg.Bytes()); err != nil {
		return nil, err
	}
	serverSignature := mac.Sum(nil)

	encoded := make([]byte, b64.EncodedLen(len(serverSignature)))
	b64.Encode(encoded, serverSignature)
	return encoded, nil
}
