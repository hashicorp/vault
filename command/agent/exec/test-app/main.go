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
	"strings"
	"syscall"
	"time"
)

var (
	port          uint
	ignoreStopSig bool
	ttl           time.Duration
	useSigusr1    bool
	stopAfter     time.Duration
	exitCode      int
)

func init() {
	flag.UintVar(&port, "port", 8000, "port to run the sample app on")
	flag.BoolVar(&ignoreStopSig, "ignore-stop-signal", false, "dont stop the server on SIGTERM")
	flag.DurationVar(&ttl, "ttl", 5*time.Second, "time to wait after getting the signal before exiting (ignored if `ignore-stop-signal` is set)")
	flag.BoolVar(&useSigusr1, "use-sigusr1", false, "use SIGUSR1 as the stop signal, instead of the default SIGTERM")
	flag.DurationVar(&stopAfter, "stop-after", 0, "stop the process after duration (overrides all other flags if set)")
	flag.IntVar(&exitCode, "exit-code", 0, "exit code to return when this script exits")
}

type Response struct {
	EnvVars   map[string]string `json:"env_vars"`
	ProcessID int               `json:"process_id"`
}

func newResponse() Response {
	respEnv := make(map[string]string, len(os.Environ()))
	for _, envVar := range os.Environ() {
		tokens := strings.Split(envVar, "=")
		respEnv[tokens[0]] = tokens[1]
	}

	return Response{
		EnvVars:   respEnv,
		ProcessID: os.Getpid(),
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	if r.URL.Query().Get("pretty") == "1" {
		encoder.SetIndent("", "  ")
	}
	if err := encoder.Encode(newResponse()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(buf.Bytes())
}

func main() {
	flag.Parse()
	logger := log.New(os.Stderr, "vault-agent-testing-sample-app: ", log.LstdFlags)

	if stopAfter > 0 {
		timer := time.AfterFunc(stopAfter, func() {
			logger.Printf("stopping the app early with exit code %d", exitCode)
			os.Exit(exitCode)
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
			signal.Notify(stopSig, syscall.SIGTERM)
		}

		<-stopSig

		if ignoreStopSig {
			logger.Println("ignoring stop signal!")
			logger.Printf("run: `kill %d` to stop\n", os.Getpid())
			return
		}

		time.Sleep(ttl)

		if err := server.Shutdown(context.Background()); err != nil {
			log.Printf("HTTP server Shutdown: %v", err)
		}

		close(idleConnsClosed)
	}()

	logger.Printf("starting server on port %d\n", port)
	if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		logger.Fatalf("HTTP server ListenAndServe: %v", err)
	}

	<-idleConnsClosed
	os.Exit(exitCode)
}
