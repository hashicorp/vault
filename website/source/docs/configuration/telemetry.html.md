---
layout: "docs"
page_title: "Telemetry - Configuration"
sidebar_current: "docs-configuration-telemetry"
description: |-
  The telemetry stanza specifies various configurations for Vault to publish
  metrics to upstream systems.
---

# `telemetry` Stanza

The `telemetry` stanza specifies various configurations for Vault to publish
metrics to upstream systems.

```hcl
telemetry {
  statsite_address = "statsite.company.local:8125"
}
```

## `telemetry` Parameters

Due to the number of configurable parameters to the `telemetry` stanza,
parameters on this page are grouped by the telemetry provider.

### Common

The following options are available on all telemetry configurations.

- `disable_hostname` `(bool: false)` - Specifies if gauge values should be
  prefixed with the local hostname.

### `statsite`

These `telemetry` parameters apply to
[statsite](https://github.com/armon/statsite).

- `statsite_address` `(string: "")` - Specifies the address of a statsite server
  to forward metrics data to.

```hcl
telemetry {
  statsite_address = "statsite.company.local:8125"
}
```

### `statsd`

These `telemetry` parameters apply to
[statsd](https://github.com/etsy/statsd).

- `statsd_address` `(string: "")` - Specifies the address of a statsd server to
  forward metrics to.

```hcl
telemetry {
  statsd_address = "statsd.company.local:8125"
}
```

### `circonus`

These `telemetry` parameters apply to [Circonus](http://circonus.com/).

- `circonus_api_token` `(string: "")` - Specifies a valid Circonus API Token
  used to create/manage check. If provided, metric management is enabled.

- `circonus_api_app` `(string: "nomad")` - Specifies a valid app name associated
  with the API token.

- `circonus_api_url` `(string: "https://api.circonus.com/v2")` - Specifies the
  base URL to use for contacting the Circonus API.

- `circonus_submission_interval` `(string: "10s")` - Specifies the interval at
  which metrics are submitted to Circonus.

- `circonus_submission_url` `(string: "")` - Specifies the
  `check.config.submission_url` field, of a Check API object, from a previously
  created HTTPTRAP check.

- `circonus_check_id` `(string: "")` - Specifies the Check ID (**not check
  bundle**) from a previously created HTTPTRAP check. The numeric portion of the
  `check._cid` field in the Check API object.

- `circonus_check_force_metric_activation` `(bool: false)` - Specifies if force
  activation of metrics which already exist and are not currently active. If
  check management is enabled, the default behavior is to add new metrics as
  they are encountered. If the metric already exists in the check, it will
  not be activated. This setting overrides that behavior.

- `circonus_check_instance_id` `(string: "<hostname>:<application>")` - Serves
  to uniquely identify the metrics coming from this *instance*.  It can be used
  to maintain metric continuity with transient or ephemeral instances as they
  move around within an infrastructure. By default, this is set to
  hostname:application name (e.g. "host123:nomad").

- `circonus_check_search_tag` `(string: <service>:<application>)` - Specifies a
  special tag which, when coupled with the instance id, helps to narrow down the
  search results when neither a Submission URL or Check ID is provided. By
  default, this is set to service:app (e.g. "service:nomad").

- `circonus_check_display_name` `(string: "")` - Specifies a name to give a
   check when it is created. This name is displayed in the Circonus UI Checks
   list.

- `circonus_check_tags` `(string: "")` - Comma separated list of additional
  tags to add to a check when it is created.

- `circonus_broker_id` `(string: "")` - Specifies the ID of a specific Circonus
  Broker to use when creating a new check. The numeric portion of `broker._cid`
  field in a Broker API object. If metric management is enabled and neither a
  Submission URL nor Check ID is provided, an attempt will be made to search for
  an existing check using Instance ID and Search Tag. If one is not found, a new
  HTTPTRAP check will be created. By default, this is a random
  Enterprise Broker is selected, or, the default Circonus Public Broker.

- `circonus_broker_select_tag` `(string: "")` - Specifies a special tag which
  will be used to select a Circonus Broker when a Broker ID is not provided. The
  best use of this is to as a hint for which broker should be used based on
  *where* this particular instance is running (e.g. a specific geo location or
  datacenter, dc:sfo).

### `dogstatsd`

These `telemetry` parameters apply to
[DogStatsD](http://docs.datadoghq.com/guides/dogstatsd/).

- `dogstatsd_addr` `(string: "")` - This provides the address of a DogStatsD
  instance. DogStatsD is a protocol-compatible flavor of statsd, with the added
  ability to decorate metrics with tags and event information. If provided,
  Vault will send various telemetry information to that instance for
  aggregation. This can be used to capture runtime information.


- `dogstatsd_tags` `(string array: [])` - This provides a list of global tags
  that will be added to all telemetry packets sent to DogStatsD. It is a list
  of strings, where each string looks like "my_tag_name:my_tag_value".
