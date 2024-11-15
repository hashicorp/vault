// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package slug

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/go-slug/internal/ignorefiles"
	"github.com/hashicorp/go-slug/internal/unpackinfo"
)

// Meta provides detailed information about a slug.
type Meta struct {
	// The list of files contained in the slug.
	Files []string

	// Total size of the slug in bytes.
	Size int64
}

// IllegalSlugError indicates the provided slug (io.Writer for Pack, io.Reader
// for Unpack) violates a rule about its contents. For example, an absolute or
// external symlink. It implements the error interface.
type IllegalSlugError struct {
	Err error
}

func (e *IllegalSlugError) Error() string {
	return fmt.Sprintf("illegal slug error: %v", e.Err)
}

// Unwrap returns the underlying issue with the provided Slug into the error
// chain.
func (e *IllegalSlugError) Unwrap() error { return e.Err }

// externalSymlink is a simple abstraction for a information about a symlink target
type externalSymlink struct {
	absTarget string
	target    string
	info      os.FileInfo
}

// PackerOption is a functional option that can configure non-default Packers.
type PackerOption func(*Packer) error

// ApplyTerraformIgnore is a PackerOption that will apply the .terraformignore
// rules and skip packing files it specifies.
func ApplyTerraformIgnore() PackerOption {
	return func(p *Packer) error {
		p.applyTerraformIgnore = true
		return nil
	}
}

// DereferenceSymlinks is a PackerOption that will allow symlinks that
// reference a target outside of the source directory by copying the link
// target, turning it into a normal file within the archive.
func DereferenceSymlinks() PackerOption {
	return func(p *Packer) error {
		p.dereference = true
		return nil
	}
}

// AllowSymlinkTarget relaxes safety checks on symlinks with targets matching
// path. Specifically, absolute symlink targets (e.g. "/foo/bar") and relative
// targets (e.g. "../foo/bar") which resolve to a path outside of the
// source/destination directories for pack/unpack operations respectively, may
// be expressly permitted, whereas they are forbidden by default. Exercise
// caution when using this option. A symlink matches path if its target
// resolves to path exactly, or if path is a parent directory of target.
func AllowSymlinkTarget(path string) PackerOption {
	return func(p *Packer) error {
		p.allowSymlinkTargets = append(p.allowSymlinkTargets, path)
		return nil
	}
}

// Packer holds options for the Pack function.
type Packer struct {
	dereference          bool
	applyTerraformIgnore bool
	allowSymlinkTargets  []string
}

// NewPacker is a constructor for Packer.
func NewPacker(options ...PackerOption) (*Packer, error) {
	p := &Packer{
		dereference:          false,
		applyTerraformIgnore: false,
	}

	for _, opt := range options {
		if err := opt(p); err != nil {
			return nil, fmt.Errorf("option failed: %w", err)
		}
	}

	return p, nil
}

// Pack at the package level is used to maintain compatibility with existing
// code that relies on this function signature. New options related to packing
// slugs should be added to the Packer struct instead.
func Pack(src string, w io.Writer, dereference bool) (*Meta, error) {
	p := Packer{
		dereference: dereference,

		// This defaults to false in NewPacker, but is true here. This matches
		// the old behavior of Pack, which always used .terraformignore.
		applyTerraformIgnore: true,
	}
	return p.Pack(src, w)
}

