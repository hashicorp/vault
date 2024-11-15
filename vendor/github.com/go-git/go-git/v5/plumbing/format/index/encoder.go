package index

import (
	"bytes"
	"errors"
	"io"
	"sort"
	"time"

	"github.com/go-git/go-git/v5/plumbing/hash"
	"github.com/go-git/go-git/v5/utils/binary"
)

var (
	// EncodeVersionSupported is the range of supported index versions
	EncodeVersionSupported uint32 = 3

	// ErrInvalidTimestamp is returned by Encode if a Index with a Entry with
	// negative timestamp values
	ErrInvalidTimestamp = errors.New("negative timestamps are not allowed")
)

// An Encoder writes an Index to an output stream.
type Encoder struct {
	w    io.Writer
	hash hash.Hash
}

// NewEncoder returns a new encoder that writes to w.
func NewEncoder(w io.Writer) *Encoder {
	h := hash.New(hash.CryptoType)
	mw := io.MultiWriter(w, h)
	return &Encoder{mw, h}
}

// Encode writes the Index to the stream of the encoder.
func (e *Encoder) Encode(idx *Index) error {
	// TODO: support v4
	// TODO: support extensions
	if idx.Version > EncodeVersionSupported {
		return ErrUnsupportedVersion
	}

	if err := e.encodeHeader(idx); err != nil {
		return err
	}

	if err := e.encodeEntries(idx); err != nil {
		return err
	}

	return e.encodeFooter()
}

func (e *Encoder) encodeHeader(idx *Index) error {
	return binary.Write(e.w,
		indexSignature,
		idx.Version,
		uint32(len(idx.Entries)),
	)
}

func (e *Encoder) encodeEntries(idx *Index) error {
	sort.Sort(byName(idx.Entries))

	for _, entry := range idx.Entries {
		if err := e.encodeEntry(entry); err != nil {
			return err
		}
		entryLength := entryHeaderLength
		if entry.IntentToAdd || entry.SkipWorktree {
			entryLength += 2
		}

		wrote := entryLength + len(entry.Name)
		if err := e.padEntry(wrote); err != nil {
			return err
		}
	}

	return nil
}

func (e *Encoder) encodeEntry(entry *Entry) error {
	sec, nsec, err := e.timeToUint32(&entry.CreatedAt)
	if err != nil {
		return err
	}

	msec, mnsec, err := e.timeToUint32(&entry.ModifiedAt)
	if err != nil {
		return err
	}

	flags := uint16(entry.Stage&0x3) << 12
	if l := len(entry.Name); l < nameMask {
		flags |= uint16(l)
	} else {
		flags |= nameMask
	}

	flow := []interface{}{
		sec, nsec,
		msec, mnsec,
		entry.Dev,
		entry.Inode,
		entry.Mode,
		entry.UID,
		entry.GID,
		entry.Size,
		entry.Hash[:],
	}

	flagsFlow := []interface{}{flags}

	if entry.IntentToAdd || entry.SkipWorktree {
		var extendedFlags uint16

		if entry.IntentToAdd {
			extendedFlags |= intentToAddMask
		}
		if entry.SkipWorktree {
			extendedFlags |= skipWorkTreeMask
		}

		flagsFlow = []interface{}{flags | entryExtended, extendedFlags}
	}

	flow = append(flow, flagsFlow...)

	if err := binary.Write(e.w, flow...); err != nil {
		return err
	}

	return binary.Write(e.w, []byte(entry.Name))
}

func (e *Encoder) timeToUint32(t *time.Time) (uint32, uint32, error) {
	if t.IsZero() {
		return 0, 0, nil
	}

	if t.Unix() < 0 || t.UnixNano() < 0 {
		return 0, 0, ErrInvalidTimestamp
	}

	return uint32(t.Unix()), uint32(t.Nanosecond()), nil
}

func (e *Encoder) padEntry(wrote int) error {
	padLen := 8 - wrote%8

	_, err := e.w.Write(bytes.Repeat([]byte{'\x00'}, padLen))
	return err
}

func (e *Encoder) encodeFooter() error {
	return binary.Write(e.w, e.hash.Sum(nil))
}

type byName []*Entry

func (l byName) Len() int           { return len(l) }
func (l byName) Swap(i, j int)      { l[i], l[j] = l[j], l[i] }
func (l byName) Less(i, j int) bool { return l[i].Name < l[j].Name }
