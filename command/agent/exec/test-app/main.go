package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	port          uint
	ignoreStopSig bool
	ttl           time.Duration
	useSigusr1    bool
	stopAfter     time.Duration
)

func init() {
	flag.UintVar(&port, "port", 8000, "port to run the sample app on")
	flag.BoolVar(&ignoreStopSig, "ignore-stop-signal", false, "dont stop the server on SIGINT")
	flag.DurationVar(&ttl, "ttl", 5*time.Second, "time to wait after getting the signal before exiting (ignored if `ignore-stop-signal` is set)")
	flag.BoolVar(&useSigusr1, "use-sigusr1", false, "use SIGUSR1 as the stop signal, instead of the default SIGINT")
	flag.DurationVar(&stopAfter, "stop-after", 0, "stop the process after duration (overrides all other flags if set)")
}

type Response struct {
	EnvVars   []string `json:"env_vars"`
	ProcessID int      `json:"process_id"`
}

func newResponse() Response {
	return Response{
		EnvVars:   os.Environ(),
		ProcessID: os.Getpid(),
	}
}

func handler(w http.ResponseWriter, _ *http.Request) {
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	if err := encoder.Encode(newResponse()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(buf.Bytes())
}

func main() {
	flag.Parse()

	if stopAfter > 0 {
		timer := time.AfterFunc(stopAfter, func() {
			os.Exit(0)
		})
		defer timer.Stop()
	}

	server := http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: http.HandlerFunc(handler),
	}
	idleConnsClosed := make(chan struct{})

	go func() {
		stopSig := make(chan os.Signal, 1)
		if useSigusr1 {
			signal.Notify(stopSig, syscall.SIGUSR1)
		} else {
			signal.Notify(stopSig, syscall.SIGINT)
		}

		<-stopSig

		if ignoreStopSig {
			fmt.Fprintln(os.Stderr, "ignoring stop signal!")
			fmt.Fprintf(os.Stderr, "run: `kill %d` to stop\n", os.Getpid())
			return
		}

		time.Sleep(ttl)

		if err := server.Shutdown(context.Background()); err != nil {
			log.Printf("HTTP server Shutdown: %v", err)
		}

		close(idleConnsClosed)
	}()

	if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}

	<-idleConnsClosed
}
