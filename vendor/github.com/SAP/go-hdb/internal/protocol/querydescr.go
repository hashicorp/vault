// SPDX-FileCopyrightText: 2014-2020 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

package protocol

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/SAP/go-hdb/internal/protocol/scanner"
)

// QueryKind is the query type of a database statement.
type QueryKind int

func (k QueryKind) String() string {
	keyword, ok := queryKindKeyword[k]
	if ok {
		return keyword
	}
	return fmt.Sprintf("cmdKind(%d)", k)
}

// Query kind constants.
const (
	QkUnknown QueryKind = iota
	QkCall
	QkSelect
	QkInsert
	QkUpdate
	QkUpsert
	QkCreate
	QkDrop
	QkSet
	QkID
)

var (
	queryKindKeyword = map[QueryKind]string{
		QkUnknown: "unknown",
		QkCall:    "call",
		QkSelect:  "select",
		QkInsert:  "insert",
		QkUpdate:  "update",
		QkUpsert:  "upsert",
		QkCreate:  "create",
		QkDrop:    "drop",
		QkSet:     "set",
		QkID:      "id",
	}
	queryKeywordKind = map[string]QueryKind{}
)

func init() {
	// build cmdKeywordKind from cmdKindKeyword
	for k, v := range queryKindKeyword {
		queryKeywordKind[v] = k
	}
}

func encodeID(id uint64) string {
	return fmt.Sprintf("%s %s", queryKindKeyword[QkID], strconv.FormatUint(id, 10))
}

var errInvalidCmdToken = errors.New("invalid command token")

const (
	bulkQuery = "bulk"
)

// QueryDescr represents a query descriptor of a database statement.
type QueryDescr struct {
	query  string
	kind   QueryKind
	isBulk bool
	id     uint64
}

func (d *QueryDescr) String() string {
	return fmt.Sprintf("query: %s kind: %s isBulk: %t", d.query, d.kind, d.isBulk)
}

// Query return the query statement of a query descriptor.
func (d *QueryDescr) Query() string { return d.query }

// Kind return the query kind of a query descriptor.
func (d *QueryDescr) Kind() QueryKind { return d.kind }

// ID return the query id of a query descriptor (legacy mode: call table output parameters).
func (d *QueryDescr) ID() uint64 { return d.id }

// IsBulk returns true if the query is a bulk statement..
func (d *QueryDescr) IsBulk() bool { return d.isBulk }

// NewQueryDescr returns a new QueryDescr instance.
func NewQueryDescr(query string, sc *scanner.Scanner) (*QueryDescr, error) {
	d := &QueryDescr{query: query}

	sc.Reset(query)

	// first token
	token, start, end := sc.Next()

	if token != scanner.Identifier {
		return nil, errInvalidCmdToken
	}

	if strings.ToLower(query[start:end]) == bulkQuery {
		d.isBulk = true
		_, start, end = sc.Next()
	}

	// kind
	keyword := strings.ToLower(query[start:end])

	d.kind = QkUnknown
	kind, ok := queryKeywordKind[keyword]
	if ok {
		d.kind = kind
	}

	// command
	d.query = query[start:] // cut off whitespaces and bulk

	// result set id query
	if d.kind == QkID {
		token, start, end = sc.Next()
		if token != scanner.Number {
			return nil, errInvalidCmdToken
		}
		var err error
		d.id, err = strconv.ParseUint(query[start:end], 10, 64)
		if err != nil {
			return nil, err
		}
	}

	// TODO release v1.0.0 - scan variables (named parameters)

	return d, nil
}
