package gore

import (
	"crypto/sha1"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strings"
	"sync"
)

// Script represents a Lua script.
type Script struct {
	body string
	sha  string
	lock sync.RWMutex
}

// NewScript returns a new Lua script
func NewScript() *Script {
	return &Script{}
}

// SetBody sets script body and its SHA value
func (s *Script) SetBody(body string) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.body = strings.TrimSpace(body)
	return s.createSHA()
}

// ReadFromFile reads the script from a file
func (s *Script) ReadFromFile(file string) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	s.body = strings.TrimSpace(string(b))
	return s.createSHA()
}

// Execute runs the script over a connection
func (s *Script) Execute(conn *Conn, keyCount int, keysAndArgs ...interface{}) (*Reply, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	if s.body == "" {
		return nil, ErrEmptyScript
	}
	args := make([]interface{}, len(keysAndArgs)+2)
	args[0] = s.sha
	args[1] = keyCount
	for i := range keysAndArgs {
		args[i+2] = keysAndArgs[i]
	}
	rep, err := NewCommand("EVALSHA", args...).Run(conn)
	if err != nil {
		return nil, err
	}
	if !rep.IsError() {
		return rep, nil
	}
	errorMessage, _ := rep.Error()
	if !strings.HasPrefix(errorMessage, "NOSCRIPT") {
		return rep, nil
	}
	args[0] = s.body
	return NewCommand("EVAL", args...).Run(conn)
}

func (s *Script) createSHA() error {
	h := sha1.New()
	_, err := io.WriteString(h, s.body)
	if err != nil {
		return err
	}
	s.sha = fmt.Sprintf("%x", h.Sum(nil))
	return nil
}

// ScriptMap is a thread-safe map from script name to its content
type ScriptMap struct {
	scripts map[string]*Script
	lock    sync.RWMutex
}

// NewScriptMap makes a new ScriptMap
func NewScriptMap() *ScriptMap {
	return &ScriptMap{
		scripts: make(map[string]*Script),
	}
}

// Load loads all script files from a folder with a regular expression pattern.
// Loaded script will be keyed by its file name. This method can be called many times
// to reload script files.
func (sm *ScriptMap) Load(folder, pattern string) error {
	sm.lock.Lock()
	sm.lock.Unlock()
	r, err := regexp.Compile(pattern)
	if err != nil {
		return err
	}
	dir, err := os.Open(folder)
	if err != nil {
		return err
	}
	defer dir.Close()
	infos, err := dir.Readdir(0)
	if err != nil {
		return err
	}
	for _, fi := range infos {
		if fi.IsDir() || !r.MatchString(fi.Name()) {
			continue
		}
		script := NewScript()
		err := script.ReadFromFile(path.Join(folder, fi.Name()))
		if err == nil {
			sm.scripts[fi.Name()] = script
		}
	}
	return nil
}

// Get a script by its name. Nil value will be returned if the name
// is not found
func (sm *ScriptMap) Get(scriptName string) *Script {
	sm.lock.RLock()
	defer sm.lock.RUnlock()
	return sm.scripts[scriptName]
}

// Add a script to the script map
func (sm *ScriptMap) Add(scriptName string, script *Script) {
	sm.lock.Lock()
	defer sm.lock.Unlock()
	sm.scripts[scriptName] = script
}

// Delete a script from the script map
func (sm *ScriptMap) Delete(scriptName string) {
	sm.lock.Lock()
	defer sm.lock.Unlock()
	delete(sm.scripts, scriptName)
}

var defaultScriptMap = NewScriptMap()

// LoadScripts loads all script files from a folder with a regular expression pattern to the default
// script map.
// Loaded script will be keyed by its file name. This method can be called many times
// to reload script files.
func LoadScripts(folder, pattern string) error {
	return defaultScriptMap.Load(folder, pattern)
}

// GetScript a script by its name from defaultScriptMap . Nil value will be returned if the name
// is not found
func GetScript(scriptName string) *Script {
	return defaultScriptMap.Get(scriptName)
}

// AddScript a script to the default script map
func AddScript(scriptName string, script *Script) {
	defaultScriptMap.Add(scriptName, script)
}

// DeleteScript a script from the default script map
func DeleteScript(scriptName string) {
	defaultScriptMap.Delete(scriptName)
}
