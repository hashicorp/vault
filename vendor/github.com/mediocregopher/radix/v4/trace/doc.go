// Package trace contains all the types provided for tracing within the radix
// package. With tracing a user is able to pull out fine-grained runtime events
// as they happen, which is useful for gathering metrics, logging, performance
// analysis, etc...
//
// Events which are eligible for tracing are those which would not be possible
// to access otherwise.
package trace
