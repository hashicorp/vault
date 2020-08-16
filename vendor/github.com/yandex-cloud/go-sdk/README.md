# Yandex.Cloud Go SDK

[![GoDoc](https://godoc.org/github.com/yandex-cloud/go-sdk?status.svg)](https://godoc.org/github.com/yandex-cloud/go-sdk)
[![CircleCI](https://circleci.com/gh/yandex-cloud/go-sdk.svg?style=shield)](https://circleci.com/gh/yandex-cloud/go-sdk)

Go SDK for Yandex.Cloud services.

**NOTE:** SDK is under development, and may make
backwards-incompatible changes.

## Installation

```bash
go get github.com/yandex-cloud/go-sdk
```

## Example usages

### Initializing SDK

```go
sdk, err := ycsdk.Build(ctx, ycsdk.Config{
	Credentials: ycsdk.OAuthToken(token),
})
if err != nil {
	log.Fatal(err)
}
```

### More examples

More examples can be found in [examples dir](examples).
