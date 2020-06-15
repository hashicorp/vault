package logical

import (
	"context"
	"encoding/binary"
	"errors"
	"github.com/pierrec/lz4"
)

type Compressor struct {
	physical Storage
}

var ErrUnsupportedVersion = errors.New("unsupported compression version")
var ErrUnsupportedAlgorithm = errors.New("unsupported compression algorithm")
var ErrSizeMismatch = errors.New("decompressed size unexpected")

const v1 = 1
const headerLen = 6
const algoLZ4 = 1

func NewCompressor(store Storage) *Compressor {
	return &Compressor{
		physical: store,
	}
}

func (c *Compressor) Put(ctx context.Context, entry *StorageEntry) error {
	ht := make([]int, 65536)
	srcLen := len(entry.Value)
	dstLen := headerLen + lz4.CompressBlockBound(srcLen)
	var dst []byte
	if dstLen <= srcLen {
		dst = entry.Value
	} else {
		dst = make([]byte, dstLen)
	}

	sz, err := lz4.CompressBlock(entry.Value, dst[headerLen:], ht)
	if err != nil {
		return err
	}
	dst[0] = v1
	dst[1] = algoLZ4
	binary.LittleEndian.PutUint32(dst[2:6], uint32(sz))
	entry.Value = dst[0 : sz+headerLen]

	return c.physical.Put(ctx, entry)
}

func (c *Compressor) Get(ctx context.Context, key string) (*StorageEntry, error) {
	entry, err := c.physical.Get(ctx, key)
	if err != nil || entry == nil {
		return entry, err
	}

	if entry.Value == nil {
		return entry, nil
	}
	if entry.Value[0] != v1 {
		return nil, ErrUnsupportedVersion
	}
	if entry.Value[1] != algoLZ4 {
		return nil, ErrUnsupportedAlgorithm
	}

	sz := binary.LittleEndian.Uint32(entry.Value[2:6])
	dst := make([]byte, 0, sz*2)
	si, err := lz4.UncompressBlock(entry.Value[headerLen:headerLen+sz], dst)
	if err != nil {
		return nil, err
	}
	if sz != uint32(si) {
		return nil, ErrSizeMismatch
	}
	entry.Value = dst
	return entry, nil
}

func (c *Compressor) Delete(ctx context.Context, key string) error {
	return c.physical.Delete(ctx, key)
}

func (c *Compressor) List(ctx context.Context, prefix string) ([]string, error) {
	return c.physical.List(ctx, prefix)
}