// Pack creates a slug from a src directory, and writes the new slug
// to w. Returns metadata about the slug and any errors.
//
// When dereference is set to true, symlinks with a target outside of
// the src directory will be dereferenced. When dereference is set to
// false symlinks with a target outside the src directory are omitted
// from the slug.
func (p *Packer) Pack(src string, w io.Writer) (*Meta, error) {
	// Gzip compress all the output data.
	gzipW, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
	if err != nil {
		// This error is only raised when an incorrect gzip level is
		// specified.
		return nil, err
	}

	// Tar the file contents.
	tarW := tar.NewWriter(gzipW)

	// Track the metadata details as we go.
	meta := &Meta{}

	info, err := os.Lstat(src)
	if err != nil {
		return nil, err
	}

	// Check if the root (src) is a symlink
	if info.Mode()&os.ModeSymlink != 0 {
		src, err = os.Readlink(src)
		if err != nil {
			return nil, err
		}
	}

	// Load the ignore rule configuration, which will use
	// defaults if no .terraformignore is configured
	var ignoreRules *ignorefiles.Ruleset
	if p.applyTerraformIgnore {
		ignoreRules = parseIgnoreFile(src)
	}

	// Ensure the source path provided is absolute
	src, err = filepath.Abs(src)
	if err != nil {
		return nil, fmt.Errorf("failed to read absolute path for source: %w", err)
	}

	// Walk the tree of files.
	err = filepath.Walk(src, p.packWalkFn(src, src, src, tarW, meta, ignoreRules))
	if err != nil {
		return nil, err
	}

	// Flush the tar writer.
	if err := tarW.Close(); err != nil {
		return nil, fmt.Errorf("failed to close the tar archive: %w", err)
	}

	// Flush the gzip writer.
	if err := gzipW.Close(); err != nil {
		return nil, fmt.Errorf("failed to close the gzip writer: %w", err)
	}

	return meta, nil
}

func (p *Packer) packWalkFn(root, src, dst string, tarW *tar.Writer, meta *Meta, ignoreRules *ignorefiles.Ruleset) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Get the relative path from the current src directory.
		subpath, err := filepath.Rel(src, path)
		if err != nil {
			return fmt.Errorf("failed to get relative path for file %q: %w", path, err)
		}
		if subpath == "." {
			return nil
		}

		if r := matchIgnoreRules(subpath, ignoreRules); r.Excluded {
			return nil
		}

		// Catch directories so we don't end up with empty directories,
		// the files are ignored correctly
		if info.IsDir() {
			if r := matchIgnoreRules(subpath+string(os.PathSeparator), ignoreRules); r.Excluded {
				if r.Dominating {
					return filepath.SkipDir
				} else {
					return nil
				}
			}
		}

		// Get the relative path from the initial root directory.
		subpath, err = filepath.Rel(root, strings.Replace(path, src, dst, 1))
		if err != nil {
			return fmt.Errorf("failed to get relative path for file %q: %w", path, err)
		}
		if subpath == "." {
			return nil
		}

		// Check the file type and if we need to write the body.
		keepFile, writeBody := checkFileMode(info.Mode())
		if !keepFile {
			return nil
		}

		fm := info.Mode()
		// An "Unknown" format is imposed because this is the default but also because
		// it imposes the simplest behavior. Notably, the mod time is preserved by rounding
		// to the nearest second. During unpacking, these rounded timestamps are restored
		// upon the corresponding file/directory/symlink.
		header := &tar.Header{
			Format:  tar.FormatUnknown,
			Name:    filepath.ToSlash(subpath),
			ModTime: info.ModTime(),
			Mode:    int64(fm.Perm()),
		}

		switch {
		case info.IsDir():
			header.Typeflag = tar.TypeDir
			header.Name += "/"

		case fm.IsRegular():
			header.Typeflag = tar.TypeReg
			header.Size = info.Size()

		case fm&os.ModeSymlink != 0:
			// Read the symlink file to find the destination.
			target, err := os.Readlink(path)
			if err != nil {
				return fmt.Errorf("failed to read symlink %q: %w", path, err)
			}

			// Check if the symlink's target falls within the root.
			if ok, err := p.validSymlink(root, path, target); ok {
				// We can simply copy the link.
				header.Typeflag = tar.TypeSymlink
				header.Linkname = filepath.ToSlash(target)
				break
			} else if !p.dereference {
				// If the target does not fall within the root and dereference
				// is set to false, we can't resolve the target and copy its
				// contents.
				return err
			}

			// Attempt to follow the external target so we can copy its contents
			resolved, err := p.resolveExternalLink(root, path)
			if err != nil {
				return err
			}

			// If the target is a directory we can recurse into the target
			// directory by calling the packWalkFn with updated arguments.
			if resolved.info.IsDir() {
				return filepath.Walk(resolved.absTarget, p.packWalkFn(root, resolved.absTarget, path, tarW, meta, ignoreRules))
			}

			// Dereference this symlink by updating the header with the target file
			// details and set writeBody to true so the body will be written.
			header.Typeflag = tar.TypeReg
			header.ModTime = resolved.info.ModTime()
			header.Mode = int64(resolved.info.Mode().Perm())
			header.Size = resolved.info.Size()
			writeBody = true

		default:
			return fmt.Errorf("unexpected file mode %v", fm)
		}

		// Write the header first to the archive.
		if err := tarW.WriteHeader(header); err != nil {
			return fmt.Errorf("failed writing archive header for file %q: %w", path, err)
		}

		// Account for the file in the list.
		meta.Files = append(meta.Files, header.Name)

		// Skip writing file data for certain file types (above).
		if !writeBody {
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("failed opening file %q for archiving: %w", path, err)
		}
		defer f.Close()

		size, err := io.Copy(tarW, f)
		if err != nil {
			return fmt.Errorf("failed copying file %q to archive: %w", path, err)
		}

		// Add the size we copied to the body.
		meta.Size += size

		return nil
	}
}

