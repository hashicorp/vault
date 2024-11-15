package mssql

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"strconv"

	"github.com/denisenkom/go-mssqldb/msdsn"
	"github.com/golang-sql/sqlexp"
)

//go:generate go run golang.org/x/tools/cmd/stringer -type token

type token byte

// token ids
const (
	tokenReturnStatus  token = 121 // 0x79
	tokenColMetadata   token = 129 // 0x81
	tokenOrder         token = 169 // 0xA9
	tokenError         token = 170 // 0xAA
	tokenInfo          token = 171 // 0xAB
	tokenReturnValue   token = 0xAC
	tokenLoginAck      token = 173 // 0xad
	tokenFeatureExtAck token = 174 // 0xae
	tokenRow           token = 209 // 0xd1
	tokenNbcRow        token = 210 // 0xd2
	tokenEnvChange     token = 227 // 0xE3
	tokenSSPI          token = 237 // 0xED
	tokenFedAuthInfo   token = 238 // 0xEE
	tokenDone          token = 253 // 0xFD
	tokenDoneProc      token = 254
	tokenDoneInProc    token = 255
)

// done flags
// https://msdn.microsoft.com/en-us/library/dd340421.aspx
const (
	doneFinal    = 0
	doneMore     = 1
	doneError    = 2
	doneInxact   = 4
	doneCount    = 0x10
	doneAttn     = 0x20
	doneSrvError = 0x100
)

// ENVCHANGE types
// http://msdn.microsoft.com/en-us/library/dd303449.aspx
const (
	envTypDatabase           = 1
	envTypLanguage           = 2
	envTypCharset            = 3
	envTypPacketSize         = 4
	envSortId                = 5
	envSortFlags             = 6
	envSqlCollation          = 7
	envTypBeginTran          = 8
	envTypCommitTran         = 9
	envTypRollbackTran       = 10
	envEnlistDTC             = 11
	envDefectTran            = 12
	envDatabaseMirrorPartner = 13
	envPromoteTran           = 15
	envTranMgrAddr           = 16
	envTranEnded             = 17
	envResetConnAck          = 18
	envStartedInstanceName   = 19
	envRouting               = 20
)

const (
	fedAuthInfoSTSURL = 0x01
	fedAuthInfoSPN    = 0x02
)

// COLMETADATA flags
// https://msdn.microsoft.com/en-us/library/dd357363.aspx
const (
	colFlagNullable = 1
	// TODO implement more flags
)

// interface for all tokens
type tokenStruct interface{}

type orderStruct struct {
	ColIds []uint16
}

type doneStruct struct {
	Status   uint16
	CurCmd   uint16
	RowCount uint64
	errors   []Error
}

func (d doneStruct) isError() bool {
	return d.Status&doneError != 0 || len(d.errors) > 0
}

func (d doneStruct) getError() Error {
	n := len(d.errors)
	if n == 0 {
		return Error{Message: "Request failed but didn't provide reason"}
	}
	err := d.errors[n-1]
	// should this return the most severe error?
	err.All = make([]Error, n)
	copy(err.All, d.errors)
	return err
}

type doneInProcStruct doneStruct

