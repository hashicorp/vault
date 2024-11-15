// Copyright (c) 2017-2022 Snowflake Computing Inc. All rights reserved.

package gosnowflake

import (
	"context"
	"database/sql/driver"
	"errors"
)

type snowflakeTx struct {
	sc  *snowflakeConn
	ctx context.Context
}

type txCommand int

const (
	commit txCommand = iota
	rollback
)

func (cmd txCommand) string() (string, error) {
	switch cmd {
	case commit:
		return "COMMIT", nil
	case rollback:
		return "ROLLBACK", nil
	}
	return "", errors.New("unsupported transaction command")
}

func (tx *snowflakeTx) Commit() error {
	return tx.execTxCommand(commit)
}

func (tx *snowflakeTx) Rollback() error {
	return tx.execTxCommand(rollback)
}

func (tx *snowflakeTx) execTxCommand(command txCommand) (err error) {
	txStr, err := command.string()
	if err != nil {
		return
	}
	if tx.sc == nil || tx.sc.rest == nil {
		return driver.ErrBadConn
	}
	_, err = tx.sc.exec(tx.ctx, txStr, false /* noResult */, false /* isInternal */, false /* describeOnly */, nil)
	if err != nil {
		return
	}
	tx.sc = nil
	return
}
