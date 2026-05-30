// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package seal_binary

import (
	"bufio"
	"fmt"
	"io"
	"maps"
	"os"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/helper/testcluster"
	"github.com/hashicorp/vault/sdk/helper/testcluster/docker"
)

func init() {
	if signed := os.Getenv("VAULT_LICENSE_CI"); signed != "" {
		if err := os.Setenv("VAULT_LICENSE", signed); err != nil {
			panic(err.Error())
		}
	}
}

func withPriorityAndDisabled(priority int, disabled bool, seal testcluster.VaultNodeSealConfig) testcluster.VaultNodeSealConfig {
	modified := seal
	modified.Config = maps.Clone(seal.Config)
	modified.Config["disabled"] = strconv.FormatBool(disabled)
	modified.Config["priority"] = strconv.Itoa(priority)
	return modified
}

type seal struct {
	base     func(name string, idx int) testcluster.VaultNodeSealConfig
	index    int
	disabled bool
	priority int
}

type step struct {
	expectedSealType string
	seals            []seal
}

func validateVaultStatusAndSealType(client *api.Client, expectedSealType string) error {
	statusResp, err := client.Sys().SealStatus()
	if err != nil {
		return fmt.Errorf("error getting vault status: %w", err)
	}

	if statusResp.Sealed {
		return fmt.Errorf("expected vault to be unsealed, but it is sealed")
	}

	if statusResp.Type != expectedSealType {
		return fmt.Errorf("unexpected seal type: expected %s, got %s", expectedSealType, statusResp.Type)
	}

	return nil
}

func dockerOptions(t *testing.T, repo, tag string) *docker.DockerClusterOptions {
	opts := docker.DefaultOptions(t)
	opts.NumCores = 1
	opts.ImageRepo, opts.ImageTag = repo, tag
	opts.VaultBinary = ""
	// Probably not reliable in CI with multi-node clusters, but we're assuming callers
	// of this func won't change NumCores to be >1.
	opts.VaultNodeConfig.StorageOptions = map[string]string{
		"performance_multiplier": "1",
	}
	return opts
}

type logScanner struct {
	wg   sync.WaitGroup
	l    sync.Mutex
	ch   chan string
	pw   *io.PipeWriter
	stop chan struct{}
}

func newLogScanner(t *testing.T, underlying io.Writer, bufLines int) (*logScanner, io.Writer) {
	pr, pw := io.Pipe()
	ls := &logScanner{
		ch:   make(chan string, bufLines),
		pw:   pw,
		stop: make(chan struct{}),
	}

	ls.wg.Add(1)
	go func() {
		defer ls.wg.Done()
		// bufio.Scanner is perfect here because hclog writes each log entry
		// ending with a newline character.
		scanner := bufio.NewScanner(pr)

		// scanner.Scan() will block until a new line is written to the pipe,
		// and it will exit automatically when pw.Close() is called.
		for scanner.Scan() {
			logLine := scanner.Text()
			underlying.Write([]byte(logLine + "\n"))
			select {
			case <-ls.stop:
				return
			case ls.ch <- logLine:
			}
		}

		if err := scanner.Err(); err != nil {
			t.Fatalf("Scanner error: %v", err)
		}
	}()

	t.Cleanup(func() {
		if err := ls.Close(); err != nil {
			t.Logf("Error closing scanner: %v", err)
		}
	})

	return ls, pw
}

func (ls *logScanner) Lines() <-chan string {
	return ls.ch
}

func (ls *logScanner) Close() error {
	ls.l.Lock()
	defer ls.l.Unlock()

	close(ls.stop)
	err := ls.pw.Close()
	ls.wg.Wait()

	return err
}

type logMatcher struct {
	targets map[string]bool
	lines   <-chan string
	done    chan struct{}
	l       sync.RWMutex
}

func newLogMatcher(lines <-chan string, targets []string) *logMatcher {
	tmap := make(map[string]bool)
	for _, target := range targets {
		tmap[target] = false
	}
	lm := &logMatcher{
		targets: tmap,
		lines:   lines,
		done:    make(chan struct{}),
	}

	go func() {
		for {
			select {
			case <-lm.done:
				return
			case line := <-lines:
				for target, ok := range tmap {
					if !ok && strings.Contains(line, target) {
						lm.l.Lock()
						tmap[target] = true
						lm.l.Unlock()
					}
				}
			}
		}
	}()

	return lm
}

func (lm *logMatcher) stop() {
	close(lm.done)
}

func (lm *logMatcher) missing() []string {
	lm.l.RLock()
	defer lm.l.RUnlock()

	var ret []string
	for target, found := range lm.targets {
		if !found {
			ret = append(ret, target)
		}
	}
	return ret
}
