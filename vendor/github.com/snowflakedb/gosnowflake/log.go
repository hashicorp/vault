// +build !sfdebug

// Wrapper for glog to replace direct use, so that glog usage remains optional.
// This file contains the no-op glog wrapper/emulator.

package gosnowflake

// glogWrapper is an empty struct to create a no-op glog wrapper
type glogWrapper struct{}

// V emulates the glog.V() call
func (glogWrapper) V(int) glogWrapper {
	return glogWrapper{}
}

// Check if the logging is enabled. Returns always False by default
func (glogWrapper) IsEnabled(int) bool {
	return false
}

// Flush emulates the glog.Flush() call
func (glogWrapper) Flush() {}

// Info emulates the glog.V(?).Info call
func (glogWrapper) Info(...interface{}) {}

// Infoln emulates the glog.V(?).Infoln call
func (glogWrapper) Infoln(...interface{}) {}

// Infof emulates the glog.V(?).Infof call
func (glogWrapper) Infof(...interface{}) {}

// InfoDepth emulates the glog.V(?).InfoDepth call
func (glogWrapper) InfoDepth(...interface{}) {}

// NOTE: Warning* and Error* methods are not emulated since they are not used.
// NOTE: Fatal* and Exit* methods are not emulated, since they also require additional calls (like os.Exit() and stack traces) to be compatible.

// glog is our glog emulator
var glog = glogWrapper{}