// ENVCHANGE stream
// http://msdn.microsoft.com/en-us/library/dd303449.aspx
func processEnvChg(ctx context.Context, sess *tdsSession) {
	size := sess.buf.uint16()
	r := &io.LimitedReader{R: sess.buf, N: int64(size)}
	for {
		var err error
		var envtype uint8
		err = binary.Read(r, binary.LittleEndian, &envtype)
		if err == io.EOF {
			return
		}
		if err != nil {
			badStreamPanic(err)
		}
		switch envtype {
		case envTypDatabase:
			sess.database, err = readBVarChar(r)
			if err != nil {
				badStreamPanic(err)
			}
			_, err = readBVarChar(r)
			if err != nil {
				badStreamPanic(err)
			}
		case envTypLanguage:
			// currently ignored
			// new value
			if _, err = readBVarChar(r); err != nil {
				badStreamPanic(err)
			}
			// old value
			if _, err = readBVarChar(r); err != nil {
				badStreamPanic(err)
			}
		case envTypCharset:
			// currently ignored
			// new value
			if _, err = readBVarChar(r); err != nil {
				badStreamPanic(err)
			}
			// old value
			if _, err = readBVarChar(r); err != nil {
				badStreamPanic(err)
			}
		case envTypPacketSize:
			packetsize, err := readBVarChar(r)
			if err != nil {
				badStreamPanic(err)
			}
			_, err = readBVarChar(r)
			if err != nil {
				badStreamPanic(err)
			}
			packetsizei, err := strconv.Atoi(packetsize)
			if err != nil {
				badStreamPanicf("Invalid Packet size value returned from server (%s): %s", packetsize, err.Error())
			}
			sess.buf.ResizeBuffer(packetsizei)
		case envSortId:
			// currently ignored
			// new value
			if _, err = readBVarChar(r); err != nil {
				badStreamPanic(err)
			}
			// old value, should be 0
			if _, err = readBVarChar(r); err != nil {
				badStreamPanic(err)
			}
		case envSortFlags:
			// currently ignored
			// new value
			if _, err = readBVarChar(r); err != nil {
				badStreamPanic(err)
			}
			// old value, should be 0
			if _, err = readBVarChar(r); err != nil {
				badStreamPanic(err)
			}
		case envSqlCollation:
			// currently ignored
			var collationSize uint8
			err = binary.Read(r, binary.LittleEndian, &collationSize)
			if err != nil {
				badStreamPanic(err)
			}

			// SQL Collation data should contain 5 bytes in length
			if collationSize != 5 {
				badStreamPanicf("Invalid SQL Collation size value returned from server: %d", collationSize)
			}

			// 4 bytes, contains: LCID ColFlags Version
			var info uint32
			err = binary.Read(r, binary.LittleEndian, &info)
			if err != nil {
				badStreamPanic(err)
			}

			// 1 byte, contains: sortID
			var sortID uint8
			err = binary.Read(r, binary.LittleEndian, &sortID)
			if err != nil {
				badStreamPanic(err)
			}

			// old value, should be 0
			if _, err = readBVarChar(r); err != nil {
				badStreamPanic(err)
			}
		case envTypBeginTran:
			tranid, err := readBVarByte(r)
			if len(tranid) != 8 {
				badStreamPanicf("invalid size of transaction identifier: %d", len(tranid))
			}
			sess.tranid = binary.LittleEndian.Uint64(tranid)
			if err != nil {
				badStreamPanic(err)
			}
			if sess.logFlags&logTransaction != 0 {
				sess.logger.Log(ctx, msdsn.LogTransaction, fmt.Sprintf("BEGIN TRANSACTION %x", sess.tranid))
			}
			_, err = readBVarByte(r)
			if err != nil {
				badStreamPanic(err)
			}
		case envTypCommitTran, envTypRollbackTran:
			_, err = readBVarByte(r)
			if err != nil {
				badStreamPanic(err)
			}
			_, err = readBVarByte(r)
			if err != nil {
				badStreamPanic(err)
			}
			if sess.logFlags&logTransaction != 0 {
				if envtype == envTypCommitTran {
					sess.logger.Log(ctx, msdsn.LogTransaction, fmt.Sprintf("COMMIT TRANSACTION %x", sess.tranid))
				} else {
					sess.logger.Log(ctx, msdsn.LogTransaction, fmt.Sprintf("ROLLBACK TRANSACTION %x", sess.tranid))
				}
			}
			sess.tranid = 0
		case envEnlistDTC:
			// currently ignored
			// new value, should be 0
			if _, err = readBVarChar(r); err != nil {
				badStreamPanic(err)
			}
			// old value
			if _, err = readBVarChar(r); err != nil {
				badStreamPanic(err)
			}
		case envDefectTran:
			// currently ignored
			// new value
			if _, err = readBVarChar(r); err != nil {
				badStreamPanic(err)
			}
			// old value, should be 0
			if _, err = readBVarChar(r); err != nil {
				badStreamPanic(err)
			}
		case envDatabaseMirrorPartner:
			sess.partner, err = readBVarChar(r)
			if err != nil {
				badStreamPanic(err)
			}
			_, err = readBVarChar(r)
			if err != nil {
				badStreamPanic(err)
			}
		case envPromoteTran:
			// currently ignored
			// old value, should be 0
			if _, err = readBVarChar(r); err != nil {
				badStreamPanic(err)
			}
			// dtc token
			// spec says it should be L_VARBYTE, so this code might be wrong
			if _, err = readBVarChar(r); err != nil {
				badStreamPanic(err)
			}
		case envTranMgrAddr:
			// currently ignored
			// old value, should be 0
			if _, err = readBVarChar(r); err != nil {
				badStreamPanic(err)
			}
			// XACT_MANAGER_ADDRESS = B_VARBYTE
			if _, err = readBVarChar(r); err != nil {
				badStreamPanic(err)
			}
		case envTranEnded:
			// currently ignored
			// old value, B_VARBYTE
			if _, err = readBVarChar(r); err != nil {
				badStreamPanic(err)
			}
			// should be 0
			if _, err = readBVarChar(r); err != nil {
				badStreamPanic(err)
			}
		case envResetConnAck:
			// currently ignored
			// old value, should be 0
			if _, err = readBVarChar(r); err != nil {
				badStreamPanic(err)
			}
			// should be 0
			if _, err = readBVarChar(r); err != nil {
				badStreamPanic(err)
			}
		case envStartedInstanceName:
			// currently ignored
			// old value, should be 0
			if _, err = readBVarChar(r); err != nil {
				badStreamPanic(err)
			}
			// instance name
			if _, err = readBVarChar(r); err != nil {
				badStreamPanic(err)
			}
		case envRouting:
			// RoutingData message is:
			// ValueLength                 USHORT
			// Protocol (TCP = 0)          BYTE
			// ProtocolProperty (new port) USHORT
			// AlternateServer             US_VARCHAR
			_, err := readUshort(r)
			if err != nil {
				badStreamPanic(err)
			}
			protocol, err := readByte(r)
			if err != nil || protocol != 0 {
				badStreamPanic(err)
			}
			newPort, err := readUshort(r)
			if err != nil {
				badStreamPanic(err)
			}
			newServer, err := readUsVarChar(r)
			if err != nil {
				badStreamPanic(err)
			}
			// consume the OLDVALUE = %x00 %x00
			_, err = readUshort(r)
			if err != nil {
				badStreamPanic(err)
			}
			sess.routedServer = newServer
			sess.routedPort = newPort
		default:
			// ignore rest of records because we don't know how to skip those
			if sess.logFlags&logDebug != 0 {
				sess.logger.Log(ctx, msdsn.LogDebug, fmt.Sprintf("WARN: Unknown ENVCHANGE record detected with type id = %d", envtype))
			}
			return
		}
	}
}

