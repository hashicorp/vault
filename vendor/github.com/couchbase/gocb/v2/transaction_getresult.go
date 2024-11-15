package gocb

import (
	"encoding/json"
	"strconv"

	"github.com/couchbase/gocbcore/v10"
)

// TransactionGetResult represents the result of a Get operation which was performed.
type TransactionGetResult struct {
	collection *Collection
	docID      string

	transcoder Transcoder
	flags      uint32

	txnMeta json.RawMessage

	coreRes *gocbcore.TransactionGetResult
}

// Content provides access to the documents contents.
func (d *TransactionGetResult) Content(valuePtr interface{}) error {
	return d.transcoder.Decode(d.coreRes.Value, d.flags, valuePtr)
}

func fromScas(scas string) (gocbcore.Cas, error) {
	i, err := strconv.ParseUint(scas, 10, 64)
	if err != nil {
		return 0, err
	}

	return gocbcore.Cas(i), nil
}

func toScas(cas gocbcore.Cas) string {
	return strconv.FormatUint(uint64(cas), 10)
}
