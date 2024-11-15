package gocbcore

import (
	"errors"
	"log"

	"github.com/couchbase/gocbcore/v10/memd"
)

var streamEndErrorMap = make(map[memd.StreamEndStatus]error)

func makeStreamEndStatusError(code memd.StreamEndStatus) error {
	err := errors.New(code.KVText())
	if streamEndErrorMap[code] != nil {
		log.Fatal("error handling setup failure")
	}
	streamEndErrorMap[code] = err
	return err
}

func getStreamEndStatusError(code memd.StreamEndStatus) error {
	if code == memd.StreamEndOK {
		return nil
	}
	if err := streamEndErrorMap[code]; err != nil {
		return err
	}
	return errors.New(code.KVText())
}

var (
	// ErrDCPStreamClosed occurs when a DCP stream is closed gracefully.
	ErrDCPStreamClosed = makeStreamEndStatusError(memd.StreamEndClosed)

	// ErrDCPStreamStateChanged occurs when a DCP stream is interrupted by failover.
	ErrDCPStreamStateChanged = makeStreamEndStatusError(memd.StreamEndStateChanged)

	// ErrDCPStreamDisconnected occurs when a DCP stream is disconnected.
	ErrDCPStreamDisconnected = makeStreamEndStatusError(memd.StreamEndDisconnected)

	// ErrDCPStreamTooSlow occurs when a DCP stream is cancelled due to the application
	// not keeping up with the rate of flow of DCP events sent by the server.
	ErrDCPStreamTooSlow = makeStreamEndStatusError(memd.StreamEndTooSlow)

	// ErrDCPBackfillFailed occurs when there was an issue starting the backfill on
	// the server e.g. the requested start seqno was behind the purge seqno.
	ErrDCPBackfillFailed = makeStreamEndStatusError(memd.StreamEndBackfillFailed)

	// ErrDCPStreamFilterEmpty occurs when all of the collections for a DCP stream are
	// dropped.
	ErrDCPStreamFilterEmpty = makeStreamEndStatusError(memd.StreamEndFilterEmpty)

	// ErrStreamIDNotEnabled occurs when dcp operations are performed using a stream ID when stream IDs are not enabled.
	ErrStreamIDNotEnabled = errors.New("stream IDs have not been enabled on this stream")

	// ErrDCPStreamIDInvalid occurs when a dcp stream ID is invalid.
	ErrDCPStreamIDInvalid = errors.New("stream ID invalid")
)