// http://msdn.microsoft.com/en-us/library/dd358180.aspx
func parseReturnStatus(r *tdsBuffer) ReturnStatus {
	return ReturnStatus(r.int32())
}

func parseOrder(r *tdsBuffer) (res orderStruct) {
	len := int(r.uint16())
	res.ColIds = make([]uint16, len/2)
	for i := 0; i < len/2; i++ {
		res.ColIds[i] = r.uint16()
	}
	return res
}

// https://msdn.microsoft.com/en-us/library/dd340421.aspx
func parseDone(r *tdsBuffer) (res doneStruct) {
	res.Status = r.uint16()
	res.CurCmd = r.uint16()
	res.RowCount = r.uint64()
	return res
}

// https://msdn.microsoft.com/en-us/library/dd340553.aspx
func parseDoneInProc(r *tdsBuffer) (res doneInProcStruct) {
	res.Status = r.uint16()
	res.CurCmd = r.uint16()
	res.RowCount = r.uint64()
	return res
}

type sspiMsg []byte

func parseSSPIMsg(r *tdsBuffer) sspiMsg {
	size := r.uint16()
	buf := make([]byte, size)
	r.ReadFull(buf)
	return sspiMsg(buf)
}

type fedAuthInfoStruct struct {
	STSURL    string
	ServerSPN string
}

type fedAuthInfoOpt struct {
	fedAuthInfoID          byte
	dataLength, dataOffset uint32
}

