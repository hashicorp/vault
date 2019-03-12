package listenerutil

import (
	"io/ioutil"
	"os"
	"strconv"
	"testing"
)

func TestUnixSocketListener(t *testing.T) {
	socket, err := ioutil.TempFile("", "socket")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(socket.Name())

	uid, gid := os.Getuid(), os.Getgid()
	l, err := UnixSocketListener(socket.Name(), &UnixSocketsConfig{
		User:  strconv.Itoa(uid),
		Group: strconv.Itoa(gid),
		Mode:  "644",
	})
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()

	fi, err := os.Stat(socket.Name())
	if err != nil {
		t.Fatal(err)
	}

	mode, err := strconv.ParseUint("644", 8, 32)
	if err != nil {
		t.Fatal(err)
	}
	if fi.Mode().Perm() != os.FileMode(mode) {
		t.Fatalf("failed to set permissions on the socket file")
	}
}
