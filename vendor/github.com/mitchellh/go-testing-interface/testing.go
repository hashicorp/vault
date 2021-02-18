package testing

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
)

// T is the interface that mimics the standard library *testing.T.
//
// In unit tests you can just pass a *testing.T struct. At runtime, outside
// of tests, you can pass in a RuntimeT struct from this package.
type T interface {
	Cleanup(func())
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Fail()
	FailNow()
	Failed() bool
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	Helper()
	Log(args ...interface{})
	Logf(format string, args ...interface{})
	Name() string
	Parallel()
	Skip(args ...interface{})
	SkipNow()
	Skipf(format string, args ...interface{})
	Skipped() bool
	TempDir() string
}

// TB is the interface common to T and B, copied from the standard library
// *testing.TB.
//
// This interface should be used as the type of the testing argument to any
// test helper function that exists in the main codebase, which may be invoked
// by tests of type *testing.T and *testing.B.
type TB interface {
	Cleanup(func())
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Fail()
	FailNow()
	Failed() bool
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	Helper()
	Log(args ...interface{})
	Logf(format string, args ...interface{})
	Name() string
	Skip(args ...interface{})
	SkipNow()
	Skipf(format string, args ...interface{})
	Skipped() bool
	TempDir() string
}

var tempDirReplacer struct {
	sync.Once
	r *strings.Replacer
}

// RuntimeT implements T and can be instantiated and run at runtime to
// mimic *testing.T behavior. Unlike *testing.T, this will simply panic
// for calls to Fatal. For calls to Error, you'll have to check the errors
// list to determine whether to exit yourself.
//
// Parallel does not do anything.
type RuntimeT struct {
	mu          sync.RWMutex // guards this group of fields
	skipped     bool
	failed      bool
	tempDirOnce sync.Once
	tempDir     string
	tempDirErr  error
	tempDirSeq  int32
	cleanup     func()    // optional function to be called at the end of the test
	cleanupName string    // Name of the cleanup function.
	cleanupPc   []uintptr // The stack trace at the point where Cleanup was called.
}

// The maximum number of stack frames to go through when skipping helper functions for
// the purpose of decorating log messages.
const maxStackLen = 50

func (t *RuntimeT) Error(args ...interface{}) {
	log.Println(fmt.Sprintln(args...))
	t.Fail()
}

func (t *RuntimeT) Errorf(format string, args ...interface{}) {
	log.Printf(format, args...)
	t.Fail()
}

func (t *RuntimeT) Fail() {
	t.failed = true
}

func (t *RuntimeT) FailNow() {
	panic("testing.T failed, see logs for output (if any)")
}

func (t *RuntimeT) Failed() bool {
	return t.failed
}

func (t *RuntimeT) Fatal(args ...interface{}) {
	log.Print(args...)
	t.FailNow()
}

func (t *RuntimeT) Fatalf(format string, args ...interface{}) {
	log.Printf(format, args...)
	t.FailNow()
}

func (t *RuntimeT) Log(args ...interface{}) {
	log.Println(fmt.Sprintln(args...))
}

func (t *RuntimeT) Logf(format string, args ...interface{}) {
	log.Println(fmt.Sprintf(format, args...))
}

func (t *RuntimeT) Name() string {
	return ""
}

func (t *RuntimeT) Parallel() {}

func (t *RuntimeT) Skip(args ...interface{}) {
	log.Print(args...)
	t.SkipNow()
}

func (t *RuntimeT) SkipNow() {
	t.skipped = true
}

func (t *RuntimeT) Skipf(format string, args ...interface{}) {
	log.Printf(format, args...)
	t.SkipNow()
}

func (t *RuntimeT) Skipped() bool {
	return t.skipped
}

// TempDir returns a temporary directory for the test to use.
// The directory is automatically removed by Cleanup when the test and
// all its subtests complete.
// Each subsequent call to t.TempDir returns a unique directory;
// if the directory creation fails, TempDir terminates the test by calling Fatal.
//
// This logic is copied from the standard go library
func (t *RuntimeT) TempDir() string {
	// Use a single parent directory for all the temporary directories
	// created by a test, each numbered sequentially.
	t.tempDirOnce.Do(func() {
		t.Helper()

		// ioutil.TempDir doesn't like path separators in its pattern,
		// so mangle the name to accommodate subtests.
		tempDirReplacer.Do(func() {
			tempDirReplacer.r = strings.NewReplacer("/", "_", "\\", "_", ":", "_")
		})
		pattern := tempDirReplacer.r.Replace(t.Name())

		t.tempDir, t.tempDirErr = ioutil.TempDir("", pattern)
		if t.tempDirErr == nil {
			t.Cleanup(func() {
				if err := os.RemoveAll(t.tempDir); err != nil {
					t.Errorf("TempDir RemoveAll cleanup: %v", err)
				}
			})
		}
	})
	if t.tempDirErr != nil {
		t.Fatalf("TempDir: %v", t.tempDirErr)
	}
	seq := atomic.AddInt32(&t.tempDirSeq, 1)
	dir := fmt.Sprintf("%s%c%03d", t.tempDir, os.PathSeparator, seq)
	if err := os.Mkdir(dir, 0777); err != nil {
		t.Fatalf("TempDir: %v", err)
	}
	return dir
}

func (t *RuntimeT) Helper() {}

// Cleanup registers a function to be called when the test and all its
// subtests complete. Cleanup functions will be called in last added,
// first called order.
//
// This logic is copied from the standard go library
func (t *RuntimeT) Cleanup(f func()) {
	t.mu.Lock()
	defer t.mu.Unlock()
	oldCleanup := t.cleanup
	oldCleanupPc := t.cleanupPc
	t.cleanup = func() {
		if oldCleanup != nil {
			defer func() {
				t.mu.Lock()
				t.cleanupPc = oldCleanupPc
				t.mu.Unlock()
				oldCleanup()
			}()
		}
		t.mu.Lock()
		t.cleanupName = callerName(0)
		t.mu.Unlock()
		f()
	}
	var pc [maxStackLen]uintptr
	// Skip two extra frames to account for this function and runtime.Callers itself.
	n := runtime.Callers(2, pc[:])
	t.cleanupPc = pc[:n]
}

// callerName gives the function name (qualified with a package path)
// for the caller after skip frames (where 0 means the current function).
//
// This logic is copied from the standard go library
func callerName(skip int) string {
	// Make room for the skip PC.
	var pc [1]uintptr
	n := runtime.Callers(skip+2, pc[:]) // skip + runtime.Callers + callerName
	if n == 0 {
		panic("testing: zero callers found")
	}
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	return frame.Function
}
