<h1 align="center"><img src="./docs/images/banner_dockertest.png" alt="ORY Dockertest"></h1>

[![Build Status](https://travis-ci.org/ory/dockertest.svg)](https://travis-ci.org/ory/dockertest?branch=master)
[![Coverage Status](https://coveralls.io/repos/github/ory/dockertest/badge.svg?branch=v3)](https://coveralls.io/github/ory/dockertest?branch=v3)

Use Docker to run your Golang integration tests against third party services on
**Microsoft Windows, Mac OSX and Linux**!

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

**Table of Contents**

- [Why should I use Dockertest?](#why-should-i-use-dockertest)
- [Installing and using Dockertest](#installing-and-using-dockertest)
  - [Using Dockertest](#using-dockertest)
  - [Examples](#examples)
- [Troubleshoot & FAQ](#troubleshoot--faq)
  - [Out of disk space](#out-of-disk-space)
  - [Removing old containers](#removing-old-containers)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## Why should I use Dockertest?

When developing applications, it is often necessary to use services that talk to
a database system. Unit Testing these services can be cumbersome because mocking
database/DBAL is strenuous. Making slight changes to the schema implies
rewriting at least some, if not all of the mocks. The same goes for API changes
in the DBAL. To avoid this, it is smarter to test these specific services
against a real database that is destroyed after testing. Docker is the perfect
system for running unit tests as you can spin up containers in a few seconds and
kill them when the test completes. The Dockertest library provides easy to use
commands for spinning up Docker containers and using them for your tests.

## Installing and using Dockertest

Using Dockertest is straightforward and simple. Check the
[releases tab](https://github.com/ory/dockertest/releases) for available
releases.

To install dockertest, run

```
go get -u github.com/ory/dockertest/v3
```

### Using Dockertest

```go
package dockertest_test

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/ory/dockertest/v3"
)

var db *sql.DB

func TestMain(m *testing.M) {
	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not construct pool: %s", err)
	}

	// uses pool to try to connect to Docker
	err = pool.Client.Ping()
	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}

	// pulls an image, creates a container based on it and runs it
	resource, err := pool.Run("mysql", "5.7", []string{"MYSQL_ROOT_PASSWORD=secret"})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	if err := pool.Retry(func() error {
		var err error
		db, err = sql.Open("mysql", fmt.Sprintf("root:secret@(localhost:%s)/mysql", resource.GetPort("3306/tcp")))
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to database: %s", err)
	}

	// as of go1.15 testing.M returns the exit code of m.Run(), so it is safe to use defer here
    defer func() {
      if err := pool.Purge(resource); err != nil {
        log.Fatalf("Could not purge resource: %s", err)
      }

    }()

	m.Run()
}

func TestSomething(t *testing.T) {
	// db.Query()
}
```

### Examples

We provide code examples for well known services in the [examples](examples/)
directory, check them out!

## Troubleshoot & FAQ

### Out of disk space

Try cleaning up the images with
[docker-cleanup-volumes](https://github.com/chadoe/docker-cleanup-volumes).

### Removing old containers

Sometimes container clean up fails. Check out
[this stackoverflow question](http://stackoverflow.com/questions/21398087/how-to-delete-dockers-images)
on how to fix this. You may also set an absolute lifetime on containers:

```go
resource.Expire(60) // Tell docker to hard kill the container in 60 seconds
```

To let stopped containers removed from file system automatically, use
`pool.RunWithOptions()` instead of `pool.Run()` with `config.AutoRemove` set to
true, e.g.:

```go
postgres, err := pool.RunWithOptions(&dockertest.RunOptions{
	Repository: "postgres",
	Tag:        "11",
	Env: []string{
		"POSTGRES_USER=test",
		"POSTGRES_PASSWORD=test",
		"listen_addresses = '*'",
	},
}, func(config *docker.HostConfig) {
	// set AutoRemove to true so that stopped container goes away by itself
	config.AutoRemove = true
	config.RestartPolicy = docker.RestartPolicy{
		Name: "no",
	}
})
```

## Running dockertest in Gitlab CI

### How to run dockertest on shared gitlab runners?

You should add docker dind service to your job which starts in sibling
container. That means database will be available on host `docker`.  
You app should be able to change db host through environment variable.

Here is the simple example of `gitlab-ci.yml`:

```yaml
stages:
  - test
go-test:
  stage: test
  image: golang:1.15
  services:
    - docker:dind
  variables:
    DOCKER_HOST: tcp://docker:2375
    DOCKER_DRIVER: overlay2
    YOUR_APP_DB_HOST: docker
  script:
    - go test ./...
```

Plus in the `pool.Retry` method that checks for connection readiness, you need
to use `$YOUR_APP_DB_HOST` instead of localhost.

### How to run dockertest on group(custom) gitlab runners?

Gitlab runner can be run in docker executor mode to save compatibility with
shared runners.  
Here is the simple register command:

```shell script
gitlab-runner register -n \
 --url https://gitlab.com/ \
 --registration-token $YOUR_TOKEN \
 --executor docker \
 --description "My Docker Runner" \
 --docker-image "docker:19.03.12" \
 --docker-privileged
```

You only need to instruct docker dind to start with disabled tls.  
Add variable `DOCKER_TLS_CERTDIR: ""` to `gitlab-ci.yml` above. It will tell
docker daemon to start on 2375 port over http.

## Running Dockertest Using GitHub Actions

```yaml
name: Test with Docker

on: [push]

jobs:
  test:
    runs-on: ubuntu-latest
    services:
      dind:
        image: docker:23.0-rc-dind-rootless
        ports:
          - 2375:2375
    steps:
      - name: Checkout code
        uses: actions/checkout@v2

      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: "1.21"

      - name: Test with Docker
        run: go test -v ./...
```

### How to run dockertest with remote Docker

Use-case: locally installed docker CLI (client), docker daemon somewhere
remotely, environment properly set (ie: `DOCKER_HOST`, etc..). For example,
remote docker can be provisioned by docker-machine.

Currently, dockertest in case of `resource.GetHostPort()` will return docker
host binding address (commonly - `localhost`) instead of remote docker host.
Universal solution is:

```go
func getHostPort(resource *dockertest.Resource, id string) string {
	dockerURL := os.Getenv("DOCKER_HOST")
	if dockerURL == "" {
		return resource.GetHostPort(id)
	}
	u, err := url.Parse(dockerURL)
	if err != nil {
		panic(err)
	}
	return u.Hostname() + ":" + resource.GetPort(id)
}
```

It will return the remote docker host concatenated with the allocated port in
case `DOCKER_HOST` env is defined. Otherwise, it will fall back to the embedded
behavior.