func parseFedAuthInfo(r *tdsBuffer) fedAuthInfoStruct {
	size := r.uint32()

	var STSURL, SPN string
	var err error

	// Each fedAuthInfoOpt is one byte to indicate the info ID,
	// then a four byte offset and a four byte length.
	count := r.uint32()
	offset := uint32(4)
	opts := make([]fedAuthInfoOpt, count)

	for i := uint32(0); i < count; i++ {
		fedAuthInfoID := r.byte()
		dataLength := r.uint32()
		dataOffset := r.uint32()
		offset += 1 + 4 + 4

		opts[i] = fedAuthInfoOpt{
			fedAuthInfoID: fedAuthInfoID,
			dataLength:    dataLength,
			dataOffset:    dataOffset,
		}
	}

	data := make([]byte, size-offset)
	r.ReadFull(data)

	for i := uint32(0); i < count; i++ {
		if opts[i].dataOffset < offset {
			badStreamPanicf("Fed auth info opt stated data offset %d is before data begins in packet at %d",
				opts[i].dataOffset, offset)
			// returns via panic
		}

		if opts[i].dataOffset+opts[i].dataLength > size {
			badStreamPanicf("Fed auth info opt stated data length %d added to stated offset exceeds size of packet %d",
				opts[i].dataOffset+opts[i].dataLength, size)
			// returns via panic
		}

		optData := data[opts[i].dataOffset-offset : opts[i].dataOffset-offset+opts[i].dataLength]
		switch opts[i].fedAuthInfoID {
		case fedAuthInfoSTSURL:
			STSURL, err = ucs22str(optData)
		case fedAuthInfoSPN:
			SPN, err = ucs22str(optData)
		default:
			err = fmt.Errorf("unexpected fed auth info opt ID %d", int(opts[i].fedAuthInfoID))
		}

		if err != nil {
			badStreamPanic(err)
		}
	}

	return fedAuthInfoStruct{
		STSURL:    STSURL,
		ServerSPN: SPN,
	}
}

type loginAckStruct struct {
	Interface  uint8
	TDSVersion uint32
	ProgName   string
	ProgVer    uint32
}

func parseLoginAck(r *tdsBuffer) loginAckStruct {
	size := r.uint16()
	buf := make([]byte, size)
	r.ReadFull(buf)
	var res loginAckStruct
	res.Interface = buf[0]
	res.TDSVersion = binary.BigEndian.Uint32(buf[1:])
	prognamelen := buf[1+4]
	var err error
	if res.ProgName, err = ucs22str(buf[1+4+1 : 1+4+1+prognamelen*2]); err != nil {
		badStreamPanic(err)
	}
	res.ProgVer = binary.BigEndian.Uint32(buf[size-4:])
	return res
}

// https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-tds/2eb82f8e-11f0-46dc-b42d-27302fa4701a
type fedAuthAckStruct struct {
	Nonce     []byte
	Signature []byte
}

func parseFeatureExtAck(r *tdsBuffer) map[byte]interface{} {
	ack := map[byte]interface{}{}

	for feature := r.byte(); feature != featExtTERMINATOR; feature = r.byte() {
		length := r.uint32()

		switch feature {
		case featExtFEDAUTH:
			// In theory we need to know the federated authentication library to
			// know how to parse, but the alternatives provide compatible structures.
			fedAuthAck := fedAuthAckStruct{}
			if length >= 32 {
				fedAuthAck.Nonce = make([]byte, 32)
				r.ReadFull(fedAuthAck.Nonce)
				length -= 32
			}
			if length >= 32 {
				fedAuthAck.Signature = make([]byte, 32)
				r.ReadFull(fedAuthAck.Signature)
				length -= 32
			}
			ack[feature] = fedAuthAck

		}

		// Skip unprocessed bytes
		if length > 0 {
			io.CopyN(ioutil.Discard, r, int64(length))
		}
	}

	return ack
}

// http://msdn.microsoft.com/en-us/library/dd357363.aspx
func parseColMetadata72(r *tdsBuffer) (columns []columnStruct) {
	count := r.uint16()
	if count == 0xffff {
		// no metadata is sent
		return nil
	}
	columns = make([]columnStruct, count)
	for i := range columns {
		column := &columns[i]
		column.UserType = r.uint32()
		column.Flags = r.uint16()

		// parsing TYPE_INFO structure
		column.ti = readTypeInfo(r)
		column.ColName = r.BVarChar()
	}
	return columns
}

