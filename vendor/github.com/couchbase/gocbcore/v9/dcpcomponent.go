package gocbcore

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"sync/atomic"

	"github.com/couchbase/gocbcore/v9/memd"
)

type dcpComponent struct {
	kvMux           *kvMux
	streamIDEnabled bool
}

func newDcpComponent(kvMux *kvMux, streamIDEnabled bool) *dcpComponent {
	return &dcpComponent{
		kvMux:           kvMux,
		streamIDEnabled: streamIDEnabled,
	}
}

func (dcp *dcpComponent) OpenStream(vbID uint16, flags memd.DcpStreamAddFlag, vbUUID VbUUID, startSeqNo,
	endSeqNo, snapStartSeqNo, snapEndSeqNo SeqNo, evtHandler StreamObserver, opts OpenStreamOptions,
	cb OpenStreamCallback) (PendingOp, error) {
	var req *memdQRequest
	var openHandled uint32
	handler := func(resp *memdQResponse, _ *memdQRequest, err error) {
		if resp == nil && err == nil {
			logWarnf("DCP event occurred with no error and no response")
			return
		}

		if err != nil {
			if resp == nil {
				if atomic.CompareAndSwapUint32(&openHandled, 0, 1) {
					// If open hasn't been handled and there's no response then it's reasonably safe to assume that
					// this occurring for the open stream request.
					cb(nil, err)
					return
				}
			}

			if resp != nil && resp.Magic == memd.CmdMagicRes {
				atomic.StoreUint32(&openHandled, 1)
				// CmdMagicRes means that this must be the open stream request response.
				cb(nil, err)
				return
			}

			var streamID uint16
			if opts.StreamOptions != nil {
				streamID = opts.StreamOptions.StreamID
			}
			evtHandler.End(vbID, streamID, err)
			return
		}

		if resp.Magic == memd.CmdMagicRes {
			atomic.StoreUint32(&openHandled, 1)
			// This is the response to the open stream request.
			numEntries := len(resp.Value) / 16
			entries := make([]FailoverEntry, numEntries)
			for i := 0; i < numEntries; i++ {
				entries[i] = FailoverEntry{
					VbUUID: VbUUID(binary.BigEndian.Uint64(resp.Value[i*16+0:])),
					SeqNo:  SeqNo(binary.BigEndian.Uint64(resp.Value[i*16+8:])),
				}
			}

			cb(entries, nil)
			return
		}

		// This is one of the stream events
		switch resp.Command {
		case memd.CmdDcpSnapshotMarker:
			vbID := resp.Vbucket
			newStartSeqNo := binary.BigEndian.Uint64(resp.Extras[0:])
			newEndSeqNo := binary.BigEndian.Uint64(resp.Extras[8:])
			snapshotType := binary.BigEndian.Uint32(resp.Extras[16:])
			var streamID uint16
			if resp.StreamIDFrame != nil {
				streamID = resp.StreamIDFrame.StreamID
			}
			evtHandler.SnapshotMarker(newStartSeqNo, newEndSeqNo, vbID, streamID, SnapshotState(snapshotType))
		case memd.CmdDcpMutation:
			vbID := resp.Vbucket
			seqNo := binary.BigEndian.Uint64(resp.Extras[0:])
			revNo := binary.BigEndian.Uint64(resp.Extras[8:])
			flags := binary.BigEndian.Uint32(resp.Extras[16:])
			expiry := binary.BigEndian.Uint32(resp.Extras[20:])
			lockTime := binary.BigEndian.Uint32(resp.Extras[24:])
			var streamID uint16
			if resp.StreamIDFrame != nil {
				streamID = resp.StreamIDFrame.StreamID
			}
			evtHandler.Mutation(seqNo, revNo, flags, expiry, lockTime, resp.Cas, resp.Datatype, vbID, resp.CollectionID, streamID, resp.Key, resp.Value)
		case memd.CmdDcpDeletion:
			vbID := resp.Vbucket
			seqNo := binary.BigEndian.Uint64(resp.Extras[0:])
			revNo := binary.BigEndian.Uint64(resp.Extras[8:])
			var deleteTime uint32
			if len(resp.Extras) == 21 {
				// Length of 21 indicates a v2 packet
				deleteTime = binary.BigEndian.Uint32(resp.Extras[16:])
			}

			var streamID uint16
			if resp.StreamIDFrame != nil {
				streamID = resp.StreamIDFrame.StreamID
			}
			evtHandler.Deletion(seqNo, revNo, deleteTime, resp.Cas, resp.Datatype, vbID, resp.CollectionID, streamID, resp.Key, resp.Value)
		case memd.CmdDcpExpiration:
			vbID := resp.Vbucket
			seqNo := binary.BigEndian.Uint64(resp.Extras[0:])
			revNo := binary.BigEndian.Uint64(resp.Extras[8:])
			var deleteTime uint32
			if len(resp.Extras) > 16 {
				deleteTime = binary.BigEndian.Uint32(resp.Extras[16:])
			}

			var streamID uint16
			if resp.StreamIDFrame != nil {
				streamID = resp.StreamIDFrame.StreamID
			}
			evtHandler.Expiration(seqNo, revNo, deleteTime, resp.Cas, vbID, resp.CollectionID, streamID, resp.Key)
		case memd.CmdDcpEvent:
			vbID := resp.Vbucket
			seqNo := binary.BigEndian.Uint64(resp.Extras[0:])
			eventCode := memd.StreamEventCode(binary.BigEndian.Uint32(resp.Extras[8:]))
			version := resp.Extras[12]
			var streamID uint16
			if resp.StreamIDFrame != nil {
				streamID = resp.StreamIDFrame.StreamID
			}

			switch eventCode {
			case memd.StreamEventCollectionCreate:
				manifestUID := binary.BigEndian.Uint64(resp.Value[0:])
				scopeID := binary.BigEndian.Uint32(resp.Value[8:])
				collectionID := binary.BigEndian.Uint32(resp.Value[12:])
				var ttl uint32
				if version == 1 {
					ttl = binary.BigEndian.Uint32(resp.Value[16:])
				}
				evtHandler.CreateCollection(seqNo, version, vbID, manifestUID, scopeID, collectionID, ttl, streamID, resp.Key)
			case memd.StreamEventCollectionDelete:
				manifestUID := binary.BigEndian.Uint64(resp.Value[0:])
				scopeID := binary.BigEndian.Uint32(resp.Value[8:])
				collectionID := binary.BigEndian.Uint32(resp.Value[12:])
				evtHandler.DeleteCollection(seqNo, version, vbID, manifestUID, scopeID, collectionID, streamID)
			case memd.StreamEventCollectionFlush:
				manifestUID := binary.BigEndian.Uint64(resp.Value[0:])
				collectionID := binary.BigEndian.Uint32(resp.Value[8:])
				evtHandler.FlushCollection(seqNo, version, vbID, manifestUID, collectionID)
			case memd.StreamEventScopeCreate:
				manifestUID := binary.BigEndian.Uint64(resp.Value[0:])
				scopeID := binary.BigEndian.Uint32(resp.Value[8:])
				evtHandler.CreateScope(seqNo, version, vbID, manifestUID, scopeID, streamID, resp.Key)
			case memd.StreamEventScopeDelete:
				manifestUID := binary.BigEndian.Uint64(resp.Value[0:])
				scopeID := binary.BigEndian.Uint32(resp.Value[8:])
				evtHandler.DeleteScope(seqNo, version, vbID, manifestUID, scopeID, streamID)
			case memd.StreamEventCollectionChanged:
				manifestUID := binary.BigEndian.Uint64(resp.Value[0:])
				collectionID := binary.BigEndian.Uint32(resp.Value[8:])
				ttl := binary.BigEndian.Uint32(resp.Value[12:])
				evtHandler.ModifyCollection(seqNo, version, vbID, manifestUID, collectionID, ttl, streamID)
			}
		case memd.CmdDcpStreamEnd:
			vbID := resp.Vbucket
			code := memd.StreamEndStatus(binary.BigEndian.Uint32(resp.Extras[0:]))
			var streamID uint16
			if resp.StreamIDFrame != nil {
				streamID = resp.StreamIDFrame.StreamID
			}
			evtHandler.End(vbID, streamID, getStreamEndStatusError(code))
			req.internalCancel(err)
		case memd.CmdDcpOsoSnapshot:
			vbID := resp.Vbucket
			snapshotType := binary.BigEndian.Uint32(resp.Extras[0:])
			var streamID uint16
			if resp.StreamIDFrame != nil {
				streamID = resp.StreamIDFrame.StreamID
			}
			evtHandler.OSOSnapshot(vbID, snapshotType, streamID)
		case memd.CmdDcpSeqNoAdvanced:
			vbID := resp.Vbucket
			seqno := binary.BigEndian.Uint64(resp.Extras[0:])
			var streamID uint16
			if resp.StreamIDFrame != nil {
				streamID = resp.StreamIDFrame.StreamID
			}
			evtHandler.SeqNoAdvanced(vbID, seqno, streamID)
		}
	}

	extraBuf := make([]byte, 48)
	binary.BigEndian.PutUint32(extraBuf[0:], uint32(flags))
	binary.BigEndian.PutUint32(extraBuf[4:], 0)
	binary.BigEndian.PutUint64(extraBuf[8:], uint64(startSeqNo))
	binary.BigEndian.PutUint64(extraBuf[16:], uint64(endSeqNo))
	binary.BigEndian.PutUint64(extraBuf[24:], uint64(vbUUID))
	binary.BigEndian.PutUint64(extraBuf[32:], uint64(snapStartSeqNo))
	binary.BigEndian.PutUint64(extraBuf[40:], uint64(snapEndSeqNo))

	var val []byte
	val = nil
	if opts.StreamOptions != nil || opts.FilterOptions != nil || opts.ManifestOptions != nil {
		convertedFilter := streamFilter{}

		if opts.FilterOptions != nil {
			// If there are collection IDs then we can assume that scope ID of 0 actually means no scope ID
			if len(opts.FilterOptions.CollectionIDs) > 0 {
				for _, cid := range opts.FilterOptions.CollectionIDs {
					convertedFilter.Collections = append(convertedFilter.Collections, fmt.Sprintf("%x", cid))
				}
			} else {
				// No collection IDs but the filter was set so even if scope ID is 0 then we use it
				convertedFilter.Scope = fmt.Sprintf("%x", opts.FilterOptions.ScopeID)
			}

		}
		if opts.ManifestOptions != nil {
			convertedFilter.ManifestUID = fmt.Sprintf("%x", opts.ManifestOptions.ManifestUID)
		}
		if opts.StreamOptions != nil {
			convertedFilter.StreamID = opts.StreamOptions.StreamID
		}

		var err error
		val, err = json.Marshal(convertedFilter)
		if err != nil {
			return nil, err
		}
	}

	req = &memdQRequest{
		Packet: memd.Packet{
			Magic:    memd.CmdMagicReq,
			Command:  memd.CmdDcpStreamReq,
			Datatype: 0,
			Cas:      0,
			Extras:   extraBuf,
			Key:      nil,
			Value:    val,
			Vbucket:  vbID,
		},
		Callback:   handler,
		ReplicaIdx: 0,
		Persistent: true,
	}
	return dcp.kvMux.DispatchDirect(req)
}

