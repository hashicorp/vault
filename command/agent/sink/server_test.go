package sink

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/go-uuid"
)

func TestServerHandler(t *testing.T) {
	t.Log("[INFO] Starting TestServerHandler")
	defer t.Log("[INFO] Finished TestServerHandler")
	core := TestCore(t)

	fs1, path1 := testFileServer(t, "")
	defer os.RemoveAll(path1)
	fs2, path2 := testFileServer(t, "")
	defer os.RemoveAll(path2)

	servers := []Server{fs1, fs2}
	go core.serverHandler.Run(servers)

	uuidStr, _ := uuid.GenerateUUID()
	core.serverHandler.TokenCh <- uuidStr

	// Give it a minute to finish writing
	time.Sleep(1 * time.Second)

	// Tell it to shut down and give it time to do so
	close(state.ShutdownCh)
	time.Sleep(1 * time.Second)

	if !core.serverHandler.Stopped() {
		t.Fatal("serverhandler did not stop")
	}

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
