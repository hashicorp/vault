package protocol

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math"

	"github.com/SAP/go-hdb/driver/internal/protocol/encoding"
)

const (
	traceMsg = "PROT"

	prefixDB     = "←"
	prefixClient = "→"

	textIni    = "INI"
	textMsgHdr = "MSH"
	textSegHdr = "SGH"
	textParHdr = "PRH"
	textPar    = "PRT"
	textSkip   = "*skipped"
)

// padding.
const padding = 8

func padBytes(size int) int {
	if r := size % padding; r != 0 {
		return padding - r
	}
	return 0
}

type partCache map[PartKind]Part

func (c *partCache) get(kind PartKind) (Part, bool) {
	if part, ok := (*c)[kind]; ok {
		return part, true
	}
	part := newGenPartReader(kind)
	if part == nil { // part cannot be instantiated generically
		return nil, false
	}
	(*c)[kind] = part
	return part, true
}

// Reader represents a protocol reader.
type Reader struct {
	dec *encoding.Decoder

	protTrace bool
	logger    *slog.Logger

	lobChunkSize int

	readFromDB bool
	prefix     string

	mh *messageHeader
	sh *segmentHeader
	ph *partHeader

	partCache partCache

	hdbErrors    *HdbErrors
	rowsAffected *rowsAffected

	cancelled bool
}

func newReader(dec *encoding.Decoder, protTrace bool, logger *slog.Logger, lobChunkSize int, readFromDB bool, prefix string) *Reader {
	return &Reader{
		dec:          dec,
		protTrace:    protTrace,
		logger:       logger,
		lobChunkSize: lobChunkSize,
		readFromDB:   readFromDB,
		prefix:       prefix,
		partCache:    partCache{},
		mh:           &messageHeader{},
		sh:           &segmentHeader{},
		ph:           &partHeader{},
		hdbErrors:    &HdbErrors{},
		rowsAffected: &rowsAffected{},
	}
}

// NewDBReader returns an instance of a database protocol reader.
func NewDBReader(dec *encoding.Decoder, protTrace bool, logger *slog.Logger, lobChunkSize int) *Reader {
	return newReader(dec, protTrace, logger, lobChunkSize, true, prefixDB)
}

// NewClientReader returns an instance of a client protocol reader.
func NewClientReader(dec *encoding.Decoder, protTrace bool, logger *slog.Logger, lobChunkSize int) *Reader {
	return newReader(dec, protTrace, logger, lobChunkSize, false, prefixClient)
}

// Cancelled returns true if reading got cancelled, else otherwise.
func (r *Reader) Cancelled() bool { return r.cancelled }

// SkipParts reads and discards all protocol parts.
func (r *Reader) SkipParts(ctx context.Context) error {
	_, err := r.IterateParts(ctx, 0, nil)
	return err
}

// SessionID returns the session ID.
func (r *Reader) SessionID() int64 { return r.mh.sessionID }

// FunctionCode returns the function code of the protocol.
func (r *Reader) FunctionCode() FunctionCode { return r.sh.functionCode }

// ReadProlog reads the protocol prolog.
func (r *Reader) ReadProlog(ctx context.Context) error {
	if r.readFromDB {
		rep := &initReply{}
		if err := rep.decode(r.dec); err != nil {
			return err
		}
		if r.protTrace {
			r.logger.LogAttrs(ctx, slog.LevelInfo, traceMsg, slog.String(r.prefix+textIni, rep.String()))
		}
		return nil
	}
	req := &initRequest{}
	if err := req.decode(r.dec); err != nil {
		return err
	}
	if r.protTrace {
		r.logger.LogAttrs(ctx, slog.LevelInfo, traceMsg, slog.String(r.prefix+textIni, req.String()))
	}
	return nil
}

func (r *Reader) skipPadding() int {
	padBytes := padBytes(int(r.ph.bufferLength))
	r.dec.Skip(padBytes)
	return padBytes
}

func (r *Reader) skipPaddingLastPart(numReadByte int64) {
	// last part:
	// skip difference between real read bytes and message header var part length
	padBytes := int64(r.mh.varPartLength) - numReadByte
	switch {
	case padBytes < 0:
		panic(fmt.Sprintf("protocol error: bytes read %d > variable part length %d", numReadByte, r.mh.varPartLength))
	case padBytes > 0:
		r.dec.Skip(int(padBytes))
	}
}

// ReadPart reads a one protocol part.
func (r *Reader) ReadPart(ctx context.Context, part Part, lobReader LobReader) (err error) {
	cntBefore := r.dec.Cnt()

	switch part := part.(type) {
	// do not return here in case of error -> read stream would be broken
	case partDecoder:
		err = part.decode(r.dec)
	case bufLenPartDecoder:
		err = part.decodeBufLen(r.dec, r.ph.bufLen())
	case numArgPartDecoder:
		err = part.decodeNumArg(r.dec, r.ph.numArg())
	case resultPartDecoder:
		if lobReader == nil {
			panic("missing lob reader") // should never happen
		}
		err = part.decodeResult(r.dec, r.ph.numArg(), lobReader, r.lobChunkSize)
	default:
		panic("invalid part decoder") // should never happen
	}
	// do not return here in case of error -> read stream would be broken

	cnt := r.dec.Cnt() - cntBefore

	if r.protTrace {
		r.logger.LogAttrs(ctx, slog.LevelInfo, traceMsg, slog.String(r.prefix+textPar, part.String()))
	}

	bufferLen := int(r.ph.bufferLength)
	switch {
	case cnt < bufferLen: // protocol buffer length > read bytes -> skip the unread bytes
		r.dec.Skip(bufferLen - cnt)
	case cnt > bufferLen: // read bytes > protocol buffer length -> should never happen
		panic(fmt.Sprintf("protocol error: read bytes %d > buffer length %d", cnt, bufferLen))
	}
	return err
}

