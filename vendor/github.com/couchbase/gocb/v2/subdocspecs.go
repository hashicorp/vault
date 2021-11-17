package gocb

import "github.com/couchbase/gocbcore/v10/memd"

// LookupInSpec is the representation of an operation available when calling LookupIn
type LookupInSpec struct {
	op      memd.SubDocOpType
	path    string
	isXattr bool
}

// MutateInSpec is the representation of an operation available when calling MutateIn
type MutateInSpec struct {
	op         memd.SubDocOpType
	createPath bool
	isXattr    bool
	path       string
	value      interface{}
	multiValue bool
}

// GetSpecOptions are the options available to LookupIn subdoc Get operations.
type GetSpecOptions struct {
	IsXattr bool
}

// GetSpec indicates a path to be retrieved from the document.  The value of the path
// can later be retrieved from the LookupResult.
// The path syntax follows query's path syntax (e.g. `foo.bar.baz`).
func GetSpec(path string, opts *GetSpecOptions) LookupInSpec {
	if opts == nil {
		opts = &GetSpecOptions{}
	}

	return LookupInSpec{
		op:      memd.SubDocOpGet,
		path:    path,
		isXattr: opts.IsXattr,
	}
}

// ExistsSpecOptions are the options available to LookupIn subdoc Exists operations.
type ExistsSpecOptions struct {
	IsXattr bool
}

// ExistsSpec is similar to Path(), but does not actually retrieve the value from the server.
// This may save bandwidth if you only need to check for the existence of a
// path (without caring for its content). You can check the status of this
// operation by using .ContentAt (and ignoring the value) or .Exists() on the LookupResult.
func ExistsSpec(path string, opts *ExistsSpecOptions) LookupInSpec {
	if opts == nil {
		opts = &ExistsSpecOptions{}
	}

	return LookupInSpec{
		op:      memd.SubDocOpExists,
		path:    path,
		isXattr: opts.IsXattr,
	}
}

// CountSpecOptions are the options available to LookupIn subdoc Count operations.
type CountSpecOptions struct {
	IsXattr bool
}

// CountSpec allows you to retrieve the number of items in an array or keys within an
// dictionary within an element of a document.
func CountSpec(path string, opts *CountSpecOptions) LookupInSpec {
	if opts == nil {
		opts = &CountSpecOptions{}
	}

	return LookupInSpec{
		op:      memd.SubDocOpGetCount,
		path:    path,
		isXattr: opts.IsXattr,
	}
}

// InsertSpecOptions are the options available to subdocument Insert operations.
type InsertSpecOptions struct {
	CreatePath bool
	IsXattr    bool
}

// InsertSpec inserts a value at the specified path within the document.
func InsertSpec(path string, val interface{}, opts *InsertSpecOptions) MutateInSpec {
	if opts == nil {
		opts = &InsertSpecOptions{}
	}

	return MutateInSpec{
		op:         memd.SubDocOpDictAdd,
		createPath: opts.CreatePath,
		isXattr:    opts.IsXattr,
		path:       path,
		value:      val,
		multiValue: false,
	}
}

// UpsertSpecOptions are the options available to subdocument Upsert operations.
type UpsertSpecOptions struct {
	CreatePath bool
	IsXattr    bool
}

// UpsertSpec creates a new value at the specified path within the document if it does not exist, if it does exist then it
// updates it.
func UpsertSpec(path string, val interface{}, opts *UpsertSpecOptions) MutateInSpec {
	if opts == nil {
		opts = &UpsertSpecOptions{}
	}

	return MutateInSpec{
		op:         memd.SubDocOpDictSet,
		createPath: opts.CreatePath,
		isXattr:    opts.IsXattr,
		path:       path,
		value:      val,
		multiValue: false,
	}
}

// ReplaceSpecOptions are the options available to subdocument Replace operations.
type ReplaceSpecOptions struct {
	IsXattr bool
}

// ReplaceSpec replaces the value of the field at path.
func ReplaceSpec(path string, val interface{}, opts *ReplaceSpecOptions) MutateInSpec {
	if opts == nil {
		opts = &ReplaceSpecOptions{}
	}

	return MutateInSpec{
		op:         memd.SubDocOpReplace,
		createPath: false,
		isXattr:    opts.IsXattr,
		path:       path,
		value:      val,
		multiValue: false,
	}
}

// RemoveSpecOptions are the options available to subdocument Remove operations.
type RemoveSpecOptions struct {
	IsXattr bool
}

// RemoveSpec removes the field at path.
func RemoveSpec(path string, opts *RemoveSpecOptions) MutateInSpec {
	if opts == nil {
		opts = &RemoveSpecOptions{}
	}

	return MutateInSpec{
		op:         memd.SubDocOpDelete,
		createPath: false,
		isXattr:    opts.IsXattr,
		path:       path,
		value:      nil,
		multiValue: false,
	}
}

// ArrayAppendSpecOptions are the options available to subdocument ArrayAppend operations.
type ArrayAppendSpecOptions struct {
	CreatePath bool
	IsXattr    bool
	// HasMultiple adds multiple values as elements to an array.
	// When used `value` in the spec must be an array type
	// ArrayAppend("path", []int{1,2,3,4}, ArrayAppendSpecOptions{HasMultiple:true}) =>
	//   "path" [..., 1,2,3,4]
	//
	// This is a more efficient version (at both the network and server levels)
	// of doing
	// spec.ArrayAppend("path", 1, nil)
	// spec.ArrayAppend("path", 2, nil)
	// spec.ArrayAppend("path", 3, nil)
	HasMultiple bool
}

