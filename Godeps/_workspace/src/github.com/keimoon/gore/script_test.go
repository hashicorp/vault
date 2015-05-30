package gore

import (
	"io/ioutil"
	"os"
	"testing"
)

var scriptSet = `
return redis.call('SET', KEYS[1], ARGV[1])
`

var scriptGet = `
return redis.call('GET', KEYS[1])
`

var scriptError = `
return redis.call('ZRANGE', KEYS[1], 0, -1)
`

func TestScript(t *testing.T) {
	conn, err := Dial("localhost:6379")
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	s := NewScript()
	s.SetBody(scriptSet)
	rep, err := s.Execute(conn, 1, "kirisame", "marisa")
	if err != nil || !rep.IsOk() {
		t.Fatal(err, rep)
	}
	s.SetBody(scriptGet)
	rep, err = s.Execute(conn, 1, "kirisame")
	if err != nil {
		t.Fatal(err)
	}
	val, err := rep.String()
	if err != nil || val != "marisa" {
		t.Fatal(err, val)
	}
	s.SetBody(scriptError)
	rep, err = s.Execute(conn, 1, "kirisame")
	if err != nil || !rep.IsError() {
		t.Fatal(err, rep)
	}
	rep, err = NewCommand("FLUSHALL").Run(conn)
	if err != nil || !rep.IsOk() {
		t.Fatal(err, "not ok")
	}
}

func TestScriptMap(t *testing.T) {
	err := os.MkdirAll("testscripts", 0755)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll("testscripts")
	err = ioutil.WriteFile("testscripts/set.lua", []byte(scriptSet), 0644)
	if err != nil {
		t.Fatal(err)
	}
	err = ioutil.WriteFile("testscripts/get.lua", []byte(scriptGet), 0644)
	if err != nil {
		t.Fatal(err)
	}
	err = ioutil.WriteFile("testscripts/error.lua", []byte(scriptError), 0644)
	if err != nil {
		t.Fatal(err)
	}
	err = LoadScripts("testscripts", ".*\\.lua")
	if err != nil {
		t.Fatal(err)
	}
	conn, err := Dial("localhost:6379")
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	s := GetScript("set.lua")
	rep, err := s.Execute(conn, 1, "kirisame", "marisa")
	if err != nil || !rep.IsOk() {
		t.Fatal(err, rep)
	}
	s = GetScript("get.lua")
	rep, err = s.Execute(conn, 1, "kirisame")
	if err != nil {
		t.Fatal(err)
	}
	val, err := rep.String()
	if err != nil || val != "marisa" {
		t.Fatal(err, val)
	}
	s = GetScript("error.lua")
	rep, err = s.Execute(conn, 1, "kirisame")
	if err != nil || !rep.IsError() {
		t.Fatal(err, rep)
	}
	rep, err = NewCommand("FLUSHALL").Run(conn)
	if err != nil || !rep.IsOk() {
		t.Fatal(err, "not ok")
	}
}
