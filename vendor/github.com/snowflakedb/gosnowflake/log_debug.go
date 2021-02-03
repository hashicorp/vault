// +build sfdebug

// Wrapper for glog to replace direct use, so that glog usage remains optional.
// This file contains the actual/operational glog wrapper.

package gosnowflake

import logger "github.com/snowflakedb/glog"

// glogWrapper wraps glog's Verbose type, enabling the use of glog.V().* methods directly
type glogWrapper struct {
	logger.Verbose
}

// V provides a wrapper for the glog.V() call
func (l *glogWrapper) V(level int32) glogWrapper {
	return glogWrapper{logger.V(logger.Level(level))}
}

func (l *glogWrapper) IsEnabled(level int32) bool {
	return bool(logger.V(logger.Level(level)))
}

// Flush calls flush on the underlying logger
func (l *glogWrapper) Flush() {
	logger.Flush()
}

// glog is our glog wrapper
var glog = glogWrapper{}
