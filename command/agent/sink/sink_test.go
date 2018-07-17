package sink

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	hclog "github.com/hashicorp/go-hclog"
	uuid "github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/helper/logging"
)

func TestSinkServer(t *testing.T) {
	log := logging.NewVaultLogger(hclog.Trace)

	fs1, path1 := testFileSink(t, log)
	defer os.RemoveAll(path1)
	fs2, path2 := testFileSink(t, log)
	defer os.RemoveAll(path2)

	ss := NewSinkServer(&SinkConfig{
		Logger: log.Named("sink.server"),
	})

	uuidStr, _ := uuid.GenerateUUID()
	in := make(chan string)
	sinks := []Sink{fs1, fs2}
	go ss.Run(in, sinks)

	// Seed a token
	in <- uuidStr

	// Give it a minute to finish writing
	time.Sleep(1 * time.Second)

	// Tell it to shut down and give it time to do so
	close(ss.ShutdownCh)
	<-ss.DoneCh

	for _, path := range []string{path1, path2} {
		fileBytes, err := ioutil.ReadFile(fmt.Sprintf("%s/token", path))
		if err != nil {
			t.Fatal(err)
		}

		if string(fileBytes) != uuidStr {
			t.Fatalf("expected %s, got %s", uuidStr, string(fileBytes))
		}
	}
}
