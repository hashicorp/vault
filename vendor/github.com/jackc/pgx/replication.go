package pgx

import (
	"context"
	"encoding/binary"
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/jackc/pgx/pgio"
	"github.com/jackc/pgx/pgproto3"
)

const (
	copyBothResponse                  = 'W'
	walData                           = 'w'
	senderKeepalive                   = 'k'
	standbyStatusUpdate               = 'r'
	initialReplicationResponseTimeout = 5 * time.Second
)

var epochNano int64

func init() {
	epochNano = time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).UnixNano()
}

// Format the given 64bit LSN value into the XXX/XXX format,
// which is the format reported by postgres.
func FormatLSN(lsn uint64) string {
	return fmt.Sprintf("%X/%X", uint32(lsn>>32), uint32(lsn))
}

// Parse the given XXX/XXX format LSN as reported by postgres,
// into a 64 bit integer as used internally by the wire procotols
func ParseLSN(lsn string) (outputLsn uint64, err error) {
	var upperHalf uint64
	var lowerHalf uint64
	var nparsed int
	nparsed, err = fmt.Sscanf(lsn, "%X/%X", &upperHalf, &lowerHalf)
	if err != nil {
		return
	}

	if nparsed != 2 {
		err = errors.New(fmt.Sprintf("Failed to parsed LSN: %s", lsn))
		return
	}

	outputLsn = (upperHalf << 32) + lowerHalf
	return
}

// The WAL message contains WAL payload entry data
type WalMessage struct {
	// The WAL start position of this data. This
	// is the WAL position we need to track.
	WalStart uint64
	// The server wal end and server time are
	// documented to track the end position and current
	// time of the server, both of which appear to be
	// unimplemented in pg 9.5.
	ServerWalEnd uint64
	ServerTime   uint64
	// The WAL data is the raw unparsed binary WAL entry.
	// The contents of this are determined by the output
	// logical encoding plugin.
	WalData []byte
}

func (w *WalMessage) Time() time.Time {
	return time.Unix(0, (int64(w.ServerTime)*1000)+epochNano)
}

func (w *WalMessage) ByteLag() uint64 {
	return (w.ServerWalEnd - w.WalStart)
}

func (w *WalMessage) String() string {
	return fmt.Sprintf("Wal: %s Time: %s Lag: %d", FormatLSN(w.WalStart), w.Time(), w.ByteLag())
}

// The server heartbeat is sent periodically from the server,
// including server status, and a reply request field
type ServerHeartbeat struct {
	// The current max wal position on the server,
	// used for lag tracking
	ServerWalEnd uint64
	// The server time, in microseconds since jan 1 2000
	ServerTime uint64
	// If 1, the server is requesting a standby status message
	// to be sent immediately.
	ReplyRequested byte
}

func (s *ServerHeartbeat) Time() time.Time {
	return time.Unix(0, (int64(s.ServerTime)*1000)+epochNano)
}

func (s *ServerHeartbeat) String() string {
	return fmt.Sprintf("WalEnd: %s ReplyRequested: %d T: %s", FormatLSN(s.ServerWalEnd), s.ReplyRequested, s.Time())
}

// The replication message wraps all possible messages from the
// server received during replication. At most one of the wal message
// or server heartbeat will be non-nil
type ReplicationMessage struct {
	WalMessage      *WalMessage
	ServerHeartbeat *ServerHeartbeat
}

// The standby status is the client side heartbeat sent to the postgresql
// server to track the client wal positions. For practical purposes,
// all wal positions are typically set to the same value.
type StandbyStatus struct {
	// The WAL position that's been locally written
	WalWritePosition uint64
	// The WAL position that's been locally flushed
	WalFlushPosition uint64
	// The WAL position that's been locally applied
	WalApplyPosition uint64
	// The client time in microseconds since jan 1 2000
	ClientTime uint64
	// If 1, requests the server to immediately send a
	// server heartbeat
	ReplyRequested byte
}

