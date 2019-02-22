package astcache

import (
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"time"

	"golang.org/x/tools/go/packages"

	"github.com/golangci/golangci-lint/pkg/logutils"
)

type File struct {
	F    *ast.File
	Fset *token.FileSet
	Name string
	Err  error
}

type Cache struct {
	m   map[string]*File // map from absolute file path to file data
	s   []*File
	log logutils.Log
}

func NewCache(log logutils.Log) *Cache {
	return &Cache{
		m:   map[string]*File{},
		log: log,
	}
}

func (c Cache) ParsedFilenames() []string {
	var keys []string
	for k := range c.m {
		keys = append(keys, k)
	}
	return keys
}

func (c Cache) normalizeFilename(filename string) string {
	if filepath.IsAbs(filename) {
		return filepath.Clean(filename)
	}

	absFilename, err := filepath.Abs(filename)
	if err != nil {
		c.log.Warnf("Can't abs-ify filename %s: %s", filename, err)
		return filename
	}

	return absFilename
}

func (c Cache) Get(filename string) *File {
	return c.m[c.normalizeFilename(filename)]
}

func (c Cache) GetAllValidFiles() []*File {
	return c.s
}

func (c *Cache) prepareValidFiles() {
	files := make([]*File, 0, len(c.m))
	for _, f := range c.m {
		if f.Err != nil || f.F == nil {
			continue
		}
		files = append(files, f)
	}
	c.s = files
}

func LoadFromFilenames(log logutils.Log, filenames ...string) *Cache {
	c := NewCache(log)

	fset := token.NewFileSet()
	for _, filename := range filenames {
		c.parseFile(filename, fset)
	}

	c.prepareValidFiles()
	return c
}

func LoadFromPackages(pkgs []*packages.Package, log logutils.Log) (*Cache, error) {
	c := NewCache(log)

	for _, pkg := range pkgs {
		c.loadFromPackage(pkg)
	}

	c.prepareValidFiles()
	return c, nil
}

func (c *Cache) loadFromPackage(pkg *packages.Package) {
	if len(pkg.Syntax) == 0 || len(pkg.GoFiles) != len(pkg.CompiledGoFiles) {
		// len(pkg.Syntax) == 0 if only filenames are loaded
		// lengths aren't equal if there are preprocessed files (cgo)
		startedAt := time.Now()

		// can't use pkg.Fset: it will overwrite offsets by preprocessed files
		fset := token.NewFileSet()
		for _, f := range pkg.GoFiles {
			c.parseFile(f, fset)
		}

		c.log.Infof("Parsed AST of all pkg.GoFiles: %s for %s", pkg.GoFiles, time.Since(startedAt))
		return
	}

	for _, f := range pkg.Syntax {
		pos := pkg.Fset.Position(f.Pos())
		if pos.Filename == "" {
			continue
		}

		c.m[pos.Filename] = &File{
			F:    f,
			Fset: pkg.Fset,
			Name: pos.Filename,
		}
	}
}

func (c *Cache) parseFile(filePath string, fset *token.FileSet) {
	if fset == nil {
		fset = token.NewFileSet()
	}

	filePath = c.normalizeFilename(filePath)

	// comments needed by e.g. golint
	f, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	c.m[filePath] = &File{
		F:    f,
		Fset: fset,
		Err:  err,
		Name: filePath,
	}
	if err != nil {
		c.log.Warnf("Can't parse AST of %s: %s", filePath, err)
	}
}
