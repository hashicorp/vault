# Circonus metrics tracking for Go applications

This library supports named counters, gauges and histograms.
It also provides convenience wrappers for registering latency
instrumented functions with Go's builtin http server.

Initializing only requires setting an ApiToken.

## Example

**rough and simple**

```go
package main

import (
	"log"
	"math/rand"
	"os"
	"time"

	cgm "github.com/circonus-labs/circonus-gometrics"
)

func main() {

    logger := log.New(os.Stdout, "", log.LstdFlags)

	logger.Println("Configuring cgm")

	cmc := &cgm.Config{}

	// Interval at which metrics are submitted to Circonus, default: 10 seconds
	// cmc.Interval = "10s" // 10 seconds

	// Enable debug messages, default: false
	cmc.Debug = true

	// Send debug messages to specific log.Logger instance
	// default: if debug stderr, else, discard
	cmc.Log = logger

    // Reset counter metrics after each submission, default: "true"
    // Change to "false" to retain (and continue submitting) the last value.
    // cmc.ResetCounters = "true"

    // Reset gauge metrics after each submission, default: "true"
    // Change to "false" to retain (and continue submitting) the last value.
    // cmc.ResetGauges = "true"

    // Reset histogram metrics after each submission, default: "true"
    // Change to "false" to retain (and continue submitting) the last value.
    // cmc.ResetHistograms = "true"

    // Reset text metrics after each submission, default: "true"
    // Change to "false" to retain (and continue submitting) the last value.
    // cmc.ResetText = "true"

	// Circonus API configuration options
	//
	// Token, no default (blank disables check manager)
	cmc.CheckManager.API.TokenKey = os.Getenv("CIRCONUS_API_TOKEN")
	// App name, default: circonus-gometrics
	cmc.CheckManager.API.TokenApp = os.Getenv("CIRCONUS_API_APP")
	// URL, default: https://api.circonus.com/v2
	cmc.CheckManager.API.URL = os.Getenv("CIRCONUS_API_URL")

	// Check configuration options
	//
	// precedence 1 - explicit submission_url
	// precedence 2 - specific check id (note: not a check bundle id)
	// precedence 3 - search using instanceId and searchTag
	// otherwise: if an applicable check is NOT specified or found, an
	//            attempt will be made to automatically create one
	//
	// Submission URL for an existing [httptrap] check
	cmc.CheckManager.Check.SubmissionURL = os.Getenv("CIRCONUS_SUBMISION_URL")

	// ID of an existing [httptrap] check (note: check id not check bundle id)
	cmc.CheckManager.Check.ID = os.Getenv("CIRCONUS_CHECK_ID")

	// if neither a submission url nor check id are provided, an attempt will be made to find an existing
	// httptrap check by using the circonus api to search for a check matching the following criteria:
	//      an active check,
	//      of type httptrap,
	//      where the target/host is equal to InstanceId - see below
	//      and the check has a tag equal to SearchTag - see below
	// Instance ID - an identifier for the 'group of metrics emitted by this process or service'
	//               this is used as the value for check.target (aka host)
	// default: 'hostname':'program name'
	// note: for a persistent instance that is ephemeral or transient where metric continuity is
	//       desired set this explicitly so that the current hostname will not be used.
	// cmc.CheckManager.Check.InstanceID = ""

	// Search tag - specific tag(s) used in conjunction with isntanceId to search for an
    // existing check. comma separated string of tags (spaces will be removed, no commas
    // in tag elements).
	// default: service:application name (e.g. service:consul service:nomad etc.)
	// cmc.CheckManager.Check.SearchTag = ""

	// Check secret, default: generated when a check needs to be created
	// cmc.CheckManager.Check.Secret = ""

	// Additional tag(s) to add when *creating* a check. comma separated string
    // of tags (spaces will be removed, no commas in tag elements).
    // (e.g. group:abc or service_role:agent,group:xyz).
    // default: none
	// cmc.CheckManager.Check.Tags = ""

	// max amount of time to to hold on to a submission url
	// when a given submission fails (due to retries) if the
	// time the url was last updated is > than this, the trap
	// url will be refreshed (e.g. if the broker is changed
	// in the UI) default 5 minutes
	// cmc.CheckManager.Check.MaxURLAge = "5m"

	// custom display name for check, default: "InstanceId /cgm"
	// cmc.CheckManager.Check.DisplayName = ""

    // force metric activation - if a metric has been disabled via the UI
	// the default behavior is to *not* re-activate the metric; this setting
	// overrides the behavior and will re-activate the metric when it is
	// encountered. "(true|false)", default "false"
	// cmc.CheckManager.Check.ForceMetricActivation = "false"

	// Broker configuration options
	//
	// Broker ID of specific broker to use, default: random enterprise broker or
	// Circonus default if no enterprise brokers are available.
	// default: only used if set
	// cmc.CheckManager.Broker.ID = ""

	// used to select a broker with the same tag(s) (e.g. can be used to dictate that a broker
	// serving a specific location should be used. "dc:sfo", "loc:nyc,dc:nyc01", "zone:us-west")
	// if more than one broker has the tag(s), one will be selected randomly from the resulting
    // list. comma separated string of tags (spaces will be removed, no commas in tag elements).
	// default: none
	// cmc.CheckManager.Broker.SelectTag = ""

	// longest time to wait for a broker connection (if latency is > the broker will
	// be considered invalid and not available for selection.), default: 500 milliseconds
	// cmc.CheckManager.Broker.MaxResponseTime = "500ms"

	// note: if broker Id or SelectTag are not specified, a broker will be selected randomly
	// from the list of brokers available to the api token. enterprise brokers take precedence
	// viable brokers are "active", have the "httptrap" module enabled, are reachable and respond
	// within MaxResponseTime.

	logger.Println("Creating new cgm instance")

	metrics, err := cgm.NewCirconusMetrics(cmc)
	if err != nil {
		panic(err)
	}

	src := rand.NewSource(time.Now().UnixNano())
	rnd := rand.New(src)

	logger.Println("Starting cgm internal auto-flush timer")
	metrics.Start()

    logger.Println("Adding ctrl-c trap")
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		logger.Println("Received CTRL-C, flushing outstanding metrics before exit")
		metrics.Flush()
		os.Exit(0)
	}()

    // Add metric tags (append to any existing tags on specified metric)
    metrics.AddMetricTags("foo", []string{"cgm:test"})
    metrics.AddMetricTags("baz", []string{"cgm:test"})

	logger.Println("Starting to send metrics")

	// number of "sets" of metrics to send
	max := 60

	for i := 1; i < max; i++ {
		logger.Printf("\tmetric set %d of %d", i, 60)

		metrics.Timing("foo", rnd.Float64()*10)
		metrics.Increment("bar")
		metrics.Gauge("baz", 10)

        if i == 35 {
            // Set metric tags (overwrite current tags on specified metric)
            metrics.SetMetricTags("baz", []string{"cgm:reset_test", "cgm:test2"})
        }

        time.Sleep(time.Second)
	}

	logger.Println("Flushing any outstanding metrics manually")
	metrics.Flush()

}
```

### HTTP Handler wrapping

```
http.HandleFunc("/", metrics.TrackHTTPLatency("/", handler_func))
```

### HTTP latency example

```
package main

import (
    "os"
    "fmt"
    "net/http"
    cgm "github.com/circonus-labs/circonus-gometrics"
)

func main() {
    cmc := &cgm.Config{}
    cmc.CheckManager.API.TokenKey = os.Getenv("CIRCONUS_API_TOKEN")

    metrics, err := cgm.NewCirconusMetrics(cmc)
    if err != nil {
        panic(err)
    }
    metrics.Start()

    http.HandleFunc("/", metrics.TrackHTTPLatency("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
    }))
    http.ListenAndServe(":8080", http.DefaultServeMux)
}

```

Unless otherwise noted, the source files are distributed under the BSD-style license found in the LICENSE file.
