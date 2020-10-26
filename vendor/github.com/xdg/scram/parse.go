// Copyright 2018 by David A. Golden. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package scram

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type c1Msg struct {
	gs2Header string
	authzID   string
	username  string
	nonce     string
	c1b       string
}

type c2Msg struct {
	cbind []byte
	nonce string
	proof []byte
	c2wop string
}

type s1Msg struct {
	nonce string
	salt  []byte
	iters int
}

type s2Msg struct {
	verifier []byte
	err      string
}

func parseField(s, k string) (string, error) {
	t := strings.TrimPrefix(s, k+"=")
	if t == s {
		return "", fmt.Errorf("error parsing '%s' for field '%s'", s, k)
	}
	return t, nil
}

func parseGS2Flag(s string) (string, error) {
	if s[0] == 'p' {
		return "", fmt.Errorf("channel binding requested but not supported")
	}

	if s == "n" || s == "y" {
		return s, nil
	}

	return "", fmt.Errorf("error parsing '%s' for gs2 flag", s)
}

func parseFieldBase64(s, k string) ([]byte, error) {
	raw, err := parseField(s, k)
	if err != nil {
		return nil, err
	}

	dec, err := base64.StdEncoding.DecodeString(raw)
	if err != nil {
		return nil, err
	}

	return dec, nil
}

func parseFieldInt(s, k string) (int, error) {
	raw, err := parseField(s, k)
	if err != nil {
		return 0, err
	}

	num, err := strconv.Atoi(raw)
	if err != nil {
		return 0, fmt.Errorf("error parsing field '%s': %v", k, err)
	}

	return num, nil
}

func parseClientFirst(c1 string) (msg c1Msg, err error) {

	fields := strings.Split(c1, ",")
	if len(fields) < 4 {
		err = errors.New("not enough fields in first server message")
		return
	}

	gs2flag, err := parseGS2Flag(fields[0])
	if err != nil {
		return
	}

	// 'a' field is optional
	if len(fields[1]) > 0 {
		msg.authzID, err = parseField(fields[1], "a")
		if err != nil {
			return
		}
	}

	// Recombine and save the gs2 header
	msg.gs2Header = gs2flag + "," + msg.authzID + ","

	// Check for unsupported extensions field "m".
	if strings.HasPrefix(fields[2], "m=") {
		err = errors.New("SCRAM message extensions are not supported")
		return
	}

	msg.username, err = parseField(fields[2], "n")
	if err != nil {
		return
	}

	msg.nonce, err = parseField(fields[3], "r")
	if err != nil {
		return
	}

	msg.c1b = strings.Join(fields[2:], ",")

	return
}

func parseClientFinal(c2 string) (msg c2Msg, err error) {
	fields := strings.Split(c2, ",")
	if len(fields) < 3 {
		err = errors.New("not enough fields in first server message")
		return
	}

	msg.cbind, err = parseFieldBase64(fields[0], "c")
	if err != nil {
		return
	}

	msg.nonce, err = parseField(fields[1], "r")
	if err != nil {
		return
	}

	// Extension fields may come between nonce and proof, so we
	// grab the *last* fields as proof.
	msg.proof, err = parseFieldBase64(fields[len(fields)-1], "p")
	if err != nil {
		return
	}

	msg.c2wop = c2[:strings.LastIndex(c2, ",")]

	return
}

func parseServerFirst(s1 string) (msg s1Msg, err error) {

	// Check for unsupported extensions field "m".
	if strings.HasPrefix(s1, "m=") {
		err = errors.New("SCRAM message extensions are not supported")
		return
	}

	fields := strings.Split(s1, ",")
	if len(fields) < 3 {
		err = errors.New("not enough fields in first server message")
		return
	}

	msg.nonce, err = parseField(fields[0], "r")
	if err != nil {
		return
	}

	msg.salt, err = parseFieldBase64(fields[1], "s")
	if err != nil {
		return
	}

	msg.iters, err = parseFieldInt(fields[2], "i")

	return
}

func parseServerFinal(s2 string) (msg s2Msg, err error) {
	fields := strings.Split(s2, ",")

	msg.verifier, err = parseFieldBase64(fields[0], "v")
	if err == nil {
		return
	}

	msg.err, err = parseField(fields[0], "e")

	return
}
