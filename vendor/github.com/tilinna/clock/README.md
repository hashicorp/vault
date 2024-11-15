# clock [![GoDoc](https://godoc.org/github.com/tilinna/clock?status.png)](https://godoc.org/github.com/tilinna/clock) [![Go Report Card](https://goreportcard.com/badge/github.com/tilinna/clock)](https://goreportcard.com/report/github.com/tilinna/clock)

A Go (golang) library for mocking standard time, optionally also with context.Context.

## Basic Usage

```go
// Use clock.Realtime() in production
mock := clock.NewMock(time.Date(2018, 1, 1, 10, 0, 0, 0, time.UTC))
fmt.Println("Time is now", mock.Now())
timer := mock.NewTimer(15 * time.Second)
mock.Add(25 * time.Second)
fmt.Println("Time is now", mock.Now())
fmt.Println("Timeout was", <-timer.C)
// Output:
// Time is now 2018-01-01 10:00:00 +0000 UTC
// Time is now 2018-01-01 10:00:25 +0000 UTC
// Timeout was 2018-01-01 10:00:15 +0000 UTC
```

## Context Usage

```go
start := time.Date(2018, 1, 1, 10, 0, 0, 0, time.UTC)
mock := clock.NewMock(start)
fmt.Println("now:", mock.Now())
ctx, cfn := mock.DeadlineContext(context.Background(), start.Add(time.Hour))
defer cfn()
fmt.Println("err:", ctx.Err())
dl, _ := ctx.Deadline()
mock.Set(dl)
fmt.Println("now:", clock.Now(ctx))
<-ctx.Done()
fmt.Println("err:", ctx.Err())
// Output:
// now: 2018-01-01 10:00:00 +0000 UTC
// err: <nil>
// now: 2018-01-01 11:00:00 +0000 UTC
// err: context deadline exceeded
```
