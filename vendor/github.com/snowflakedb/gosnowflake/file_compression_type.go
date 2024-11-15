// Copyright (c) 2021-2022 Snowflake Computing Inc. All rights reserved.

package gosnowflake

import (
	"bytes"
	"strings"

	"github.com/gabriel-vasile/mimetype"
)

type compressionType struct {
	name          string
	fileExtension string
	mimeSubtypes  []string
	isSupported   bool
}

var compressionTypes = map[string]*compressionType{
	"GZIP": {
		"GZIP",
		".gz",
		[]string{"gzip", "x-gzip"},
		true,
	},
	"DEFLATE": {
		"DEFLATE",
		".deflate",
		[]string{"zlib", "deflate"},
		true,
	},
	"RAW_DEFLATE": {
		"RAW_DEFLATE",
		".raw_deflate",
		[]string{"raw_deflate"},
		true,
	},
	"BZIP2": {
		"BZIP2",
		".bz2",
		[]string{"bzip2", "x-bzip2", "x-bz2", "x-bzip", "bz2"},
		true,
	},
	"LZIP": {
		"LZIP",
		".lz",
		[]string{"lzip", "x-lzip"},
		false,
	},
	"LZMA": {
		"LZMA",
		".lzma",
		[]string{"lzma", "x-lzma"},
		false,
	},
	"LZO": {
		"LZO",
		".lzo",
		[]string{"lzo", "x-lzo"},
		false,
	},
	"XZ": {
		"XZ",
		".xz",
		[]string{"xz", "x-xz"},
		false,
	},
	"COMPRESS": {
		"COMPRESS",
		".Z",
		[]string{"compress", "x-compress"},
		false,
	},
	"PARQUET": {
		"PARQUET",
		".parquet",
		[]string{"parquet"},
		true,
	},
	"ZSTD": {
		"ZSTD",
		".zst",
		[]string{"zstd", "x-zstd"},
		true,
	},
	"BROTLI": {
		"BROTLI",
		".br",
		[]string{"br", "x-br"},
		true,
	},
	"ORC": {
		"ORC",
		".orc",
		[]string{"orc"},
		true,
	},
}

var mimeSubTypeToCompression map[string]*compressionType
var extensionToCompression map[string]*compressionType

func init() {
	mimeSubTypeToCompression = make(map[string]*compressionType)
	extensionToCompression = make(map[string]*compressionType)
	for _, meta := range compressionTypes {
		extensionToCompression[meta.fileExtension] = meta
		for _, subtype := range meta.mimeSubtypes {
			mimeSubTypeToCompression[subtype] = meta
		}
	}
	mimetype.Extend(func(raw []byte, limit uint32) bool {
		return bytes.HasPrefix(raw, []byte("PAR1"))
	}, "snowflake/parquet", ".parquet")
	mimetype.Extend(func(raw []byte, limit uint32) bool {
		return bytes.HasPrefix(raw, []byte("ORC"))
	}, "snowflake/orc", ".orc")
}

func lookupByMimeSubType(mimeSubType string) *compressionType {
	if val, ok := mimeSubTypeToCompression[strings.ToLower(mimeSubType)]; ok {
		return val
	}
	return nil
}

func lookupByExtension(extension string) *compressionType {
	if val, ok := extensionToCompression[strings.ToLower(extension)]; ok {
		return val
	}
	return nil
}
