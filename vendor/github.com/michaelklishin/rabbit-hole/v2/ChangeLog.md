## Changes Between 2.15.0 and 2.16.0 (in development)

No changes yet

## Changes Between 2.14.0 and 2.15.0 (May 18, 2023)

### Further Improve Shovel Support

Contributed by @ChunyiLyu.

GitHub issue: [#266](https://github.com/michaelklishin/rabbit-hole/pull/268)

## Changes Between 2.13.0 and 2.14.0 (May 11, 2023)

### Allow Setting Default Queue Type for Virtual Hosts

GitHub issue: [#261](https://github.com/michaelklishin/rabbit-hole/pull/261)

Contributed by @mqhenning.

### Correctly Set (Shovel) Destination Publish Properties

Correctly pass `dest-publish-properties` as a map.

GitHub issue: [#262](https://github.com/michaelklishin/rabbit-hole/pull/262)

Contributed by @Galvill.


## Changes Between 2.12.0 and 2.13.0 (Feb 23, 2023)

### Avoids a Panic When Destination Type is Not Set

GitHub issue: [#253](https://github.com/michaelklishin/rabbit-hole/pull/253)

### Federation message TTL is No Longer Set When Omitted

GitHub issue: [#233](https://github.com/michaelklishin/rabbit-hole/pull/233)


## Changes Between 2.11.0 and 2.12.0 (Dec 12, 2021)

### Support for Definition Uploads

GitHub issue: [#220](https://github.com/michaelklishin/rabbit-hole/pull/220)

Contributed by @shubhang93.

### Support to User Limits

GitHub issue: [#217](https://github.com/michaelklishin/rabbit-hole/pull/217)

Contributed by @aitorpazos.

### Listing of All Virtual Host Limits

GitHub issue: [#217](https://github.com/michaelklishin/rabbit-hole/pull/217)

Contributed by @aitorpazos.

### Listing of Connection in a Virtual Host

GitHub issue: [#211](https://github.com/michaelklishin/rabbit-hole/pull/211)

Contributed by @needsaholiday.

## Changes Between 2.10.0 and 2.11.0 (Sep 16, 2021)

This release contains **minor breaking public API changes**.

### Avoid returning empty queue and exchange properties

struct values in `ExportedDefinitions`, `QueueInfo`, and `ExchangeInfo` have all been changed to pointers. This is to avoid having the empty struct values returned when exporting definitions and listing queues and exchanges.

Updated [`ExportedDefinitions`](https://github.com/michaelklishin/rabbit-hole/blob/v2.11.0/definitions.go#L6), [`ExchangeInfo`](https://github.com/michaelklishin/rabbit-hole/blob/v2.11.0/exchanges.go#L23) and [`QueueInfo`](https://github.com/michaelklishin/rabbit-hole/blob/v2.11.0/queues.go#L87).

Contributed by @mkuratczyk.

PRs: [#208](https://github.com/michaelklishin/rabbit-hole/pull/208) and [#209](https://github.com/michaelklishin/rabbit-hole/pull/209)

### Support for a special returned value of Queue and Exchange AutoDelete

`QueueInfo` and `ExchangeInfo` now use a special type for `AutoDelete` because returned value from RabbitMQ server could be a boolean or a string value "undefined".

```go
// AutoDelete is a boolean but RabbitMQ may return the string "undefined"
type AutoDelete bool

// ExchangeInfo represents and exchange and its properties.
type ExchangeInfo struct {
	Name       string                 `json:"name"`
	Vhost      string                 `json:"vhost,omitempty"`
	Type       string                 `json:"type"`
	Durable    bool                   `json:"durable"`
	AutoDelete AutoDelete             `json:"auto_delete"`
...
}

// QueueInfo represents a queue, its properties and key metrics.
type QueueInfo struct {
	// Queue name
	Name string `json:"name"`
	// Queue type
	Type string `json:"type"`
	// Virtual host this queue belongs to
	Vhost string `json:"vhost,omitempty"`
	// Is this queue auto-deleted?
	AutoDelete AutoDelete `json:"auto_delete"`
         ...
}
```
Contributed by @mkuratczyk.

PR: [#207](https://github.com/michaelklishin/rabbit-hole/pull/207)

### Support listing definitions for a specific vhost

Contributed by @mkuratczyk.

PR: [#206](https://github.com/michaelklishin/rabbit-hole/pull/206)

### Support for vhost limits

Contributed by @Sauci.

PR: [#200](https://github.com/michaelklishin/rabbit-hole/pull/200)

### Return more fields for `NodeInfo` and `QueueInfo`

Contributed by  @hjweddie.

PR: [#198](https://github.com/michaelklishin/rabbit-hole/pull/198)


## Changes Between 2.9.0 and 2.10.0 (Jun 3, 2021)

This release contains very **minor breaking public API changes**.

### `ShovelDefinition.SourceDeleteAfter` Type Now Matches That of `ShovelDefinition.DeleteAfter`

`ShovelDefinition.SourceDeleteAfter` type has changed to match that of
`ShovelDefinition.DeleteAfter`.

GitHub issue: [#197](https://github.com/michaelklishin/rabbit-hole/pull/197)


## Changes Between 2.8.0 and 2.9.0 (Jun 2, 2021)

This release contains **minor breaking public API changes**.

### Support for Lists of Federation Upstream URIs

Federation definition now uses a dedicated type, `URISet`, to represent
a set of URIs that will be tried sequentially until the link
can successfully connect and authenticate:

``` go
def1 := FederationDefinition{
            Uri: URISet{"amqp://hostname/%2f"},
        }
```

`URISet` has now replaced `ShovelURISet`:

``` go
sDef := ShovelDefinition{
            SourceURI:         URISet([]string{"amqp://127.0.0.1/%2f"}),
            SourceQueue:       "mySourceQueue",
            DestinationURI:    ShovelURISet([]string{"amqp://host1/%2f"}),
            DestinationQueue:  "myDestQueue",
            AddForwardHeaders: true,
            AckMode:           "on-confirm",
            DeleteAfter:       "never",
        }
```

GitHub issues: [#193](https://github.com/michaelklishin/rabbit-hole/pull/193), [#194](https://github.com/michaelklishin/rabbit-hole/pull/194)

### Support for Operator Policies

Contributed by @MrLuje.

GitHub issues: [#188](https://github.com/michaelklishin/rabbit-hole/issues/188), [#190](https://github.com/michaelklishin/rabbit-hole/pull/190)

### Declared Queue Type is Correctly Propagated

GitHub issue: [#189](https://github.com/michaelklishin/rabbit-hole/pull/189)


## Changes Between 2.7.0 and 2.8.0 (Apr 12, 2021)

### Global Runtime Parameters

The library now supports global runtime parameters:

``` go
// list all global parameters
params, err := rmqc.ListGlobalParameters()
// => []GlobalRuntimeParameter, error
// get a global parameter
p, err := rmqc.GetGlobalParameter("name")
// => *GlobalRuntimeParameter, error
// declare or update a global parameter
resp, err := rmqc.PutGlobalParameter("name", map[string]interface{
    endpoints: "amqp://server-name",
})
// => *http.Response, error
// delete a global parameter
resp, err := rmqc.DeleteGlobalParameter("name")
// => *http.Response, error
```

Contributed by @ChunyiLyu.

GitHub issue: [#180](https://github.com/michaelklishin/rabbit-hole/pull/180)

## Changes Between 2.6.0 and 2.7.0 (Mar 30, 2021)

This release contains **minor breaking public API changes**
and targets RabbitMQ 3.8.x (the only [supported version at the time of writing](https://www.rabbitmq.com/versions.html))
exclusively.

### Support for Modern Health Check Endpoints

The client now supports [modern health check endpoints](https://www.rabbitmq.com/monitoring.html#health-checks)
(same checks as provided by `rabbitmq-diagnostics`):

``` go
import (
       "github.com/michaelklishin/rabbit-hole/v2"
)

rmqc, _ = NewClient("http://127.0.0.1:15672", "username", "$3KrEƮ")

res1, err1 := rmqc.HealthCheckAlarms()

res2, err2 := rmqc.HealthCheckLocalAlarms()

res3, err3 := rmqc.HealthCheckCertificateExpiration(1, DAYS)

res4, err4 := rmqc.HealthCheckPortListener(5672)

res5, err5 := rmqc.HealthCheckProtocolListener(AMQP091)

res6, err6 := rmqc.HealthCheckVirtualHosts()

res7, err7 := rmqc.HealthCheckNodeIsMirrorSyncCritical()

res8, err8 := rmqc.HealthCheckNodeIsQuorumCritical()
```

Contributed by Martin @mkrueger-sabio Krueger.

GitHub issue: [#173](https://github.com/michaelklishin/rabbit-hole/pull/173)

### Support for Inspecting Shovel Status

`ListShovelStatus` is a new function that returns a list of
Shovel status reports for a virtual host:

``` go
res, err := rmqc.ListShovelStatus("a-virtual-host")
```

Contributed by Martin @mkrueger-sabio Krueger.

GitHub issue: [#178](https://github.com/michaelklishin/rabbit-hole/pull/178)

### Support for Lists of Shovel URIs

Shovel definition now uses a dedicated type, `ShovelURISet`, to represent
a set of URIs that will be tried sequentially until the Shovel
can successfully connect and authenticate:

``` go
sDef := ShovelDefinition{
            SourceURI:         ShovelURISet([]string{"amqp://host2/%2f", "amqp://host3/%2f"}),
            SourceQueue:       "mySourceQueue",
            DestinationURI:    ShovelURISet([]string{"amqp://host1/%2f"}),
            DestinationQueue:  "myDestQueue",
            AddForwardHeaders: true,
            AckMode:           "on-confirm",
            DeleteAfter:       "never",
        }
```

Source and destination URI sets are only supported by the Shovel plugin in
RabbitMQ 3.8.x.

Originally suggested by @pathcl in #172.

### Definition Export

`rabbithole.ListDefinitions` is a new function that retuns
[exported definitions from a cluster](https://www.rabbitmq.com/definitions.html)
as a typed Go data structure.

Contributed by @pathcl.

GitHub issue: [#170](https://github.com/michaelklishin/rabbit-hole/pull/170)

### User Tags as Array

For forward compatibility with RabbitMQ 3.9, as of this
version the list of user tags is returned as an array
intead of a comma-separated string.

Compatibility with earlier RabbitMQ HTTP API versions, such as 3.8,
has been preserved.

### Optional Federation Parameters are Now Marked with `omitempty`

Contributed by Michał @michalkurzeja Kurzeja.

GitHub issue: [#177](https://github.com/michaelklishin/rabbit-hole/pull/177)


## Changes Between 2.5.0 and 2.6.0 (Nov 25, 2020)

### Feature Flag Management

The client now can list and enable feature flags
using the `ListFeatureFlags` and `EnableFeatureFlag` functions.

Contributed by David Ansari.

GitHub issue: [#167](https://github.com/michaelklishin/rabbit-hole/pull/167)

## Changes Between 2.4.0 and 2.5.0 (Sep 28th, 2020)

### Shovels: Support for Numerical Delete-After Values

The `delete-after` Shovel parameter now can be deserialised to
a numerical TTL value as well as special string values such as `"never"`.

Contributed by Michal @mkuratczyk Kuratczyk.

GitHub issue: [#164](https://github.com/michaelklishin/rabbit-hole/pull/164)


## Changes Between 2.3.0 and 2.4.0 (Aug 4th, 2020)

### More Thorough Error Checking of HTTP[S] Requests

Suggested by @mammothbane.

GitHub issue: [#158](https://github.com/michaelklishin/rabbit-hole/issues/158)

### Salt Generation Helper Now Uses `crypto/rand` Instead of `math/rand`

Suggested by @mammothbane.

GitHub issue: [#160](https://github.com/michaelklishin/rabbit-hole/issues/160)

## More Standardized Response Errors

Error responses (`40x` with the exception of `404` in response to `DELETE` operations,
`50x`) HTTP API response errors are now always wrapped into`ErrorResponse`,
even if they do not carry a JSON payload.


## Changes Between 2.2.0 and 2.3.0 (July 11th, 2020)

### New Endpoints for Listing Federation Links

Contributed by @niclic.

GitHub issue: [#155](https://github.com/michaelklishin/rabbit-hole/pull/155)

### Support for More Shovel Parameters (e.g. for AMQP 1.0 Sources and Destinations)

Contributed by @akurz.

GitHub issue: [#155](https://github.com/michaelklishin/rabbit-hole/pull/157)

### Conditional Exclusion of Expiration Field

Contributed by @niclic.

GitHub issue: [#154](https://github.com/michaelklishin/rabbit-hole/pull/154)


## Changes Between 2.1.0 and 2.2.0 (May 21st, 2020)

### [Runtime Parameter](https://www.rabbitmq.com/parameters.html) and [Federation Upstream](https://www.rabbitmq.com/federation.html) Management

Contributed by @niclic.

GitHub issue: [#150](https://github.com/michaelklishin/rabbit-hole/pull/150)

### Improved Error Reporting

Contributed by @niclic.

GitHub issue: [michaelklishin/rabbit-hole#152](https://github.com/michaelklishin/rabbit-hole/pull/152)

### Fixed a null Pointer in HTTP Response Handling

Contributed by @justabaka.

GitHub issue: [#148](https://github.com/michaelklishin/rabbit-hole/pull/148)


## Changes Between 2.0.0 and 2.1.0 (Feb 1st, 2020)

### Corrects Package Version

See [Semantic Go Module Import Versioning](https://github.com/golang/go/wiki/Modules#semantic-import-versioning) for details

GitHub issue: [#146](https://github.com/michaelklishin/rabbit-hole/issues/146)

### New Endpoint, `DELETE /topic-permissions/{vhost}/{user}/{exchange}`

Contributed by Barnaby Shearer.

GitHub issues: [#147](https://github.com/michaelklishin/rabbit-hole/pull/147)

### Exposed Client Connection Time Field

Available in RabbitMQ 3.7 and later versions.

Contributed by @kgrieco.

GitHub issue: [#144](https://github.com/michaelklishin/rabbit-hole/pull/144)

### Authentication Failures Now Return a Reasonable Error

Contributed by @mazamats.

GitHub issues: [#145](https://github.com/michaelklishin/rabbit-hole/pull/145), [#112](https://github.com/michaelklishin/rabbit-hole/issues/112)


## Changes Between 1.5.0 and 2.0.0 (October 8th, 2019)

### Go 1.9 through 1.11 Support Dropped

This library now only supports Go 1.12 and 1.13 (two most recent minor GA releases).

### Unroutable Message Metric Support

The `drop_unroutable` metric is specific to RabbitMQ 3.8.

Contributed by David Ansari and Feroz Jilla.

### Support for Exchange Ingress and Egress Rates

Contributed by Rajendra N Acharya.

### Eager Synchronization of Classic Queue

It is now possible to initiate an eager sync of a classic mirrored queue and cancel it.

Contributed by Jaroslaw Bochniak.

GitHub issue: [#143](https://github.com/michaelklishin/rabbit-hole/pull/143)

### Queue Status JSON Serialization Fixed

Contributed by Andrew Wang.

### GET /api/consumers Support

Contributed by Thomas Hudry.

GitHub issue: [#140](https://github.com/michaelklishin/rabbit-hole/pull/140)

### http.Transport Replaced by http.RoundTripper

HTTP client configuration now uses `http.RoundTripper`.

GitHub issue: [#123](https://github.com/michaelklishin/rabbit-hole/pull/123).

Contributed by Radek Simko.

### Go Modules Support

GitHub issues: [#124](https://github.com/michaelklishin/rabbit-hole/pull/124), [#128](https://github.com/michaelklishin/rabbit-hole/pull/128).

Contributed by Radek Simko and Gerhard Lazu.


## Changes Between 1.4.0 and 1.5.0 (February 13th, 2019)

### More Binding Management Functions

`ListExchangeBindings`, `ListExchangeBindingsWithSource`, `ListExchangeBindingsWithDestination`,
and `ListExchangeBindingsBetween` are new functions that list bindings,
in particular between exchanges.

GitHub issue: [#109](https://github.com/michaelklishin/rabbit-hole/pull/109).

### Password Hash Generation Helpers

It is now possible to specify a `password_hash` when creating a user.
Helper functions `GenerateSalt` and `SaltedPasswordHashSHA256` make this more
straightforward compared to implementing [the algorithm](http://www.rabbitmq.com/passwords.html#computing-password-hash)
directly.

GitHub issue: [#119](https://github.com/michaelklishin/rabbit-hole/pull/119)

### Paginated Queue Listing

A new function, `PagedListQueuesWithParameters`, can list queues with pagination support.

GitHub issue: [#118](https://github.com/michaelklishin/rabbit-hole/pull/118)

### More `NodeInfo` and `QueueInfo` Attributes

GitHub issue: [#115](https://github.com/michaelklishin/rabbit-hole/issues/115)

### URL.Opaque Left to Its Own Devices

The client no longer messes with `URL.Opaque` as it doesn't seem to
be necessary any more for correct %-encoding of URL path.

GitHub issue: [#121](https://github.com/michaelklishin/rabbit-hole/issues/121)


## Changes Between 1.0.0 and 1.1.0 (Dec 1st, 2015)

### More Complete Message Stats Information

Message stats now include fields such as `deliver_get` and `redeliver`.

GH issue: [#73](https://github.com/michaelklishin/rabbit-hole/pull/73).

Contributed by Edward Wilde.


## 1.0 (first tagged release, Dec 25th, 2015)

### TLS Support

`rabbithole.NewTLSClient` is a new function which works
much like `NewClient` but additionally accepts a transport.

Contributed by @[GrimTheReaper](https://github.com/GrimTheReaper).

### Federation Support

It is now possible to create federation links
over HTTP API.

Contributed by [Ryan Grenz](https://github.com/grenzr-bskyb).

### Core Operations Support

Most common HTTP API operations (listing and management of
vhosts, users, permissions, queues, exchanges, and bindings)
are supported by the client.
