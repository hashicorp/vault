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
	logger := log.New(os.Stderr, "test-app: ", log.LstdFlags)

	if err := run(logger); err != nil {
		log.Fatalf("error: %v\n", err)
	}

	logger.Printf("exit code: %d\n", exitCode)

	os.Exit(exitCode)
}

func run(logger *log.Logger) error {
	/* */ logger.Println("run: started")
	defer logger.Println("run: done")

	ctx, cancelContextFunc := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancelContextFunc()

	flag.Parse()

	server := http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: http.HandlerFunc(handler),
	}

	doneCh := make(chan struct{})

	go func() {
		stopSig := make(chan os.Signal, 1)
		if useSigusr1 {
			signal.Notify(stopSig, syscall.SIGUSR1)
		} else {
			signal.Notify(stopSig, syscall.SIGTERM)
		}

		select {
		case <-ctx.Done():
			logger.Println("context done: exiting")

		case s := <-stopSig:
			logger.Printf("signal %q: received\n", s)

			if ignoreStopSig {
				logger.Printf("signal %q: ignored; run `kill -9 %d` to stop\n", s, os.Getpid())
				return
			}

			logger.Printf("signal %q: sleeping for %v simulate cleanup\n", s, ttl)

			time.Sleep(ttl)

		case <-time.After(stopAfter):
			logger.Printf("stopping after: %v\n", stopAfter)
		}

		if err := server.Shutdown(context.Background()); err != nil {
			log.Printf("server shutdown error: %v", err)
		}

		close(doneCh)
	}()

	logger.Printf("server %s: started\n", server.Addr)

	if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("could not start the server: %v", err)
	}

	logger.Printf("server %s: done\n", server.Addr)

	<-doneCh

	return nil
}