func (dcp *dcpComponent) CloseStream(vbID uint16, opts CloseStreamOptions, cb CloseStreamCallback) (PendingOp, error) {
	handler := func(_ *memdQResponse, _ *memdQRequest, err error) {
		cb(err)
	}

	var streamFrame *memd.StreamIDFrame
	if opts.StreamOptions != nil {
		if !dcp.streamIDEnabled {
			return nil, errStreamIDNotEnabled
		}

		streamFrame = &memd.StreamIDFrame{
			StreamID: opts.StreamOptions.StreamID,
		}
	}

	req := &memdQRequest{
		Packet: memd.Packet{
			Magic:         memd.CmdMagicReq,
			Command:       memd.CmdDcpCloseStream,
			Datatype:      0,
			Cas:           0,
			Extras:        nil,
			Key:           nil,
			Value:         nil,
			Vbucket:       vbID,
			StreamIDFrame: streamFrame,
		},
		Callback:      handler,
		ReplicaIdx:    0,
		Persistent:    false,
		RetryStrategy: newFailFastRetryStrategy(),
	}

	return dcp.kvMux.DispatchDirect(req)
}

func (dcp *dcpComponent) GetFailoverLog(vbID uint16, cb GetFailoverLogCallback) (PendingOp, error) {
	handler := func(resp *memdQResponse, _ *memdQRequest, err error) {
		if err != nil {
			cb(nil, err)
			return
		}

		numEntries := len(resp.Value) / 16
		entries := make([]FailoverEntry, numEntries)
		for i := 0; i < numEntries; i++ {
			entries[i] = FailoverEntry{
				VbUUID: VbUUID(binary.BigEndian.Uint64(resp.Value[i*16+0:])),
				SeqNo:  SeqNo(binary.BigEndian.Uint64(resp.Value[i*16+8:])),
			}
		}
		cb(entries, nil)
	}

	req := &memdQRequest{
		Packet: memd.Packet{
			Magic:    memd.CmdMagicReq,
			Command:  memd.CmdDcpGetFailoverLog,
			Datatype: 0,
			Cas:      0,
			Extras:   nil,
			Key:      nil,
			Value:    nil,
			Vbucket:  vbID,
		},
		Callback:      handler,
		ReplicaIdx:    0,
		Persistent:    false,
		RetryStrategy: newFailFastRetryStrategy(),
	}
	return dcp.kvMux.DispatchDirect(req)
}