// ErrSkipped is used by the caller of the iterator to indicate that no read was executed.
var ErrSkipped = errors.New("perts iterator: skipped read")

// IterateParts iterates through all protocol parts.
func (r *Reader) IterateParts(ctx context.Context, offset int, fn func(kind PartKind, attrs PartAttributes) error) (int64, error) {
	var hdbErrors *HdbErrors
	var rowsAffected *rowsAffected

	if err := r.mh.decode(r.dec); err != nil {
		return 0, err
	}

	var numReadByte int64 = 0 // header bytes are not calculated in header varPartBytes: start with zero
	if r.protTrace {
		r.logger.LogAttrs(ctx, slog.LevelInfo, traceMsg, slog.String(r.prefix+textMsgHdr, r.mh.String()))
	}

	for range int(r.mh.noOfSegm) {

		// check if ctx got cancelled
		if err := ctx.Err(); err != nil {
			r.cancelled = true
			return 0, err
		}

		if err := r.sh.decode(r.dec); err != nil {
			return 0, err
		}

		numReadByte += segmentHeaderSize

		if r.protTrace {
			r.logger.LogAttrs(ctx, slog.LevelInfo, traceMsg, slog.String(r.prefix+textSegHdr, r.sh.String()))
		}

		lastPart := int(r.sh.noOfParts) - 1
		for j := range lastPart + 1 { // <=
			if err := r.ph.decode(r.dec); err != nil {
				return 0, err
			}
			kind := r.ph.partKind

			numReadByte += partHeaderSize

			if r.protTrace {
				r.logger.LogAttrs(ctx, slog.LevelInfo, traceMsg, slog.String(r.prefix+textParHdr, r.ph.String()))
			}

			cntBefore := r.dec.Cnt()

			switch kind {
			case pkRowsAffected:
				if err := r.ReadPart(ctx, r.rowsAffected, nil); err != nil {
					return 0, err
				}
				rowsAffected = r.rowsAffected
			case pkError:
				if err := r.ReadPart(ctx, r.hdbErrors, nil); err != nil {
					return 0, err
				}
				hdbErrors = r.hdbErrors
			default:
				err := ErrSkipped
				// caller must not handle hdb errors and rows affected.
				if fn != nil {
					if err = fn(kind, r.ph.partAttributes); err != nil && err != ErrSkipped { //nolint:errorlint
						return 0, err
					}
				}
				if err == ErrSkipped { //nolint:errorlint
					// if trace is on or mandatory parts need to be read we cannot skip
					if r.protTrace {
						if part, ok := r.partCache.get(kind); ok {
							if err := r.ReadPart(ctx, part, nil); err != nil {
								return 0, err
							}
						} else {
							r.dec.Skip(int(r.ph.bufferLength))
							r.logger.LogAttrs(ctx, slog.LevelInfo, traceMsg, slog.String(r.prefix+textSkip, kind.String()))
						}
					} else {
						r.dec.Skip(int(r.ph.bufferLength))
					}
				}
			}

			numReadByte += int64(r.dec.Cnt()) - int64(cntBefore)

			if j != lastPart { // not last part
				numReadByte += int64(r.skipPadding())
			}

		}
	}

	r.skipPaddingLastPart(numReadByte)

	if err := r.dec.Error(); err != nil {
		r.dec.ResetError()
		return 0, err
	}

	var numRow int64
	if rowsAffected != nil {
		numRow = rowsAffected.Total()
	}

	if hdbErrors == nil {
		return numRow, nil
	}

	if rowsAffected != nil { // link statement to error
		j := 0
		for i, rows := range rowsAffected.rows {
			if rows == raExecutionFailed {
				hdbErrors.setStmtNo(j, offset+i)
				j++
			}
		}
	}
	if hdbErrors.onlyWarnings {
		for _, err := range hdbErrors.errs {
			r.logger.LogAttrs(ctx, slog.LevelWarn, err.Error())
		}
		return numRow, nil
	}
	return numRow, hdbErrors
}

const defaultSessionID = -1

// Writer represents a protocol writer.
type Writer struct {
	wr  *bufio.Writer
	enc *encoding.Encoder

	protTrace bool
	logger    *slog.Logger

	sv     map[string]string
	svSent bool

	sessionID int64

	// reuse header
	mh *messageHeader
	sh *segmentHeader
	ph *partHeader

	cancelledOrError bool
}

