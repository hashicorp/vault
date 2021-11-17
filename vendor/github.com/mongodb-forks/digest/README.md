[![Go Reference](https://pkg.go.dev/badge/github.com/mongodb-forks/digest.svg)](https://pkg.go.dev/github.com/mongodb-forks/digest)
[![GO tests](https://github.com/mongodb-forks/digest/actions/workflows/go-test.yml/badge.svg)](https://github.com/mongodb-forks/digest/actions/workflows/go-test.yml)
[![golangci-lint](https://github.com/mongodb-forks/digest/actions/workflows/golangci-lint.yml/badge.svg)](https://github.com/mongodb-forks/digest/actions/workflows/golangci-lint.yml)
# Golang HTTP Digest Authentication

## Overview

This is a fork of the (unmaintained) code.google.com/p/mlab-ns2/gae/ns/digest package.
There's a descriptor leak in the original package, so this fork was created to patch
the leak.

### Update 2020

This is a fork of the now unmaintained fork of [digest](https://github.com/bobziuchkovski/digest).
This implementation now supports the SHA-256 algorithm which was added as part of [rfc 7616](https://tools.ietf.org/html/rfc7616).

## Usage

```go
t := NewTransport("myUserName", "myP@55w0rd")
req, err := http.NewRequest("GET", "http://notreal.com/path?arg=1", nil)
if err != nil {
	return err
}
resp, err := t.RoundTrip(req)
if err != nil {
	return err
}
```
Or it can be used as a client:
```go
c, err := t.Client()
if err != nil {
	return err
}
resp, err := c.Get("http://notreal.com/path?arg=1")
if err != nil {
	return err
}
```
## Contributing

**Contributions are welcome!**

The code is linted with [golangci-lint](https://golangci-lint.run/).  This library also defines *git hooks* that format and lint the code.

Before submitting a PR, please run `make setup link-git-hooks` to set up your local development environment.

## Original Authors

- Bipasa Chattopadhyay <bipasa@cs.unc.edu>
- Eric Gavaletz <gavaletz@gmail.com>
- Seon-Wook Park <seon.wook@swook.net>
- Bob Ziuchkovski (@bobziuchkovski)

## License

Apache 2.0
