## Overview

Package `statsd` provides a Go [dogstatsd](http://docs.datadoghq.com/guides/dogstatsd/) client.  Dogstatsd extends Statsd, adding tags
and histograms.

## Get the code

    $ go get github.com/DataDog/datadog-go/statsd

## Usage

```go
// Create the client
c, err := statsd.New("127.0.0.1:8125")
if err != nil {
    log.Fatal(err)
}
// Prefix every metric with the app name
c.Namespace = "flubber."
// Send the EC2 availability zone as a tag with every metric
c.Tags = append(c.Tags, "us-east-1a")

// Do some metrics!
err = c.Gauge("request.queue_depth", 12, nil, 1)
err = c.Timing("request.duration", duration, nil, 1) // Uses a time.Duration!
err = c.TimeInMilliseconds("request", 12, nil, 1)
err = c.Incr("request.count_total", nil, 1)
err = c.Decr("request.count_total", nil, 1)
err = c.Count("request.count_total", 2, nil, 1)
```

## Buffering Client

DogStatsD accepts packets with multiple statsd payloads in them.  Using the BufferingClient via `NewBufferingClient` will buffer up commands and send them when the buffer is reached or after 100msec.

## Unix Domain Sockets Client

DogStatsD version 6 accepts packets through a Unix Socket datagram connection. You can use this protocol by giving a
`unix:///path/to/dsd.socket` addr argument to the `New` or `NewBufferingClient`.

With this protocol, writes can become blocking if the server's receiving buffer is full. Our default behaviour is to
timeout and drop the packet after 1 ms. You can set a custom timeout duration via the `SetWriteTimeout` method.

The default mode is to pass write errors from the socket to the caller. This includes write errors the library will
automatically recover from (DogStatsD server not ready yet or is restarting). You can drop these errors and emulate
the UDP behaviour by setting the `SkipErrors` property to `true`. Please note that packets will be dropped in both modes.

## Development

Run the tests with:

    $ go test

## Documentation

Please see: http://godoc.org/github.com/DataDog/datadog-go/statsd

## License

go-dogstatsd is released under the [MIT license](http://www.opensource.org/licenses/mit-license.php).

## Credits

Original code by [ooyala](https://github.com/ooyala/go-dogstatsd).