// Create a standby status struct, which sets all the WAL positions
// to the given wal position, and the client time to the current time.
// The wal positions are, in order:
// WalFlushPosition
// WalApplyPosition
// WalWritePosition
//
// If only one position is provided, it will be used as the value for all 3
// status fields. Note you must provide either 1 wal position, or all 3
// in order to initialize the standby status.
func NewStandbyStatus(walPositions ...uint64) (status *StandbyStatus, err error) {
	if len(walPositions) == 1 {
		status = new(StandbyStatus)
		status.WalFlushPosition = walPositions[0]
		status.WalApplyPosition = walPositions[0]
		status.WalWritePosition = walPositions[0]
	} else if len(walPositions) == 3 {
		status = new(StandbyStatus)
		status.WalFlushPosition = walPositions[0]
		status.WalApplyPosition = walPositions[1]
		status.WalWritePosition = walPositions[2]
	} else {
		err = errors.New(fmt.Sprintf("Invalid number of wal positions provided, need 1 or 3, got %d", len(walPositions)))
		return
	}
	status.ClientTime = uint64((time.Now().UnixNano() - epochNano) / 1000)
	return
}

func ReplicationConnect(config ConnConfig) (r *ReplicationConn, err error) {
	if config.RuntimeParams == nil {
		config.RuntimeParams = make(map[string]string)
	}
	config.RuntimeParams["replication"] = "database"

	c, err := Connect(config)
	if err != nil {
		return
	}
	return &ReplicationConn{c: c}, nil
}

type ReplicationConn struct {
	c *Conn
}

// Send standby status to the server, which both acts as a keepalive
// message to the server, as well as carries the WAL position of the
// client, which then updates the server's replication slot position.
func (rc *ReplicationConn) SendStandbyStatus(k *StandbyStatus) (err error) {
	buf := rc.c.wbuf
	buf = append(buf, copyData)
	sp := len(buf)
	buf = pgio.AppendInt32(buf, -1)

	buf = append(buf, standbyStatusUpdate)
	buf = pgio.AppendInt64(buf, int64(k.WalWritePosition))
	buf = pgio.AppendInt64(buf, int64(k.WalFlushPosition))
	buf = pgio.AppendInt64(buf, int64(k.WalApplyPosition))
	buf = pgio.AppendInt64(buf, int64(k.ClientTime))
	buf = append(buf, k.ReplyRequested)

	pgio.SetInt32(buf[sp:], int32(len(buf[sp:])))

	_, err = rc.c.conn.Write(buf)
	if err != nil {
		rc.c.die(err)
	}

	return
}

func (rc *ReplicationConn) Close() error {
	return rc.c.Close()
}

func (rc *ReplicationConn) IsAlive() bool {
	return rc.c.IsAlive()
}

func (rc *ReplicationConn) CauseOfDeath() error {
	return rc.c.CauseOfDeath()
}

