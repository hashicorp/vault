# atomic
    import "github.com/natefinch/atomic"
atomic is a go package for atomic file writing

By default, writing to a file in go (and generally any language) can fail
partway through... you then have a partially written file, which probably was
truncated when the write began, and bam, now you've lost data.

This go package avoids this problem, by writing first to a temp file, and then
overwriting the target file in an atomic way.  This is easy on linux, os.Rename
just is atomic.  However, on Windows, os.Rename is not atomic, and so bad things
can happen.  By wrapping the windows API moveFileEx, we can ensure that the move
is atomic, and we can be safe in knowing that either the move succeeds entirely,
or neither file will be modified.


## func ReplaceFile
``` go
func ReplaceFile(source, destination string) error
```
ReplaceFile atomically replaces the destination file or directory with the
source.  It is guaranteed to either replace the target file entirely, or not
change either file.


## func WriteFile
``` go
func WriteFile(filename string, r io.Reader) (err error)
```
WriteFile atomically writes the contents of r to the specified filepath.  If
an error occurs, the target file is guaranteed to be either fully written, or
not written at all.  WriteFile overwrites any file that exists at the
location (but only if the write fully succeeds, otherwise the existing file
is unmodified).