// resolveExternalSymlink attempts to recursively follow target paths if we
// encounter a symbolic link chain. It returns path information about the final
// target pointing to a regular file or directory.
func (p *Packer) resolveExternalLink(root string, path string) (*externalSymlink, error) {
	// Read the symlink file to find the destination.
	target, err := os.Readlink(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read symlink %q: %w", path, err)
	}

	// Get the absolute path of the symlink target.
	absTarget := target
	if !filepath.IsAbs(absTarget) {
		absTarget = filepath.Join(filepath.Dir(path), target)
	}
	if !filepath.IsAbs(absTarget) {
		absTarget = filepath.Join(root, absTarget)
	}

	// Get the file info for the target.
	info, err := os.Lstat(absTarget)
	if err != nil {
		return nil, fmt.Errorf("failed to get file info from file %q: %w", target, err)
	}

	// Recurse if the symlink resolves to another symlink
	if info.Mode()&os.ModeSymlink != 0 {
		return p.resolveExternalLink(root, absTarget)
	}

	return &externalSymlink{
		absTarget: absTarget,
		target:    target,
		info:      info,
	}, err
}

// Unpack is used to read and extract the contents of a slug to the dst
// directory, which must be an absolute path. Symlinks within the slug
// are supported, provided their targets are relative and point to paths
// within the destination directory.
func Unpack(r io.Reader, dst string) error {
	p := &Packer{}
	return p.Unpack(r, dst)
}

