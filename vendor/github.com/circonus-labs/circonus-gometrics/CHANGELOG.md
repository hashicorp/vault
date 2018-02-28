# v2.1.1

* dep dependencies
* fix two instances of shadowed variables
* fix several documentation typos
* simplify (gofmt -s)
* remove an inefficient use of regexp.MatchString

# v2.1.0

* Add unix socket capability for SubmissionURL `http+unix://...`
* Add `RecordCountForValue` function to histograms

# v2.0.0

* gauges as `interface{}`
   * change: `GeTestGauge(string) (string,error)` ->  `GeTestGauge(string) (interface{},error)`
   * add: `AddGauge(string, interface{})` to add a delta value to an existing gauge
* prom output candidate
* Add `CHANGELOG.md` to repository
