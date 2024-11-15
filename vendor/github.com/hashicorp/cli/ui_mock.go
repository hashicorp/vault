// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cli

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
	"sync"
)

// NewMockUi returns a fully initialized MockUi instance
// which is safe for concurrent use.
func NewMockUi() *MockUi {
	m := new(MockUi)
	m.once.Do(m.init)
	return m
}

// MockUi is a mock UI that is used for tests and is exported publicly
// for use in external tests if needed as well. Do not instantite this
// directly since the buffers will be initialized on the first write. If
// there is no write then you will get a nil panic. Please use the
// NewMockUi() constructor function instead. You can fix your code with
//
//  sed -i -e 's/new(cli.MockUi)/cli.NewMockUi()/g' *_test.go
type MockUi struct {
	InputReader  io.Reader
	ErrorWriter  *syncBuffer
	OutputWriter *syncBuffer

	once sync.Once
}

func (u *MockUi) Ask(query string) (string, error) {
	u.once.Do(u.init)

	var result string
	fmt.Fprint(u.OutputWriter, query)
	r := bufio.NewReader(u.InputReader)
	line, err := r.ReadString('\n')
	if err != nil {
		return "", err
	}
	result = strings.TrimRight(line, "\r\n")

	return result, nil
}

func (u *MockUi) AskSecret(query string) (string, error) {
	return u.Ask(query)
}

func (u *MockUi) Error(message string) {
	u.once.Do(u.init)

	fmt.Fprint(u.ErrorWriter, message)
	fmt.Fprint(u.ErrorWriter, "\n")
}

func (u *MockUi) Info(message string) {
	u.Output(message)
}

func (u *MockUi) Output(message string) {
	u.once.Do(u.init)

	fmt.Fprint(u.OutputWriter, message)
	fmt.Fprint(u.OutputWriter, "\n")
}

func (u *MockUi) Warn(message string) {
	u.once.Do(u.init)

	fmt.Fprint(u.ErrorWriter, message)
	fmt.Fprint(u.ErrorWriter, "\n")
}

func (u *MockUi) init() {
	u.ErrorWriter = new(syncBuffer)
	u.OutputWriter = new(syncBuffer)
}

type syncBuffer struct {
	sync.RWMutex
	b bytes.Buffer
}

func (b *syncBuffer) Write(data []byte) (int, error) {
	b.Lock()
	defer b.Unlock()
	return b.b.Write(data)
}

func (b *syncBuffer) Read(data []byte) (int, error) {
	b.RLock()
	defer b.RUnlock()
	return b.b.Read(data)
}

func (b *syncBuffer) Reset() {
	b.Lock()
	b.b.Reset()
	b.Unlock()
}

func (b *syncBuffer) String() string {
	return string(b.Bytes())
}

func (b *syncBuffer) Bytes() []byte {
	b.RLock()
	data := b.b.Bytes()
	b.RUnlock()
	return data
}
