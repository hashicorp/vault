package gocbcore

import (
	"encoding/json"
	"errors"
	"io"
)

type rowStreamState int

const (
	rowStreamStateStart    rowStreamState = 0
	rowStreamStateRows     rowStreamState = 1
	rowStreamStatePostRows rowStreamState = 2
	rowStreamStateEnd      rowStreamState = 3
)

type rowStreamer struct {
	decoder    *json.Decoder
	rowsAttrib string
	attribs    map[string]json.RawMessage
	state      rowStreamState
}

func newRowStreamer(stream io.Reader, rowsAttrib string) (*rowStreamer, error) {
	decoder := json.NewDecoder(stream)

	streamer := &rowStreamer{
		decoder:    decoder,
		rowsAttrib: rowsAttrib,
		attribs:    make(map[string]json.RawMessage),
		state:      rowStreamStateStart,
	}

	if err := streamer.begin(); err != nil {
		return nil, err
	}

	return streamer, nil
}

func (s *rowStreamer) begin() error {
	if s.state != rowStreamStateStart {
		return errors.New("unexpected parsing state during begin")
	}

	// Read the opening { for the result
	t, err := s.decoder.Token()
	if err != nil {
		return err
	}
	if delim, ok := t.(json.Delim); !ok || delim != '{' {
		return errors.New("expected an opening brace for the result")
	}

	for {
		if !s.decoder.More() {
			// We reached the end of the object
			s.state = rowStreamStateEnd
			break
		}

		// Read the attribute name
		t, err = s.decoder.Token()
		if err != nil {
			return err
		}
		key, keyOk := t.(string)
		if !keyOk {
			return errors.New("expected an object property name")
		}

		if key == s.rowsAttrib {
			// Read the opening [ for the rows
			t, err = s.decoder.Token()
			if err != nil {
				return err
			}

			if t == nil {
				continue
			}

			if delim, ok := t.(json.Delim); !ok || delim != '[' {
				return errors.New("expected an opening bracket for the rows")
			}

			s.state = rowStreamStateRows
			break
		}

		// Read the attribute value
		var value json.RawMessage
		err = s.decoder.Decode(&value)
		if err != nil {
			return err
		}

		// Save the attribute for the meta-data
		s.attribs[key] = value
	}

	return nil
}

func (s *rowStreamer) readRow() (json.RawMessage, error) {
	if s.state < rowStreamStateRows {
		return nil, errors.New("unexpected parsing state during readRow")
	}

	// If we've already read all rows or rows is null, we return nil
	if s.state > rowStreamStateRows {
		return nil, nil
	}

	// If there are no more rows, mark the rows finished and
	// return nil to signal that we are at the end
	if !s.decoder.More() {
		s.state = rowStreamStatePostRows
		return nil, nil
	}

	// Decode this row and return a raw message
	var msg json.RawMessage
	err := s.decoder.Decode(&msg)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

func (s *rowStreamer) end() error {
	if s.state < rowStreamStatePostRows {
		return errors.New("unexpected parsing state during end")
	}

	// Check if we've already read everything
	if s.state > rowStreamStatePostRows {
		return nil
	}

	// Read the ending ] for the rows
	t, err := s.decoder.Token()
	if err != nil {
		return err
	}
	if delim, ok := t.(json.Delim); !ok || delim != ']' {
		return errors.New("expected an ending bracket for the rows")
	}

	for {
		if !s.decoder.More() {
			// We reached the end of the object
			s.state = rowStreamStateEnd
			break
		}

		// Read the attribute name
		t, err := s.decoder.Token()
		if err != nil {
			return err
		}

		key, keyOk := t.(string)
		if !keyOk {
			return errors.New("expected an object property name")
		}

		// Read the attribute value
		var value json.RawMessage
		err = s.decoder.Decode(&value)
		if err != nil {
			return err
		}

		// Save the attribute for the meta-data
		s.attribs[key] = value
	}

	return nil
}

func (s *rowStreamer) NextRowBytes() (json.RawMessage, error) {
	return s.readRow()
}

func (s *rowStreamer) Finalize() (json.RawMessage, error) {
	// Make sure we've read until the end of the object
	for {
		row, err := s.readRow()
		if err != nil {
			return nil, err
		}

		if row == nil {
			break
		}
	}

	// Read the rest of the result object
	err := s.end()
	if err != nil {
		return nil, err
	}

	// Reconstruct the non-rows JSON to a raw message
	metaBytes, err := json.Marshal(s.attribs)
	if err != nil {
		return nil, err
	}

	return json.RawMessage(metaBytes), nil
}

func (s *rowStreamer) EarlyAttrib(key string) json.RawMessage {
	val, ok := s.attribs[key]
	if !ok {
		return nil
	}

	return val
}
