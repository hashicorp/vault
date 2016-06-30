# [ory.am](https://ory.am)/dockertest

[![Build Status](https://travis-ci.org/ory-am/dockertest.svg)](https://travis-ci.org/ory-am/dockertest?branch=master)
[![Coverage Status](https://coveralls.io/repos/ory-am/dockertest/badge.svg?branch=master&service=github)](https://coveralls.io/github/ory-am/dockertest?branch=master)

Use Docker to run your Go language integration tests against third party services on **Microsoft Windows, Mac OSX and Linux**!
Dockertest uses [docker-machine](https://docs.docker.com/machine/) (aka [Docker Toolbox](https://www.docker.com/toolbox)) to spin up images on Windows and Mac OSX as well.
Dockertest is based on [docker.go](https://github.com/camlistore/camlistore/blob/master/pkg/test/dockertest/docker.go)
from [camlistore](https://github.com/camlistore/camlistore).

This fork detects automatically, if [Docker Toolbox](https://www.docker.com/toolbox)
is installed. If it is, Docker integration on Windows and Mac OSX can be used without any additional work.
To avoid port collisions when using docker-machine, Dockertest chooses a random port to bind the requested image.

Dockertest ships with support for these backends:
* PostgreSQL
* MySQL
* MongoDB
* NSQ
* Redis
* Elastic Search
* RethinkDB
* RabbitMQ
* Mockserver
* ActiveMQ

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**

- [Why should I use Dockertest?](#why-should-i-use-dockertest)
- [Installing and using Dockertest](#installing-and-using-dockertest)
  - [Start a container](#start-a-container)
- [Write awesome tests](#write-awesome-tests)
  - [Setting up Travis-CI](#setting-up-travis-ci)
- [Troubleshoot & FAQ](#troubleshoot-&-faq)
  - [I want to use a specific image version](#i-want-to-use-a-specific-image-version)
  - [My build is broken!](#my-build-is-broken)
  - [Out of disk space](#out-of-disk-space)
  - [Removing old containers](#removing-old-containers)
  - [Customized database] (#Customized-database)

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

Using Dockertest is straightforward and simple. Check the [releases tab](https://github.com/ory-am/dockertest/releases)
for available releases.

To install dockertest, run

```
go get gopkg.in/ory-am/dockertest.vX
```

where `X` is your desired version. For example:

```
go get gopkg.in/ory-am/dockertest.v2
```

**Note:**  
When using the Docker Toolbox (Windows / OSX), make sure that the VM is started by running `docker-machine start default`.

### Start a container

```go
package main

import (
	"gopkg.in/ory-am/dockertest.v2"
	"gopkg.in/mgo.v2"
	"time"
)

func main() {
	var db *mgo.Session
	c, err := dockertest.ConnectToMongoDB(15, time.Millisecond*500, func(url string) bool {
	    // This callback function checks if the image's process is responsive.
	    // Sometimes, docker images are booted but the process (in this case MongoDB) is still doing maintenance
	    // before being fully responsive which might cause issues like "TCP Connection reset by peer".	
		var err error
		db, err = mgo.Dial(url)
		if err != nil {
			return false
		}
		
		// Sometimes, dialing the database is not enough because the port is already open but the process is not responsive.
		// Most database conenctors implement a ping function which can be used to test if the process is responsive.
		// Alternatively, you could execute a query to see if an error occurs or not.
		return db.Ping() == nil
	})
	
	if err != nil {
	    log.Fatalf("Could not connect to database: %s", err)
	}
	
	// Close db connection and kill the container when we leave this function body.
    defer db.Close()
	defer c.KillRemove()
	
	// The image is now responsive.
}
```

You can start PostgreSQL and MySQL in a similar fashion.

It is also possible to start a custom container (in this example, a RabbitMQ container):

```go
	c, ip, port, err := dockertest.SetupCustomContainer("rabbitmq", 5672, 10*time.Second)
	if err != nil {
		log.Fatalf("Could not setup container: %s", err
	}
	defer c.KillRemove()

	err = dockertest.ConnectToCustomContainer(fmt.Sprintf("%v:%v", ip, port), 15, time.Millisecond*500, func(url string) bool {
		amqp, err := amqp.Dial(fmt.Sprintf("amqp://%v", url))
		if err != nil {
			return false
		}
		defer amqp.Close()
		return true
	})

	...
```

## Write awesome tests

It is a good idea to start up the container only once when running tests.

```go

import (
	"fmt"
	"testing"
    "log"
	"os"

	"database/sql"
	_ "github.com/lib/pq"
	"gopkg.in/ory-am/dockertest.v2"
)

var db *sql.DB

func TestMain(m *testing.M) {
	c, err := dockertest.ConnectToPostgreSQL(15, time.Second, func(url string) bool {
	    // Check if postgres is responsive...
		var err error
		db, err = sql.Open("postgres", url)
		if err != nil {
			return false
		}
		return db.Ping() == nil
	})
	if err != nil {
		log.Fatalf("Could not connect to database: %s", err)
	}
	
	// Execute tasks like setting up schemata.
	
	// Run tests
	result := m.Run()
	
	// Close database connection.
	db.Close()
	
	// Clean up image.
	c.KillRemove()
	
	// Exit tests.
	os.Exit(result)
}

func TestFunction(t *testing.T) {
    // db.Exec(...
}
```

### Setting up Travis-CI

You can run the Docker integration on Travis easily:

```yml
# Sudo is required for docker
sudo: required

# Enable docker
services:
  - docker

# In Travis, we need to bind to 127.0.0.1 in order to get a working connection. This environment variable
# tells dockertest to do that.
env:
  - DOCKERTEST_BIND_LOCALHOST=true
```

## Troubleshoot & FAQ

### I need to use a specific container version for XYZ

You can specify a container version by setting environment variables or globals. For more information, check [vars.go](vars.go).

### My build is broken!

With v2, we removed all `Open*` methods to reduce duplicate code, unnecessary dependencies and make maintenance easier.
If you relied on these, run `go get gopkg.in/ory-am/dockertest.v1` and replace
`import "github.com/ory-am/dockertest"` with `import "gopkg.in/ory-am/dockertest.v1"`.

### Out of disk space

Try cleaning up the images with [docker-cleanup-volumes](https://github.com/chadoe/docker-cleanup-volumes).

### Removing old containers

Sometimes container clean up fails. Check out
[this stackoverflow question](http://stackoverflow.com/questions/21398087/how-to-delete-dockers-images) on how to fix this.

### Customized database

I am using postgres (or mysql) driver, how do I use customized database instead of default one?
You can alleviate this helper function to do that, see testcase or example below:

```go

func TestMain(m *testing.M) {
	if c, err := dockertest.ConnectToPostgreSQL(15, time.Second, func(url string) bool {
        customizedDB := "cherry" // here I am connecting cherry database
        newURL, err := SetUpPostgreDatabase(customizedDB, url)

        // or use SetUpMysqlDatabase for mysql driver

        if err != nil {
                log.Fatal(err)
        }
        db, err := sql.Open("postgres", newURL)
        if err != nil {
            return false
        }
        return db.Ping() == nil
    }); err != nil {
        log.Fatal(err)
    }
```

*Thanks to our sponsors: Ory GmbH & Imarum GmbH*
