// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package testing

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"strings"
	"sync"
	"testing"

	"go.uber.org/atomic"
)

const (
	ExpectedNamespace = "default"
	ExpectedPodName   = "shell-demo"
)

// Pull real-life-based testing data in from files at compile time.
// We decided to embed them in the test binary because of past issues
// with reading files that we encountered on CI workers.

//go:embed ca.crt
var caCrt string

//go:embed resp-get-pod.json
var getPodResponse string

//go:embed resp-not-found.json
var notFoundResponse string

//go:embed resp-update-pod.json
var updatePodTagsResponse string

//go:embed token
var token string

var (
	// ReturnGatewayTimeouts toggles whether the test server should return,
	// well, gateway timeouts...
	ReturnGatewayTimeouts = atomic.NewBool(false)

	pathToFiles = func() string {
		wd, _ := os.Getwd()
		repoName := "vault-enterprise"
		if !strings.Contains(wd, repoName) {
			repoName = "vault"
		}
		pathParts := strings.Split(wd, repoName)
		return pathParts[0] + "vault/serviceregistration/kubernetes/testing/"
	}()
)

// Conf returns the info needed to configure the client to point at
// the test server. This must be done by the caller to avoid an import
// cycle between the client and the testserver. Example usage:
//
//	client.Scheme = testConf.ClientScheme
//	client.TokenFile = testConf.PathToTokenFile
//	client.RootCAFile = testConf.PathToRootCAFile
//	if err := os.Setenv(client.EnvVarKubernetesServiceHost, testConf.ServiceHost); err != nil {
//		t.Fatal(err)
//	}
//	if err := os.Setenv(client.EnvVarKubernetesServicePort, testConf.ServicePort); err != nil {
//		t.Fatal(err)
//	}
type Conf struct {
	ClientScheme, PathToTokenFile, PathToRootCAFile, ServiceHost, ServicePort string
}

// Server returns an http test server that can be used to test
// Kubernetes client code. It also retains the current state,
// and a func to close the server and to clean up any temporary
// files.
func Server(t *testing.T) (testState *State, testConf *Conf, closeFunc func()) {
	testState = &State{m: &sync.Map{}}
	testConf = &Conf{
		ClientScheme: "http://",
	}

	// We're going to have multiple close funcs to call.
	var closers []func()
	closeFunc = func() {
		for _, closer := range closers {
			closer()
		}
	}

	// Plant our token in a place where it can be read for the config.
	tmpToken, err := ioutil.TempFile("", "token")
	if err != nil {
		t.Fatal(err)
	}
	closers = append(closers, func() {
		os.Remove(tmpToken.Name())
	})
	if _, err = tmpToken.WriteString(token); err != nil {
		closeFunc()
		t.Fatal(err)
	}
	if err := tmpToken.Close(); err != nil {
		closeFunc()
		t.Fatal(err)
	}
	testConf.PathToTokenFile = tmpToken.Name()

	tmpCACrt, err := ioutil.TempFile("", "ca.crt")
	if err != nil {
		closeFunc()
		t.Fatal(err)
	}
	closers = append(closers, func() {
		os.Remove(tmpCACrt.Name())
	})
	if _, err = tmpCACrt.WriteString(caCrt); err != nil {
		closeFunc()
		t.Fatal(err)
	}
	if err := tmpCACrt.Close(); err != nil {
		closeFunc()
		t.Fatal(err)
	}
	testConf.PathToRootCAFile = tmpCACrt.Name()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if ReturnGatewayTimeouts.Load() {
			w.WriteHeader(504)
			return
		}
		namespace, podName, err := parsePath(r.URL.Path)
		if err != nil {
			w.WriteHeader(400)
			w.Write([]byte(fmt.Sprintf("unable to parse %s: %s", r.URL.Path, err.Error())))
			return
		}

		switch {
		case namespace != ExpectedNamespace, podName != ExpectedPodName:
			w.WriteHeader(404)
			w.Write([]byte(notFoundResponse))
			return
		case r.Method == http.MethodGet:
			w.WriteHeader(200)
			w.Write([]byte(getPodResponse))
			return
		case r.Method == http.MethodPatch:
			var patches []interface{}
			if err := json.NewDecoder(r.Body).Decode(&patches); err != nil {
				w.WriteHeader(400)
				w.Write([]byte(fmt.Sprintf("unable to decode patches %s: %s", r.URL.Path, err.Error())))
				return
			}
			for _, patch := range patches {
				patchMap := patch.(map[string]interface{})
				p := patchMap["path"].(string)
				testState.store(p, patchMap)
			}
			w.WriteHeader(200)
			w.Write([]byte(updatePodTagsResponse))
			return
		default:
			w.WriteHeader(400)
			w.Write([]byte(fmt.Sprintf("unexpected request method: %s", r.Method)))
		}
	}))
	closers = append(closers, ts.Close)

	// ts.URL example: http://127.0.0.1:35681
	urlFields := strings.Split(ts.URL, "://")
	if len(urlFields) != 2 {
		closeFunc()
		t.Fatal("received unexpected test url: " + ts.URL)
	}
	urlFields = strings.Split(urlFields[1], ":")
	if len(urlFields) != 2 {
		closeFunc()
		t.Fatal("received unexpected test url: " + ts.URL)
	}
	testConf.ServiceHost = urlFields[0]
	testConf.ServicePort = urlFields[1]
	return testState, testConf, closeFunc
}

type State struct {
	m *sync.Map
}

func (s *State) NumPatches() int {
	l := 0
	f := func(key, value interface{}) bool {
		l++
		return true
	}
	s.m.Range(f)
	return l
}

func (s *State) Get(key string) map[string]interface{} {
	v, ok := s.m.Load(key)
	if !ok {
		return nil
	}
	patch, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}
	return patch
}

func (s *State) store(k string, p map[string]interface{}) {
	s.m.Store(k, p)
}

// The path should be formatted like this:
// fmt.Sprintf("/api/v1/namespaces/%s/pods/%s", namespace, podName)
func parsePath(urlPath string) (namespace, podName string, err error) {
	original := urlPath
	podName = path.Base(urlPath)
	urlPath = strings.TrimSuffix(urlPath, "/pods/"+podName)
	namespace = path.Base(urlPath)
	if original != fmt.Sprintf("/api/v1/namespaces/%s/pods/%s", namespace, podName) {
		return "", "", fmt.Errorf("received unexpected path: %s", original)
	}
	return namespace, podName, nil
}

func readFile(fileName string) (string, error) {
	b, err := ioutil.ReadFile(pathToFiles + fileName)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
