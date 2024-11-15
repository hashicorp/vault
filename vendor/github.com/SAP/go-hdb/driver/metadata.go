package driver

import (
	"context"
	"reflect"

	p "github.com/SAP/go-hdb/driver/internal/protocol"
)

// ColumnType equals sql.ColumnType.
type ColumnType interface {
	DatabaseTypeName() string
	DecimalSize() (precision, scale int64, ok bool)
	Length() (length int64, ok bool)
	Name() string
	Nullable() (nullable bool, ok bool)
	ScanType() reflect.Type
}

// ParameterType extends ColumnType with stored procedure metadata.
type ParameterType interface {
	ColumnType
	In() bool
	Out() bool
	InOut() bool
}

// StmtMetadata provides access to the parameter and result metadata of a prepared statement.
type StmtMetadata interface {
	ParameterTypes() []ParameterType
	ColumnTypes() []ColumnType
}

// use unexported type to avoid key collisions.
type stmtMetadataCtxKeyType struct{}

var stmtMetadataCtxKey stmtMetadataCtxKeyType

// WithStmtMetadata can be used to add a statement metadata reference to the context used for a Prepare call.
// The Prepare call will set the stmtMetadata reference on successful preparation.
func WithStmtMetadata(ctx context.Context, stmtMetadata *StmtMetadata) context.Context {
	return context.WithValue(ctx, stmtMetadataCtxKey, stmtMetadata)
}

var (
	_ StmtMetadata  = (*prepareResult)(nil)
	_ ParameterType = (*p.ParameterField)(nil)
	_ ColumnType    = (*p.ResultField)(nil)
)