// http://msdn.microsoft.com/en-us/library/dd357254.aspx
func parseRow(r *tdsBuffer, columns []columnStruct, row []interface{}) {
	for i, column := range columns {
		row[i] = column.ti.Reader(&column.ti, r)
	}
}

// http://msdn.microsoft.com/en-us/library/dd304783.aspx
func parseNbcRow(r *tdsBuffer, columns []columnStruct, row []interface{}) {
	bitlen := (len(columns) + 7) / 8
	pres := make([]byte, bitlen)
	r.ReadFull(pres)
	for i, col := range columns {
		if pres[i/8]&(1<<(uint(i)%8)) != 0 {
			row[i] = nil
			continue
		}
		row[i] = col.ti.Reader(&col.ti, r)
	}
}

// http://msdn.microsoft.com/en-us/library/dd304156.aspx
func parseError72(r *tdsBuffer) (res Error) {
	length := r.uint16()
	_ = length // ignore length
	res.Number = r.int32()
	res.State = r.byte()
	res.Class = r.byte()
	res.Message = r.UsVarChar()
	res.ServerName = r.BVarChar()
	res.ProcName = r.BVarChar()
	res.LineNo = r.int32()
	return
}

// http://msdn.microsoft.com/en-us/library/dd304156.aspx
func parseInfo(r *tdsBuffer) (res Error) {
	length := r.uint16()
	_ = length // ignore length
	res.Number = r.int32()
	res.State = r.byte()
	res.Class = r.byte()
	res.Message = r.UsVarChar()
	res.ServerName = r.BVarChar()
	res.ProcName = r.BVarChar()
	res.LineNo = r.int32()
	return
}

// https://msdn.microsoft.com/en-us/library/dd303881.aspx
func parseReturnValue(r *tdsBuffer) (nv namedValue) {
	/*
		ParamOrdinal
		ParamName
		Status
		UserType
		Flags
		TypeInfo
		CryptoMetadata
		Value
	*/
	r.uint16()
	nv.Name = r.BVarChar()
	r.byte()
	r.uint32() // UserType (uint16 prior to 7.2)
	r.uint16()
	ti := readTypeInfo(r)
	nv.Value = ti.Reader(&ti, r)
	return
}

