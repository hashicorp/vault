# Contributing

Thanks for your interest in improving this project! Before we get technical,
make sure you have reviewed the [code of conduct](code-of-conduct.md),
[Developer Certificate of Origin](DCO), and [OWNERS](OWNERS.md) files. Code will
be licensed according to [LICENSE.txt](LICENSE.txt).

## Pull Requests

When creating a pull request, please refer to an open issue. If there is no
issue open for the pull request you are creating, please create one. Frequently,
pull requests may be merged or closed while the underlying issue being addressed
is not fully addressed. Issues are a place to discuss the problem in need of a
solution. Pull requests are a place to discuss an implementation of one
particular answer to that problem.  A pull request may not address all (or any)
of the problems expressed in the issue, so it is important to track these
separately.

## Code Quality

### Documentation

All public functions and variables should include at least a short description
of the functionality they provide. Comments should be formatted according to
<https://golang.org/doc/effective_go.html#commentary>.

Documentation at <https://godoc.org/github.com/packethost/packngo> will be
generated from these comments.

The documentation provided for packngo fields and functions should be at or
better than the quality provided at <https://metal.equinix.com/developers/api/>.
When the API documentation provides a lengthy description, a linking to the
related API documentation will benefit users.

### Linters

`golangci-lint` is used to verify that the style of the code remains consistent.

Before committing, it's a good idea to run `goimports -w .`.
([goimports](https://pkg.go.dev/golang.org/x/tools/cmd/goimports?tab=doc)) and
`gofmt -w *.go`. ([gofmt](https://golang.org/cmd/gofmt/))

`make lint` can be used to verify style before creating a pull request.

## Building and Testing

The [Makefile](./Makefile) contains the targets to build, lint and test:

```sh
make build
make lint
make test
```

These normally will be run in a docker image of golang. To run locally, just run
with `BUILD=local`:

```sh
make build BUILD=local
make lint BUILD=local
make test BUILD=local
```

### Acceptance Tests

If you want to run tests against the actual Equinix Metal API, you must set the
environment variable `PACKET_TEST_ACTUAL_API` to a non-empty string and set
`PACKNGO_TEST_RECORDER` to `disabled`. The device tests wait for the device
creation, so it's best to run a few in parallel.

To run a particular test, you can do

```sh
PACKNGO_TEST_ACTUAL_API=1 go test -v -run=TestAccDeviceBasic --timeout=2h
```

If you want to see HTTP requests, set the `PACKNGO_DEBUG` env var to non-empty
string, for example:

```sh
PACKNGO_DEBUG=1 PACKNGO_TEST_ACTUAL_API=1 go test -v -run=TestAccVolumeUpdate
```

### Test Fixtures

By default, `go test ./...` will skip most of the tests unless
`PACKNGO_TEST_ACTUAL_API` is non-empty.

With the `PACKNGO_TEST_ACTUAL_API` environment variable set, tests will be run
against the Equinix Metal API, creating real infrastructure and incurring costs.

The `PACKNGO_TEST_RECORDER` variable can be used to record and playback API
responses to test code changes without the delay and costs of making actual API
calls. When unset, `PACKNGO_TEST_RECORDER` acts as though it was set to
`disabled`. This is the default behavior. This default behavior may change in
the future once fixtures are available for all tests.

When `PACKNGO_TEST_RECORDER` is set to `play`, tests will playback API responses
from recorded HTTP response fixtures. This is idea for refactoring and making
changes to request and response handling without introducing changes to the data
sent or received by the Equinix Metal API.

When adding support for new end-points, recorded test sessions should be added.
Record the HTTP interactions to fixtures by setting the environment variable
`PACKNGO_TEST_RECORDER` to `record`.

The fixtures are automatically named according to the test they were run from.
They are placed in `fixtures/`.  The API token used during authentication is
automatically removed from these fixtures. Nonetheless, caution should be
exercised before committing any fixtures into the project.  Account details
includes API tokens, contact, and payment details could easily be leaked by
committing fixtures that haven't been thoroughly reviewed.

### Automation (CI/CD)

Today, Drone tests pull requests using tests defined in
[.drone.yml](.drone.yml).

See [RELEASE.md](RELEASE.md) for details on the release process.
