[![GoDoc](https://godoc.org/github.com/xdg/stringprep?status.svg)](https://godoc.org/github.com/xdg/stringprep)
[![Build Status](https://travis-ci.org/xdg/stringprep.svg?branch=master)](https://travis-ci.org/xdg/stringprep)

# stringprep – Go implementation of RFC-3454 stringprep and RFC-4013 SASLprep

## Synopsis

```
    import "github.com/xdg/stringprep"

    prepped := stringprep.SASLprep.Prepare("TrustNô1")

```

## Description

This library provides an implementation of the stringprep algorithm
(RFC-3454) in Go, including all data tables.

A pre-built SASLprep (RFC-4013) profile is provided as well.

## Copyright and License

Copyright 2018 by David A. Golden. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License"). You may
obtain a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
