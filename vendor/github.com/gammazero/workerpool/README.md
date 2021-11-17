# workerpool
[![Build Status](https://travis-ci.org/gammazero/workerpool.svg)](https://travis-ci.org/gammazero/workerpool)
[![Go Report Card](https://goreportcard.com/badge/github.com/gammazero/workerpool)](https://goreportcard.com/report/github.com/gammazero/workerpool)
[![codecov](https://codecov.io/gh/gammazero/workerpool/branch/master/graph/badge.svg)](https://codecov.io/gh/gammazero/workerpool)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](https://github.com/gammazero/workerpool/blob/master/LICENSE)

Concurrency limiting goroutine pool. Limits the concurrency of task execution, not the number of tasks queued. Never blocks submitting tasks, no matter how many tasks are queued.

[![GoDoc](https://godoc.org/github.com/gammazero/workerpool?status.svg)](https://godoc.org/github.com/gammazero/workerpool)

This implementation builds on ideas from the following:

- http://marcio.io/2015/07/handling-1-million-requests-per-minute-with-golang
- http://nesv.github.io/golang/2014/02/25/worker-queues-in-go.html

## Installation
To install this package, you need to setup your Go workspace.  The simplest way to install the library is to run:
```
$ go get github.com/gammazero/workerpool
```

## Example
```go
package main

import (
	"fmt"
	"github.com/gammazero/workerpool"
)

func main() {
	wp := workerpool.New(2)
	requests := []string{"alpha", "beta", "gamma", "delta", "epsilon"}

	for _, r := range requests {
		r := r
		wp.Submit(func() {
			fmt.Println("Handling request:", r)
		})
	}

	wp.StopWait()
}
```