// NewWriter returns an instance of a protocol writer.
func NewWriter(wr *bufio.Writer, enc *encoding.Encoder, protTrace bool, logger *slog.Logger, sv map[string]string) *Writer {
	return &Writer{
		wr:        wr,
		enc:       enc,
		protTrace: protTrace,
		logger:    logger,
		sv:        sv,
		sessionID: defaultSessionID,
		mh:        new(messageHeader),
		sh:        new(segmentHeader),
		ph:        new(partHeader),
	}
}

const (
	productVersionMajor  = 4
	productVersionMinor  = 20
	protocolVersionMajor = 4
	protocolVersionMinor = 1
)

// CancelledOrError returns true if writing got cancelled or raised an error, else otherwise.
func (w *Writer) CancelledOrError() bool { return w.cancelledOrError }

// WriteProlog writes the protocol prolog.
func (w *Writer) WriteProlog(ctx context.Context) error {
	req := &initRequest{}
	req.product.major = productVersionMajor
	req.product.minor = productVersionMinor
	req.protocol.major = protocolVersionMajor
	req.protocol.minor = protocolVersionMinor
	req.numOptions = 1
	req.endianess = littleEndian
	if err := req.encode(w.enc); err != nil {
		return err
	}
	if w.protTrace {
		w.logger.LogAttrs(ctx, slog.LevelInfo, traceMsg, slog.String(prefixClient+textIni, req.String()))
	}
	return w.wr.Flush()
}

// SetSessionID sets the session ID after a successful authentication.
func (w *Writer) SetSessionID(sessionID int64) { w.sessionID = sessionID }

func (w *Writer) Write(ctx context.Context, messageType MessageType, commit bool, parts ...PartEncoder) error {
	err := w._write(ctx, messageType, commit, parts...)
	if err != nil {
		w.cancelledOrError = true
	}
	return err
}

func (w *Writer) _write(ctx context.Context, messageType MessageType, commit bool, parts ...PartEncoder) error {
	// check on session variables to be send as ClientInfo
	if w.sv != nil && !w.svSent && messageType.ClientInfoSupported() {
		parts = append([]PartEncoder{(*clientInfo)(&w.sv)}, parts...)
		w.svSent = true
	}

	numPart := len(parts)
	partSize := make([]int, numPart)
	size := int64(segmentHeaderSize + numPart*partHeaderSize) // int64 to hold MaxUInt32 in 32bit OS

	for i, part := range parts {
		s := part.size()
		size += int64(s + padBytes(s))
		partSize[i] = s // buffer size (expensive calculation)
	}

	if size > math.MaxUint32 {
		return fmt.Errorf("message size %d exceeds maximum message header value %d", size, int64(math.MaxUint32)) // int64: without cast overflow error in 32bit OS
	}

	bufferSize := size

	w.mh.sessionID = w.sessionID
	w.mh.varPartLength = uint32(size)     //nolint: gosec
	w.mh.varPartSize = uint32(bufferSize) //nolint: gosec
	w.mh.noOfSegm = 1

	if err := w.mh.encode(w.enc); err != nil {
		return err
	}
	if w.protTrace {
		w.logger.LogAttrs(ctx, slog.LevelInfo, traceMsg, slog.String(prefixClient+textMsgHdr, w.mh.String()))
	}

	if size > math.MaxInt32 {
		return fmt.Errorf("message size %d exceeds maximum part header value %d", size, math.MaxInt32)
	}

	w.sh.messageType = messageType
	w.sh.commit = commit
	w.sh.segmentKind = skRequest
	w.sh.segmentLength = int32(size) //nolint: gosec
	w.sh.segmentOfs = 0
	w.sh.noOfParts = int16(numPart) //nolint: gosec
	w.sh.segmentNo = 1

	if err := w.sh.encode(w.enc); err != nil {
		return err
	}
	if w.protTrace {
		w.logger.LogAttrs(ctx, slog.LevelInfo, traceMsg, slog.String(prefixClient+textSegHdr, w.sh.String()))
	}

	bufferSize -= segmentHeaderSize

	for i, part := range parts {
		// check if ctx got cancelled
		if err := ctx.Err(); err != nil {
			return err
		}

		size := partSize[i]
		pad := padBytes(size)

		w.ph.partKind = part.kind()
		if err := w.ph.setNumArg(part.numArg()); err != nil {
			return err
		}
		w.ph.bufferLength = int32(size)     //nolint: gosec
		w.ph.bufferSize = int32(bufferSize) //nolint: gosec

		if err := w.ph.encode(w.enc); err != nil {
			return err
		}
		if w.protTrace {
			w.logger.LogAttrs(ctx, slog.LevelInfo, traceMsg, slog.String(prefixClient+textParHdr, w.ph.String()))
		}

		if err := part.encode(w.enc); err != nil {
			return err
		}
		if w.protTrace {
			w.logger.LogAttrs(ctx, slog.LevelInfo, traceMsg, slog.String(prefixClient+textPar, part.String()))
		}

		w.enc.Zeroes(pad)

		bufferSize -= int64(partHeaderSize + size + pad)
	}
	return w.wr.Flush()
}
