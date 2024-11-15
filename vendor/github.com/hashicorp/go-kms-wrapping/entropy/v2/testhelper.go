// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package entropy

type mockSourcer struct{}

// simulates a successful sourcer
func (m *mockSourcer) GetRandom(bytes int) ([]byte, error) {
	return make([]byte, bytes), nil
}

// provide a mock entropy.Reader
func NewMockRandomReader() *Reader {
	return &Reader{new(mockSourcer)}
}
