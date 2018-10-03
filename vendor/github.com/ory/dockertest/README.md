<h1 align="center"><img src="./docs/images/banner_dockertest.png" alt="ORY Dockertest"></h1>

[![Build Status](https://travis-ci.org/ory/dockertest.svg)](https://travis-ci.org/ory/dockertest?branch=master)
[![Coverage Status](https://coveralls.io/repos/github/ory/dockertest/badge.svg?branch=v3)](https://coveralls.io/github/ory/dockertest?branch=v3)

Use Docker to run your Go language integration tests against third party services on **Microsoft Windows, Mac OSX and Linux**!

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**

- [Why should I use Dockertest?](#why-should-i-use-dockertest)
- [Installing and using Dockertest](#installing-and-using-dockertest)
  - [Using Dockertest](#using-dockertest)
  - [Examples](#examples)
  - [Setting up Travis-CI](#setting-up-travis-ci)
- [Troubleshoot & FAQ](#troubleshoot-&-faq)
  - [Out of disk space](#out-of-disk-space)
  - [Removing old containers](#removing-old-containers)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## Why should I use Dockertest?

When developing applications, it is often necessary to use services that talk to a database system.
Unit Testing these services can be cumbersome because mocking database/DBAL is strenuous. Making slight changes to the
schema implies rewriting at least some, if not all of the mocks. The same goes for API changes in the DBAL.
To avoid this, it is smarter to test these specific services against a real database that is destroyed after testing.
Docker is the perfect system for running unit tests as you can spin up containers in a few seconds and kill them when
the test completes. The Dockertest library provides easy to use commands for spinning up Docker containers and using
them for your tests.

## Installing and using Dockertest

Using Dockertest is straightforward and simple. Check the [releases tab](https://github.com/ory/dockertest/releases)
for available releases.

To install dockertest, run

```
dep ensure -add github.com/ory/dockertest@v3.x.y
```

### Using Dockertest

```go
package dockertest_test

import (
	"testing"
	"log"
	"github.com/ory/dockertest"
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"fmt"
	"os"
)

var db *sql.DB

func TestMain(m *testing.M) {
	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
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
		log.Fatalf("Could not connect to docker: %s", err)
	}

	code := m.Run()
	
	// You can't defer this because os.Exit doesn't care for defer
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}
	
	os.Exit(code)
}

func TestSomething(t *testing.T) {
	// db.Query()
}
```

### Examples

We provide code examples for well known services in the [examples](examples/) directory, check them out!

### Setting up Travis-CI

You can run the Docker integration on Travis easily:

```yml
# Sudo is required for docker
sudo: required

# Enable docker
services:
  - docker
```

## Troubleshoot & FAQ

### Out of disk space

Try cleaning up the images with [docker-cleanup-volumes](https://github.com/chadoe/docker-cleanup-volumes).

### Removing old containers

Sometimes container clean up fails. Check out
[this stackoverflow question](http://stackoverflow.com/questions/21398087/how-to-delete-dockers-images) on how to fix this. You may also set an absolute lifetime on containers: 

```go
resource.Expire(60) // Tell docker to hard kill the container in 60 seconds
```