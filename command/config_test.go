package command

import (
	"path/filepath"
	"reflect"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	config, err := LoadConfig(filepath.Join(FixturePath, "config.hcl"))
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	expected := &Config{
		TokenHelper: "foo",
	}
	if !reflect.DeepEqual(expected, config) {
		t.Fatalf("bad: %#v", config)
	}
}