func (rc *ReplicationConn) readReplicationMessage() (r *ReplicationMessage, err error) {
	msg, err := rc.c.rxMsg()
	if err != nil {
		return
	}

	switch msg := msg.(type) {
	case *pgproto3.NoticeResponse:
		pgError := rc.c.rxErrorResponse((*pgproto3.ErrorResponse)(msg))
		if rc.c.shouldLog(LogLevelInfo) {
			rc.c.log(LogLevelInfo, pgError.Error(), nil)
		}
	case *pgproto3.ErrorResponse:
		err = rc.c.rxErrorResponse(msg)
		if rc.c.shouldLog(LogLevelError) {
			rc.c.log(LogLevelError, err.Error(), nil)
		}
		return
	case *pgproto3.CopyBothResponse:
		// This is the tail end of the replication process start,
		// and can be safely ignored
		return
	case *pgproto3.CopyData:
		msgType := msg.Data[0]
		rp := 1

		switch msgType {
		case walData:
			walStart := binary.BigEndian.Uint64(msg.Data[rp:])
			rp += 8
			serverWalEnd := binary.BigEndian.Uint64(msg.Data[rp:])
			rp += 8
			serverTime := binary.BigEndian.Uint64(msg.Data[rp:])
			rp += 8
			walData := msg.Data[rp:]
			walMessage := WalMessage{WalStart: walStart,
				ServerWalEnd: serverWalEnd,
				ServerTime:   serverTime,
				WalData:      walData,
			}

			return &ReplicationMessage{WalMessage: &walMessage}, nil
		case senderKeepalive:
			serverWalEnd := binary.BigEndian.Uint64(msg.Data[rp:])
			rp += 8
			serverTime := binary.BigEndian.Uint64(msg.Data[rp:])
			rp += 8
			replyNow := msg.Data[rp]
			rp += 1
			h := &ServerHeartbeat{ServerWalEnd: serverWalEnd, ServerTime: serverTime, ReplyRequested: replyNow}
			return &ReplicationMessage{ServerHeartbeat: h}, nil
		default:
			if rc.c.shouldLog(LogLevelError) {
				rc.c.log(LogLevelError, "Unexpected data playload message type", map[string]interface{}{"type": msgType})
			}
		}
	default:
		if rc.c.shouldLog(LogLevelError) {
			rc.c.log(LogLevelError, "Unexpected replication message type", map[string]interface{}{"type": msg})
		}
	}
	return
}

// Wait for a single replication message.
//
// Properly using this requires some knowledge of the postgres replication mechanisms,
// as the client can receive both WAL data (the ultimate payload) and server heartbeat
// updates. The caller also must send standby status updates in order to keep the connection
// alive and working.
//
// This returns the context error when there is no replication message before
// the context is canceled.
func (rc *ReplicationConn) WaitForReplicationMessage(ctx context.Context) (*ReplicationMessage, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	go func() {
		select {
		case <-ctx.Done():
			if err := rc.c.conn.SetDeadline(time.Now()); err != nil {
				rc.Close() // Close connection if unable to set deadline
				return
			}
			rc.c.closedChan <- ctx.Err()
		case <-rc.c.doneChan:
		}
	}()

	r, opErr := rc.readReplicationMessage()

	var err error
	select {
	case err = <-rc.c.closedChan:
		if err := rc.c.conn.SetDeadline(time.Time{}); err != nil {
			rc.Close() // Close connection if unable to disable deadline
			return nil, err
		}

		if opErr == nil {
			err = nil
		}
	case rc.c.doneChan <- struct{}{}:
		err = opErr
	}

	return r, err
}

func (rc *ReplicationConn) sendReplicationModeQuery(sql string) (*Rows, error) {
	rc.c.lastActivityTime = time.Now()

	rows := rc.c.getRows(sql, nil)

	if err := rc.c.lock(); err != nil {
		rows.fatal(err)
		return rows, err
	}
	rows.unlockConn = true

	err := rc.c.sendSimpleQuery(sql)
	if err != nil {
		rows.fatal(err)
	}

	msg, err := rc.c.rxMsg()
	if err != nil {
		return nil, err
	}

	switch msg := msg.(type) {
	case *pgproto3.RowDescription:
		rows.fields = rc.c.rxRowDescription(msg)
		// We don't have c.PgTypes here because we're a replication
		// connection. This means the field descriptions will have
		// only OIDs. Not much we can do about this.
	default:
		if e := rc.c.processContextFreeMsg(msg); e != nil {
			rows.fatal(e)
			return rows, e
		}
	}

	return rows, rows.err
}

// Execute the "IDENTIFY_SYSTEM" command as documented here:
// https://www.postgresql.org/docs/9.5/static/protocol-replication.html
//
// This will return (if successful) a result set that has a single row
// that contains the systemid, current timeline, xlogpos and database
// name.
//
// NOTE: Because this is a replication mode connection, we don't have
// type names, so the field descriptions in the result will have only
// OIDs and no DataTypeName values
func (rc *ReplicationConn) IdentifySystem() (r *Rows, err error) {
	return rc.sendReplicationModeQuery("IDENTIFY_SYSTEM")
}

