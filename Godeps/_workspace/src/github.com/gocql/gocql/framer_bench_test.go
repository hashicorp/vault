package gocql

import (
	"compress/gzip"
	"io/ioutil"
	"os"
	"testing"
)

func readGzipData(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r, err := gzip.NewReader(f)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	return ioutil.ReadAll(r)
}

func BenchmarkParseRowsFrame(b *testing.B) {
	data, err := readGzipData("testdata/frames/bench_parse_result.gz")
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		framer := &framer{
			header: &frameHeader{
				version: protoVersion4 | 0x80,
				op:      opResult,
				length:  len(data),
			},
			rbuf: data,
		}

		_, err = framer.parseFrame()
		if err != nil {
			b.Fatal(err)
		}
	}
}
