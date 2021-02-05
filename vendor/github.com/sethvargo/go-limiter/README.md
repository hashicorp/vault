# Go Rate Limiter

[![GoDoc](https://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](https://pkg.go.dev/mod/github.com/sethvargo/go-limiter)
[![GitHub Actions](https://img.shields.io/github/workflow/status/sethvargo/go-limiter/Test?style=flat-square)](https://github.com/sethvargo/go-limiter/actions?query=workflow%3ATest)


This package provides a rate limiter in Go (Golang), suitable for use in HTTP
servers and distributed workloads. It's specifically designed for
configurability and flexibility without compromising throughput.


## Usage

1.  Create a store. This example uses an in-memory store:

    ```golang
    store, err := memorystore.New(&memorystore.Config{
      // Number of tokens allowed per interval.
      Tokens: 15,

      // Interval until tokens reset.
      Interval: time.Minute,
    })
    if err != nil {
      log.Fatal(err)
    }
    ```

1.  Determine the limit by calling `Take()` on the store:

    ```golang
    ctx := context.Background()

    // key is the unique value upon which you want to rate limit, like an IP or
    // MAC address.
    key := "127.0.0.1"
    tokens, remaining, reset, ok, err := store.Take(ctx, key)

    // tokens is the configured tokens (15 in this example).
    _ = tokens

    // remaining is the number of tokens remaining (14 now).
    _ = remaining

    // reset is the unix nanoseconds at which the tokens will replenish.
    _ = reset

    // ok indicates whether the take was successful. If the key is over the
    // configured limit, ok will be false.
    _ = ok

    // Here's a more realistic example:
    if !ok {
      return fmt.Errorf("rate limited: retry at %v", reset)
    }
    ```

There's also HTTP middleware via the `httplimit` package. After creating a
store, wrap Go's standard HTTP handler:

```golang
middleware, err := httplimit.NewMiddleware(store, httplimit.IPKeyFunc())
if err != nil {
  log.Fatal(err)
}

mux1 := http.NewServeMux()
mux1.Handle("/", middleware.Handle(doWork)) // doWork is your original handler
```

The middleware automatically set the following headers, conforming to the latest
RFCs:

- `X-RateLimit-Limit` - configured rate limit (constant).
- `X-RateLimit-Remaining` - number of remaining tokens in current interval.
- `X-RateLimit-Reset` - UTC time when the limit resets.
- `Retry-After` - Time at which to retry


## Why _another_ Go rate limiter?

I really wanted to learn more about the topic and possibly implementations. The
existing packages in the Go ecosystem either lacked flexibility or traded
flexibility for performance. I wanted to write a package that was highly
extensible while still offering the highest levels of performance.


### Speed and performance

How fast is it? You can run the benchmarks yourself, but here's a few sample
benchmarks with 100,000 unique keys. I added commas to the output for clarity,
but you can run the benchmarks via `make benchmarks`:

```text
$ make benchmarks
BenchmarkSethVargoMemory/memory/serial-7      13,706,899      81.7 ns/op       16 B/op     1 allocs/op
BenchmarkSethVargoMemory/memory/parallel-7     7,900,639       151 ns/op       61 B/op     3 allocs/op
BenchmarkSethVargoMemory/sweep/serial-7       19,601,592      58.3 ns/op        0 B/op     0 allocs/op
BenchmarkSethVargoMemory/sweep/parallel-7     21,042,513      55.2 ns/op        0 B/op     0 allocs/op
BenchmarkThrottled/memory/serial-7             6,503,260       176 ns/op        0 B/op     0 allocs/op
BenchmarkThrottled/memory/parallel-7           3,936,655       297 ns/op        0 B/op     0 allocs/op
BenchmarkThrottled/sweep/serial-7              6,901,432       171 ns/op        0 B/op     0 allocs/op
BenchmarkThrottled/sweep/parallel-7            5,948,437       202 ns/op        0 B/op     0 allocs/op
BenchmarkTollbooth/memory/serial-7             3,064,309       368 ns/op        0 B/op     0 allocs/op
BenchmarkTollbooth/memory/parallel-7           2,658,014       448 ns/op        0 B/op     0 allocs/op
BenchmarkTollbooth/sweep/serial-7              2,769,937       430 ns/op      192 B/op     3 allocs/op
BenchmarkTollbooth/sweep/parallel-7            2,216,211       546 ns/op      192 B/op     3 allocs/op
BenchmarkUber/memory/serial-7                 13,795,612      94.2 ns/op        0 B/op     0 allocs/op
BenchmarkUber/memory/parallel-7                7,503,214       159 ns/op        0 B/op     0 allocs/op
BenchmarkUlule/memory/serial-7                 2,964,438       405 ns/op       24 B/op     2 allocs/op
BenchmarkUlule/memory/parallel-7               2,441,778       469 ns/op       24 B/op     2 allocs/op
```

There's likely still optimizations to be had, pull requests are welcome!


### Ecosystem

Many of the existing packages in the ecosystem take dependencies on other
packages. I'm an advocate of very thin libraries, and I don't think a rate
limiter should be pulling external packages. That's why **go-limit uses only the
Go standard library**.


### Flexible and extensible

Most of the existing rate limiting libraries make a strong assumption that rate
limiting is only for HTTP services. Baked in that assumption are more
assumptions like rate limiting by "IP address" or are limited to a resolution of
"per second". While go-limit supports rate limiting at the HTTP layer, it can
also be used to rate limit literally anything. It rate limits on a user-defined
arbitrary string key.


### Stores

#### Memory

Memory is the fastest store, but only works on a single container/virtual
machine since there's no way to share the state.
[Learn more](https://pkg.go.dev/github.com/sethvargo/go-limiter/memorystore).

#### Redis

Redis uses Redis + Lua as a shared pool, but comes at a performance cost.
[Learn more](https://pkg.go.dev/github.com/sethvargo/go-redisstore).

#### Noop

Noop does no rate limiting, but still implements the interface - useful for
testing and local development.
[Learn more](https://pkg.go.dev/github.com/sethvargo/go-limiter/noopstore).
