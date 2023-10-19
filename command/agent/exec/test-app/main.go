// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package main

// This is a test application that is used by TestExecServer_Run to verify
// the behavior of vault agent running as a process supervisor.
//
// The app will automatically exit after 1 minute or the --stop-after interval,
// whichever comes first. It also can serve its loaded environment variables on
// the given --port. This app will also return the given --exit-code and
// terminate on SIGTERM unless --use-sigusr1 is specified.

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

var (
	port                 uint
	sleepAfterStopSignal time.Duration
	useSigusr1StopSignal bool
	stopAfter            time.Duration
	exitCode             int
	logToStdout          bool
)

func init() {
	flag.UintVar(&port, "port", 34000, "port to run the test app on")
	flag.DurationVar(&sleepAfterStopSignal, "sleep-after-stop-signal", 1*time.Second, "time to sleep after getting the signal before exiting")
	flag.BoolVar(&useSigusr1StopSignal, "use-sigusr1", false, "use SIGUSR1 as the stop signal, instead of the default SIGTERM")
	flag.DurationVar(&stopAfter, "stop-after", 60*time.Second, "stop the process after duration (overrides all other flags if set)")
	flag.IntVar(&exitCode, "exit-code", 0, "exit code to return when this script exits")
	flag.BoolVar(&logToStdout, "log-to-stdout", false, "send logs to stdout instead of stderr")
}

type Response struct {
	EnvironmentVariables map[string]string `json:"environment_variables"`
	ProcessID            int               `json:"process_id"`
}

func newResponse() Response {
	respEnv := make(map[string]string, len(os.Environ()))
	for _, envVar := range os.Environ() {
		tokens := strings.Split(envVar, "=")
		respEnv[tokens[0]] = tokens[1]
	}

	return Response{
		EnvironmentVariables: respEnv,
		ProcessID:            os.Getpid(),
	}
}

func index(w http.ResponseWriter, r *http.Request) {
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

func exit(exitCh chan<- struct{}) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
		exitCh <- struct{}{}
	}
}

func main() {
	flag.Parse()
	var logOut io.Writer = os.Stderr
	if logToStdout {
		logOut = os.Stdout
	}
	logger := log.New(logOut, "test-app: ", log.LstdFlags)

	doneCh := make(chan struct{})
	exitCh := make(chan struct{})

	mux := http.NewServeMux()
	mux.HandleFunc("/", index)
	mux.HandleFunc("/quit", exit(exitCh))

	server := http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      mux,
		ReadTimeout:  20 * time.Second,
		WriteTimeout: 20 * time.Second,
		IdleTimeout:  20 * time.Second,
	}

	go func() {
		defer close(doneCh)

		stopSignal := make(chan os.Signal, 1)
		if useSigusr1StopSignal {
			signal.Notify(stopSignal, syscall.SIGUSR1)
		} else {
			signal.Notify(stopSignal, syscall.SIGTERM)
		}

		select {

		case s := <-stopSignal:
			logger.Printf("signal %q: received\n", s)

			if sleepAfterStopSignal > 0 {
				logger.Printf("signal %q: sleeping for %v simulate cleanup\n", s, sleepAfterStopSignal)
				time.Sleep(sleepAfterStopSignal)
			}

		case <-time.After(stopAfter):
			logger.Printf("stopping after: %v\n", stopAfter)

		case <-exitCh:
			logger.Println("exiting after endpoint call")
		}

		if err := server.Shutdown(context.Background()); err != nil {
			log.Printf("server shutdown error: %v", err)
		}
	}()

	logger.Printf("server %s: started\n", server.Addr)

	if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		logger.Fatalf("could not start the server: %v", err)
	}

	logger.Printf("server %s: done\n", server.Addr)

	<-doneCh

	logger.Printf("exit code: %d\n", exitCode)

	os.Exit(exitCode)
}