// Execute the "TIMELINE_HISTORY" command as documented here:
// https://www.postgresql.org/docs/9.5/static/protocol-replication.html
//
// This will return (if successful) a result set that has a single row
// that contains the filename of the history file and the content
// of the history file. If called for timeline 1, typically this will
// generate an error that the timeline history file does not exist.
//
// NOTE: Because this is a replication mode connection, we don't have
// type names, so the field descriptions in the result will have only
// OIDs and no DataTypeName values
func (rc *ReplicationConn) TimelineHistory(timeline int) (r *Rows, err error) {
	return rc.sendReplicationModeQuery(fmt.Sprintf("TIMELINE_HISTORY %d", timeline))
}

// Start a replication connection, sending WAL data to the given replication
// receiver. This function wraps a START_REPLICATION command as documented
// here:
// https://www.postgresql.org/docs/9.5/static/protocol-replication.html
//
// Once started, the client needs to invoke WaitForReplicationMessage() in order
// to fetch the WAL and standby status. Also, it is the responsibility of the caller
// to periodically send StandbyStatus messages to update the replication slot position.
//
// This function assumes that slotName has already been created. In order to omit the timeline argument
// pass a -1 for the timeline to get the server default behavior.
func (rc *ReplicationConn) StartReplication(slotName string, startLsn uint64, timeline int64, pluginArguments ...string) (err error) {
	queryString := fmt.Sprintf("START_REPLICATION SLOT %s LOGICAL %s", slotName, FormatLSN(startLsn))
	if timeline >= 0 {
		timelineOption := fmt.Sprintf("TIMELINE %d", timeline)
		pluginArguments = append(pluginArguments, timelineOption)
	}

	if len(pluginArguments) > 0 {
		queryString += fmt.Sprintf(" ( %s )", strings.Join(pluginArguments, ", "))
	}

	if err = rc.c.sendQuery(queryString); err != nil {
		return
	}

	ctx, cancelFn := context.WithTimeout(context.Background(), initialReplicationResponseTimeout)
	defer cancelFn()

	// The first replication message that comes back here will be (in a success case)
	// a empty CopyBoth that is (apparently) sent as the confirmation that the replication has
	// started. This call will either return nil, nil or if it returns an error
	// that indicates the start replication command failed
	var r *ReplicationMessage
	r, err = rc.WaitForReplicationMessage(ctx)
	if err != nil && r != nil {
		if rc.c.shouldLog(LogLevelError) {
			rc.c.log(LogLevelError, "Unexpected replication message", map[string]interface{}{"msg": r, "err": err})
		}
	}

	return
}

// Create the replication slot, using the given name and output plugin.
func (rc *ReplicationConn) CreateReplicationSlot(slotName, outputPlugin string) (err error) {
	_, err = rc.c.Exec(fmt.Sprintf("CREATE_REPLICATION_SLOT %s LOGICAL %s", slotName, outputPlugin))
	return
}

// Create the replication slot, using the given name and output plugin, and return the consistent_point and snapshot_name values.
func (rc *ReplicationConn) CreateReplicationSlotEx(slotName, outputPlugin string) (consistentPoint string, snapshotName string, err error) {
	var dummy string
	var rows *Rows
	rows, err = rc.sendReplicationModeQuery(fmt.Sprintf("CREATE_REPLICATION_SLOT %s LOGICAL %s", slotName, outputPlugin))
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&dummy, &consistentPoint, &snapshotName, &dummy)
	}
	return
}

// Drop the replication slot for the given name
func (rc *ReplicationConn) DropReplicationSlot(slotName string) (err error) {
	_, err = rc.c.Exec(fmt.Sprintf("DROP_REPLICATION_SLOT %s", slotName))
	return
}