// Unpack unpacks the archive data in r into directory dst.
func (p *Packer) Unpack(r io.Reader, dst string) error {
	// Track directory times and permissions so they can be restored after all files
	// are extracted. This metadata modification is delayed because extracting files
	// into a new directory would necessarily change its timestamps. By way of
	// comparison, see
	// https://www.gnu.org/software/tar/manual/html_node/Directory-Modification-Times-and-Permissions.html
	// for more details about how tar attempts to preserve file metadata.
	directoriesExtracted := []unpackinfo.UnpackInfo{}

	// Decompress as we read.
	uncompressed, err := gzip.NewReader(r)
	if err != nil {
		return fmt.Errorf("failed to decompress slug: %w", err)
	}

	// Untar as we read.
	untar := tar.NewReader(uncompressed)

	// Unpackage all the contents into the directory.
	for {
		header, err := untar.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to untar slug: %w", err)
		}

		// If the entry has no name, ignore it.
		if header.Name == "" {
			continue
		}

		info, err := unpackinfo.NewUnpackInfo(dst, header)
		if err != nil {
			return &IllegalSlugError{Err: err}
		}

		// Make the directories to the path.
		dir := filepath.Dir(info.Path)

		// Timestamps and permissions will be restored after all files are extracted.
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %q: %w", dir, err)
		}

		// Handle symlinks, directories, non-regular files
		if info.IsSymlink() {
			if ok, err := p.validSymlink(dst, header.Name, header.Linkname); ok {
				// Create the symlink.
				if err = os.Symlink(header.Linkname, info.Path); err != nil {
					return fmt.Errorf("failed creating symlink (%q -> %q): %w",
						header.Name, header.Linkname, err)
				}
			} else {
				return err
			}

			if err := info.RestoreInfo(); err != nil {
				return err
			}

			continue
		}

		if info.IsDirectory() {
			// Restore directory info after all files are extracted because
			// the extraction process changes directory's timestamps.
			directoriesExtracted = append(directoriesExtracted, info)
			continue
		}

		// The remaining logic only applies to regular files
		if !info.IsRegular() {
			continue
		}

		// Open a handle to the destination.
		fh, err := os.Create(info.Path)
		if err != nil {
			// This mimics tar's behavior wrt the tar file containing duplicate files
			// and it allowing later ones to clobber earlier ones even if the file
			// has perms that don't allow overwriting. The file permissions will be restored
			// once the file contents are copied.
			if os.IsPermission(err) {
				os.Chmod(info.Path, 0600)
				fh, err = os.Create(info.Path)
			}

			if err != nil {
				return fmt.Errorf("failed creating file %q: %w", info.Path, err)
			}
		}

		// Copy the contents of the file.
		_, err = io.Copy(fh, untar)
		fh.Close()
		if err != nil {
			return fmt.Errorf("failed to copy slug file %q: %w", info.Path, err)
		}

		if err := info.RestoreInfo(); err != nil {
			return err
		}
	}

	for _, dir := range directoriesExtracted {
		if err := dir.RestoreInfo(); err != nil {
			return err
		}
	}

	return nil
}

// Given a "root" directory, the path to a symlink within said root, and the
// target of said symlink, validSymlink checks that the target either falls
// into root somewhere, or is explicitly allowed per the Packer's config.
func (p *Packer) validSymlink(root, path, target string) (bool, error) {
	// Get the absolute path to root.
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return false, fmt.Errorf("failed making path %q absolute: %w", root, err)
	}

	// Get the absolute path to the file path.
	absPath := path
	if !filepath.IsAbs(absPath) {
		absPath = filepath.Join(absRoot, path)
	}

	// Get the absolute path of the symlink target.
	var absTarget string
	if filepath.IsAbs(target) {
		absTarget = filepath.Clean(target)
	} else {
		absTarget = filepath.Join(filepath.Dir(absPath), target)
	}

	// Target falls within root.
	if strings.HasPrefix(absTarget, absRoot) {
		return true, nil
	}

	// The link target is outside of root. Check if it is allowed.
	for _, prefix := range p.allowSymlinkTargets {
		// Ensure prefix is absolute.
		if !filepath.IsAbs(prefix) {
			prefix = filepath.Join(absRoot, prefix)
		}

		// Exact match is allowed.
		if absTarget == prefix {
			return true, nil
		}

		// Prefix match of a directory is allowed.
		if !strings.HasSuffix(prefix, "/") {
			prefix += "/"
		}
		if strings.HasPrefix(absTarget, prefix) {
			return true, nil
		}
	}

	return false, &IllegalSlugError{
		Err: fmt.Errorf(
			"invalid symlink (%q -> %q) has external target",
			path, target,
		),
	}
}

// checkFileMode is used to examine an os.FileMode and determine if it should
// be included in the archive, and if it has a data body which needs writing.
func checkFileMode(m os.FileMode) (keep, body bool) {
	switch {
	case m.IsDir():
		return true, false

	case m.IsRegular():
		return true, true

	case m&os.ModeSymlink != 0:
		return true, false
	}

	return false, false
}
