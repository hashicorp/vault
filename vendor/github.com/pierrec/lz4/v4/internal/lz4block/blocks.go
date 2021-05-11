// Package lz4block provides LZ4 BlockSize types and pools of buffers.
package lz4block

import "sync"

const (
	Block64Kb uint32 = 1 << (16 + iota*2)
	Block256Kb
	Block1Mb
	Block4Mb
)

var (
	BlockPool64K  = sync.Pool{New: func() interface{} { return make([]byte, Block64Kb) }}
	BlockPool256K = sync.Pool{New: func() interface{} { return make([]byte, Block256Kb) }}
	BlockPool1M   = sync.Pool{New: func() interface{} { return make([]byte, Block1Mb) }}
	BlockPool4M   = sync.Pool{New: func() interface{} { return make([]byte, Block4Mb) }}
)

func Index(b uint32) BlockSizeIndex {
	switch b {
	case Block64Kb:
		return 4
	case Block256Kb:
		return 5
	case Block1Mb:
		return 6
	case Block4Mb:
		return 7
	}
	return 0
}

func IsValid(b uint32) bool {
	return Index(b) > 0
}

type BlockSizeIndex uint8

func (b BlockSizeIndex) IsValid() bool {
	switch b {
	case 4, 5, 6, 7:
		return true
	}
	return false
}

func (b BlockSizeIndex) Get() []byte {
	var buf interface{}
	switch b {
	case 4:
		buf = BlockPool64K.Get()
	case 5:
		buf = BlockPool256K.Get()
	case 6:
		buf = BlockPool1M.Get()
	case 7:
		buf = BlockPool4M.Get()
	}
	return buf.([]byte)
}

func (b BlockSizeIndex) Put(buf []byte) {
	// Safeguard: do not allow invalid buffers.
	switch c := uint32(cap(buf)); b {
	case 4:
		if c == Block64Kb {
			BlockPool64K.Put(buf[:c])
		}
	case 5:
		if c == Block256Kb {
			BlockPool256K.Put(buf[:c])
		}
	case 6:
		if c == Block1Mb {
			BlockPool1M.Put(buf[:c])
		}
	case 7:
		if c == Block4Mb {
			BlockPool4M.Put(buf[:c])
		}
	}
}

type CompressionLevel uint32

const Fast CompressionLevel = 0