func processSingleResponse(ctx context.Context, sess *tdsSession, ch chan tokenStruct, outs outputs) {
	firstResult := true
	defer func() {
		if err := recover(); err != nil {
			if sess.logFlags&logErrors != 0 {
				sess.logger.Log(ctx, msdsn.LogErrors, fmt.Sprintf("Intercepted panic %v", err))
			}
			ch <- err
		}
		close(ch)
	}()

	packet_type, err := sess.buf.BeginRead()
	if err != nil {
		if sess.logFlags&logErrors != 0 {
			sess.logger.Log(ctx, msdsn.LogErrors, fmt.Sprintf("BeginRead failed %v", err))
		}
		ch <- err
		return
	}
	if packet_type != packReply {
		badStreamPanic(fmt.Errorf("unexpected packet type in reply: got %v, expected %v", packet_type, packReply))
	}
	var columns []columnStruct
	errs := make([]Error, 0, 5)
	for tokens := 0; ; tokens += 1 {
		token := token(sess.buf.byte())
		if sess.logFlags&logDebug != 0 {
			sess.logger.Log(ctx, msdsn.LogDebug, fmt.Sprintf("got token %v", token))
		}
		switch token {
		case tokenSSPI:
			ch <- parseSSPIMsg(sess.buf)
			return
		case tokenFedAuthInfo:
			ch <- parseFedAuthInfo(sess.buf)
			return
		case tokenReturnStatus:
			returnStatus := parseReturnStatus(sess.buf)
			ch <- returnStatus
		case tokenLoginAck:
			loginAck := parseLoginAck(sess.buf)
			ch <- loginAck
		case tokenFeatureExtAck:
			featureExtAck := parseFeatureExtAck(sess.buf)
			ch <- featureExtAck
		case tokenOrder:
			order := parseOrder(sess.buf)
			ch <- order
		case tokenDoneInProc:
			done := parseDoneInProc(sess.buf)

			ch <- done
			if sess.logFlags&logRows != 0 && done.Status&doneCount != 0 {
				sess.logger.Log(ctx, msdsn.LogRows, fmt.Sprintf("(%d rows affected)", done.RowCount))

				if outs.msgq != nil {
					_ = sqlexp.ReturnMessageEnqueue(ctx, outs.msgq, sqlexp.MsgRowsAffected{Count: int64(done.RowCount)})
				}
			}
			if done.Status&doneMore == 0 {
				if outs.msgq != nil {
					// For now we ignore ctx->Done errors that ReturnMessageEnqueue might return
					// It's not clear how to handle them correctly here, and data/sql seems
					// to set Rows.Err correctly when ctx expires already
					_ = sqlexp.ReturnMessageEnqueue(ctx, outs.msgq, sqlexp.MsgNextResultSet{})
				}
				return
			}
		case tokenDone, tokenDoneProc:
			done := parseDone(sess.buf)
			done.errors = errs
			if outs.msgq != nil {
				errs = make([]Error, 0, 5)
			}
			if sess.logFlags&logDebug != 0 {
				sess.logger.Log(ctx, msdsn.LogDebug, fmt.Sprintf("got DONE or DONEPROC status=%d", done.Status))
			}
			if done.Status&doneSrvError != 0 {
				ch <- ServerError{done.getError()}
				if outs.msgq != nil {
					_ = sqlexp.ReturnMessageEnqueue(ctx, outs.msgq, sqlexp.MsgNextResultSet{})
				}
				return
			}
			if sess.logFlags&logRows != 0 && done.Status&doneCount != 0 {
				sess.logger.Log(ctx, msdsn.LogRows, fmt.Sprintf("(%d row(s) affected)", done.RowCount))
			}
			ch <- done
			if done.Status&doneCount != 0 {
				if outs.msgq != nil {
					_ = sqlexp.ReturnMessageEnqueue(ctx, outs.msgq, sqlexp.MsgRowsAffected{Count: int64(done.RowCount)})
				}
			}
			if done.Status&doneMore == 0 {
				if outs.msgq != nil {
					_ = sqlexp.ReturnMessageEnqueue(ctx, outs.msgq, sqlexp.MsgNextResultSet{})
				}
				return
			}
		case tokenColMetadata:
			columns = parseColMetadata72(sess.buf)
			ch <- columns

			if outs.msgq != nil {
				if !firstResult {
					_ = sqlexp.ReturnMessageEnqueue(ctx, outs.msgq, sqlexp.MsgNextResultSet{})
				}
				_ = sqlexp.ReturnMessageEnqueue(ctx, outs.msgq, sqlexp.MsgNext{})
			}
			firstResult = false

		case tokenRow:
			row := make([]interface{}, len(columns))
			parseRow(sess.buf, columns, row)
			ch <- row
		case tokenNbcRow:
			row := make([]interface{}, len(columns))
			parseNbcRow(sess.buf, columns, row)
			ch <- row
		case tokenEnvChange:
			processEnvChg(ctx, sess)
		case tokenError:
			err := parseError72(sess.buf)
			if sess.logFlags&logDebug != 0 {
				sess.logger.Log(ctx, msdsn.LogDebug, fmt.Sprintf("got ERROR %d %s", err.Number, err.Message))
			}
			errs = append(errs, err)
			if sess.logFlags&logErrors != 0 {
				sess.logger.Log(ctx, msdsn.LogErrors, err.Message)
			}
			if outs.msgq != nil {
				_ = sqlexp.ReturnMessageEnqueue(ctx, outs.msgq, sqlexp.MsgError{Error: err})
			}
		case tokenInfo:
			info := parseInfo(sess.buf)
			if sess.logFlags&logDebug != 0 {
				sess.logger.Log(ctx, msdsn.LogDebug, fmt.Sprintf("got INFO %d %s", info.Number, info.Message))
			}
			if sess.logFlags&logMessages != 0 {
				sess.logger.Log(ctx, msdsn.LogMessages, info.Message)
			}
			if outs.msgq != nil {
				_ = sqlexp.ReturnMessageEnqueue(ctx, outs.msgq, sqlexp.MsgNotice{Message: info})
			}
		case tokenReturnValue:
			nv := parseReturnValue(sess.buf)
			if len(nv.Name) > 0 {
				name := nv.Name[1:] // Remove the leading "@".
				if ov, has := outs.params[name]; has {
					err = scanIntoOut(name, nv.Value, ov)
					if err != nil {
						fmt.Println("scan error", err)
						ch <- err
					}
				}
			}
		default:
			badStreamPanic(fmt.Errorf("unknown token type returned: %v", token))
		}
	}
}