// ArrayAppendSpec adds an element(s) to the end (i.e. right) of an array
func ArrayAppendSpec(path string, val interface{}, opts *ArrayAppendSpecOptions) MutateInSpec {
	if opts == nil {
		opts = &ArrayAppendSpecOptions{}
	}

	return MutateInSpec{
		op:         memd.SubDocOpArrayPushLast,
		createPath: opts.CreatePath,
		isXattr:    opts.IsXattr,
		path:       path,
		value:      val,
		multiValue: opts.HasMultiple,
	}
}

// ArrayPrependSpecOptions are the options available to subdocument ArrayPrepend operations.
type ArrayPrependSpecOptions struct {
	CreatePath bool
	IsXattr    bool
	// HasMultiple adds multiple values as elements to an array.
	// When used `value` in the spec must be an array type
	// ArrayPrepend("path", []int{1,2,3,4}, ArrayPrependSpecOptions{HasMultiple:true}) =>
	//   "path" [1,2,3,4, ....]
	//
	// This is a more efficient version (at both the network and server levels)
	// of doing
	// spec.ArrayPrepend("path", 1, nil)
	// spec.ArrayPrepend("path", 2, nil)
	// spec.ArrayPrepend("path", 3, nil)
	HasMultiple bool
}

// ArrayPrependSpec adds an element to the beginning (i.e. left) of an array
func ArrayPrependSpec(path string, val interface{}, opts *ArrayPrependSpecOptions) MutateInSpec {
	if opts == nil {
		opts = &ArrayPrependSpecOptions{}
	}

	return MutateInSpec{
		op:         memd.SubDocOpArrayPushFirst,
		createPath: opts.CreatePath,
		isXattr:    opts.IsXattr,
		path:       path,
		value:      val,
		multiValue: opts.HasMultiple,
	}
}

// ArrayInsertSpecOptions are the options available to subdocument ArrayInsert operations.
type ArrayInsertSpecOptions struct {
	CreatePath bool
	IsXattr    bool
	// HasMultiple adds multiple values as elements to an array.
	// When used `value` in the spec must be an array type
	// ArrayInsert("path[1]", []int{1,2,3,4}, ArrayInsertSpecOptions{HasMultiple:true}) =>
	//   "path" [..., 1,2,3,4]
	//
	// This is a more efficient version (at both the network and server levels)
	// of doing
	// spec.ArrayInsert("path[2]", 1, nil)
	// spec.ArrayInsert("path[3]", 2, nil)
	// spec.ArrayInsert("path[4]", 3, nil)
	HasMultiple bool
}

// ArrayInsertSpec inserts an element at a given position within an array. The position should be
// specified as part of the path, e.g. path.to.array[3]
func ArrayInsertSpec(path string, val interface{}, opts *ArrayInsertSpecOptions) MutateInSpec {
	if opts == nil {
		opts = &ArrayInsertSpecOptions{}
	}

	return MutateInSpec{
		op:         memd.SubDocOpArrayInsert,
		createPath: opts.CreatePath,
		isXattr:    opts.IsXattr,
		path:       path,
		value:      val,
		multiValue: opts.HasMultiple,
	}
}

// ArrayAddUniqueSpecOptions are the options available to subdocument ArrayAddUnique operations.
type ArrayAddUniqueSpecOptions struct {
	CreatePath bool
	IsXattr    bool
}

// ArrayAddUniqueSpec adds an dictionary add unique operation to this mutation operation set.
func ArrayAddUniqueSpec(path string, val interface{}, opts *ArrayAddUniqueSpecOptions) MutateInSpec {
	if opts == nil {
		opts = &ArrayAddUniqueSpecOptions{}
	}

	return MutateInSpec{
		op:         memd.SubDocOpArrayAddUnique,
		createPath: opts.CreatePath,
		isXattr:    opts.IsXattr,
		path:       path,
		value:      val,
		multiValue: false,
	}
}

// CounterSpecOptions are the options available to subdocument Increment and Decrement operations.
type CounterSpecOptions struct {
	CreatePath bool
	IsXattr    bool
}

// IncrementSpec adds an increment operation to this mutation operation set.
func IncrementSpec(path string, delta int64, opts *CounterSpecOptions) MutateInSpec {
	if opts == nil {
		opts = &CounterSpecOptions{}
	}

	return MutateInSpec{
		op:         memd.SubDocOpCounter,
		createPath: opts.CreatePath,
		isXattr:    opts.IsXattr,
		path:       path,
		value:      delta,
		multiValue: false,
	}
}

// DecrementSpec adds a decrement operation to this mutation operation set.
func DecrementSpec(path string, delta int64, opts *CounterSpecOptions) MutateInSpec {
	if opts == nil {
		opts = &CounterSpecOptions{}
	}

	return MutateInSpec{
		op:         memd.SubDocOpCounter,
		createPath: opts.CreatePath,
		isXattr:    opts.IsXattr,
		path:       path,
		value:      -delta,
		multiValue: false,
	}
}
