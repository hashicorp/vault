// +build memprofiler

package command

import (
	"os"
	"path/filepath"
	"runtime"
	"runtime/pprof"
	"time"
)

func init() {
	memProfilerEnabled = true
}

func (c *ServerCommand) startMemProfiler() {
	profileDir := filepath.Join(os.TempDir(), "vaultprof")
	if err := os.MkdirAll(profileDir, 0700); err != nil {
		c.logger.Debug("could not create profile directory", "error", err)
		return
	}

	go func() {
		for {
			filename := filepath.Join(profileDir, time.Now().UTC().Format("20060102_150405")) + ".pprof"
			f, err := os.Create(filename)
			if err != nil {
				c.logger.Debug("could not create memory profile", "error", err)
			}
			runtime.GC()
			if err := pprof.WriteHeapProfile(f); err != nil {
				c.logger.Debug("could not write memory profile", "error", err)
			}
			f.Close()
			c.logger.Debug("wrote memory profile", "filename", filename)
			time.Sleep(5 * time.Minute)
		}
	}()
}