type tokenProcessor struct {
	tokChan    chan tokenStruct
	ctx        context.Context
	sess       *tdsSession
	outs       outputs
	lastRow    []interface{}
	rowCount   int64
	firstError error
	// whether to skip sending attention when ctx is done
	noAttn bool
}

func startReading(sess *tdsSession, ctx context.Context, outs outputs) *tokenProcessor {
	tokChan := make(chan tokenStruct, 5)
	go processSingleResponse(ctx, sess, tokChan, outs)
	return &tokenProcessor{
		tokChan: tokChan,
		ctx:     ctx,
		sess:    sess,
		outs:    outs,
	}
}

func (t *tokenProcessor) iterateResponse() error {
	for {
		tok, err := t.nextToken()
		if err == nil {
			if tok == nil {
				return t.firstError
			} else {
				switch token := tok.(type) {
				case []columnStruct:
					t.sess.columns = token
				case []interface{}:
					t.lastRow = token
				case doneInProcStruct:
					if token.Status&doneCount != 0 {
						t.rowCount += int64(token.RowCount)
					}
				case doneStruct:
					if token.Status&doneCount != 0 {
						t.rowCount += int64(token.RowCount)
					}
					if token.isError() && t.firstError == nil {
						t.firstError = token.getError()
					}
				case ReturnStatus:
					if t.outs.returnStatus != nil {
						*t.outs.returnStatus = token
					}
					/*case error:
					if resultError == nil {
						resultError = token
					}*/
				}
			}
		} else {
			return err
		}
	}
}

func (t tokenProcessor) nextToken() (tokenStruct, error) {
	// we do this separate non-blocking check on token channel to
	// prioritize it over cancellation channel
	select {
	case tok, more := <-t.tokChan:
		err, more := tok.(error)
		if more {
			// this is an error and not a token
			return nil, err
		} else {
			return tok, nil
		}
	default:
		// there are no tokens on the channel, will need to wait
	}

	select {
	case tok, more := <-t.tokChan:
		if more {
			err, ok := tok.(error)
			if ok {
				// this is an error and not a token
				return nil, err
			} else {
				return tok, nil
			}
		} else {
			// completed reading response
			return nil, nil
		}
	case <-t.ctx.Done():
		// It seems the Message function on t.outs.msgq doesn't get the Done if it comes here instead
		if t.outs.msgq != nil {
			_ = sqlexp.ReturnMessageEnqueue(t.ctx, t.outs.msgq, sqlexp.MsgNextResultSet{})
		}
		if t.noAttn {
			return nil, t.ctx.Err()
		}
		if err := sendAttention(t.sess.buf); err != nil {
			// unable to send attention, current connection is bad
			// notify caller and close channel
			return nil, err
		}

		// now the server should send cancellation confirmation
		// it is possible that we already received full response
		// just before we sent cancellation request
		// in this case current response would not contain confirmation
		// and we would need to read one more response

		// first lets finish reading current response and look
		// for confirmation in it
		if readCancelConfirmation(t.tokChan) {
			// we got confirmation in current response
			return nil, t.ctx.Err()
		}
		// we did not get cancellation confirmation in the current response
		// read one more response, it must be there
		t.tokChan = make(chan tokenStruct, 5)
		go processSingleResponse(t.ctx, t.sess, t.tokChan, t.outs)
		if readCancelConfirmation(t.tokChan) {
			return nil, t.ctx.Err()
		}
		// we did not get cancellation confirmation, something is not
		// right, this connection is not usable anymore
		return nil, errors.New("did not get cancellation confirmation from the server")
	}
}

func readCancelConfirmation(tokChan chan tokenStruct) bool {
	for tok := range tokChan {
		switch tok := tok.(type) {
		default:
		// just skip token
		case doneStruct:
			if tok.Status&doneAttn != 0 {
				// got cancellation confirmation, exit
				return true
			}
		}
	}
	return false
}