func (dcp *dcpComponent) GetVbucketSeqnos(serverIdx int, state memd.VbucketState, opts GetVbucketSeqnoOptions, cb GetVBucketSeqnosCallback) (PendingOp, error) {
	handler := func(resp *memdQResponse, _ *memdQRequest, err error) {
		if err != nil {
			cb(nil, err)
			return
		}

		var vbs []VbSeqNoEntry

		numVbs := len(resp.Value) / 10
		for i := 0; i < numVbs; i++ {
			vbs = append(vbs, VbSeqNoEntry{
				VbID:  binary.BigEndian.Uint16(resp.Value[i*10:]),
				SeqNo: SeqNo(binary.BigEndian.Uint64(resp.Value[i*10+2:])),
			})
		}

		cb(vbs, nil)
	}

	var extraBuf []byte

	if opts.FilterOptions == nil {
		extraBuf = make([]byte, 4)
		binary.BigEndian.PutUint32(extraBuf[0:], uint32(state))
	} else {
		if !dcp.kvMux.SupportsCollections() {
			return nil, errCollectionsUnsupported
		}

		extraBuf = make([]byte, 8)
		binary.BigEndian.PutUint32(extraBuf[0:], uint32(state))
		binary.BigEndian.PutUint32(extraBuf[4:], opts.FilterOptions.CollectionID)
	}
	req := &memdQRequest{
		Packet: memd.Packet{
			Magic:    memd.CmdMagicReq,
			Command:  memd.CmdGetAllVBSeqnos,
			Datatype: 0,
			Cas:      0,
			Extras:   extraBuf,
			Key:      nil,
			Value:    nil,
			Vbucket:  0,
		},
		Callback:      handler,
		ReplicaIdx:    -serverIdx,
		Persistent:    false,
		RetryStrategy: newFailFastRetryStrategy(),
	}

	return dcp.kvMux.DispatchDirect(req)
}
