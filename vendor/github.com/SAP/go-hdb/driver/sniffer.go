package driver

import (
	"context"
	"errors"
	"io"
	"log"
	"log/slog"
	"net"
	"sync"

	p "github.com/SAP/go-hdb/driver/internal/protocol"
	"github.com/SAP/go-hdb/driver/internal/protocol/encoding"
	"github.com/SAP/go-hdb/driver/unicode/cesu8"
)

// A Sniffer is a simple proxy for logging hdb protocol requests and responses.
type Sniffer struct {
	logger *slog.Logger
	conn   net.Conn
	dbConn net.Conn
}

// NewSniffer creates a new sniffer instance. The conn parameter is the net.Conn connection, where the Sniffer
// is listening for hdb protocol calls. The dbAddr is the hdb host port address in "host:port" format.
func NewSniffer(conn net.Conn, dbConn net.Conn) *Sniffer {
	return &Sniffer{
		logger: slog.Default().With(slog.String("conn", conn.RemoteAddr().String())),
		conn:   conn,
		dbConn: dbConn,
	}
}

func pipeData(wg *sync.WaitGroup, conn net.Conn, dbConn net.Conn, wr io.Writer) {
	defer wg.Done()

	mwr := io.MultiWriter(dbConn, wr)
	trd := io.TeeReader(conn, mwr)
	buf := make([]byte, 1000)

	var err error
	for err == nil {
		_, err = trd.Read(buf)
	}
}

func readMsg(ctx context.Context, prd *p.Reader) error {
	// TODO complete for non generic parts, see internal/protocol/parts/newGenPartReader for details
	_, err := prd.IterateParts(ctx, 0, nil)
	// _, err := prd.IterateParts(ctx, 0, func(kind p.PartKind, attrs p.PartAttributes, read func(part p.Part)) {})
	return err
}

func logData(ctx context.Context, wg *sync.WaitGroup, prd *p.Reader) {
	defer wg.Done()

	if err := prd.ReadProlog(ctx); err != nil {
		panic(err)
	}

	var err error
	for !errors.Is(err, io.EOF) {
		err = readMsg(ctx, prd)
	}
}

// Run starts the protocol request and response logging.
func (s *Sniffer) Run() error {
	clientRd, clientWr := io.Pipe()
	dbRd, dbWr := io.Pipe()

	ctx := context.Background()
	wg := &sync.WaitGroup{}

	wg.Add(4)
	go pipeData(wg, s.conn, s.dbConn, clientWr)
	go pipeData(wg, s.dbConn, s.conn, dbWr)

	defaultDecoder := cesu8.DefaultDecoder()

	clientDec := encoding.NewDecoder(clientRd, defaultDecoder, false)
	dbDec := encoding.NewDecoder(dbRd, defaultDecoder, false)

	pClientRd := p.NewClientReader(clientDec, true, s.logger, defaultLobChunkSize)
	pDBRd := p.NewDBReader(dbDec, true, s.logger, defaultLobChunkSize)

	go logData(ctx, wg, pClientRd)
	go logData(ctx, wg, pDBRd)

	wg.Wait()
	log.Println("end run")

	return nil
}
