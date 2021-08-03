package main // import "github.com/hashicorp/vault"

import (
	"log"
	"os"
	"runtime/pprof"
	"time"

	"github.com/hashicorp/vault/command"
)

func main() {
	f, err := os.Create("/tmp/vault.pprof")
	if err != nil {
		log.Fatal(err)
	}
	pprof.StartCPUProfile(f)
	go func() {
		time.Sleep(30 * time.Second)
		pprof.StopCPUProfile()
	}()
	os.Exit(command.Run(os.Args[1:]))
}
