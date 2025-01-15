Vault SDK libs
=================

This package provides the `sdk` package which contains code useful for
developing Vault plugins.

Although we try not to break functionality, we reserve the right to reorganize
the code at will and may occasionally cause breaks if they are warranted. As
such we expect the tag of this module will stay less than `v1.0.0`.

For any major changes we will try to give advance notice in the CHANGES section
of Vault's CHANGELOG.md.

## Metrics Emission and Compatibility

This module can emit metrics using either `github.com/armon/go-metrics` or `github.com/hashicorp/go-metrics`. Choosing between the libraries is controlled via build tags. 

**Build Tags**
* `armonmetrics` - Using this tag will cause metrics to be routed to `armon/go-metrics`
* `hashicorpmetrics` - Using this tag will cause all metrics to be routed to `hashicorp/go-metrics`

If no build tag is specified, the default behavior is to use `armon/go-metrics`. 

**Deprecating `armon/go-metrics`**

Emitting metrics to `armon/go-metrics` is officially deprecated. Usage of `armon/go-metrics` will remain the default until mid-2025 with opt-in support continuing to the end of 2025.

**Migration**
To migrate an application currently using the older `armon/go-metrics` to instead use `hashicorp/go-metrics` the following should be done.

1. Upgrade libraries using `armon/go-metrics` to consume `hashicorp/go-metrics/compat` instead. This should involve only changing import statements. All repositories in the `hashicorp` namespace
2. Update an applications library dependencies to those that have the compatibility layer configured.
3. Update the application to use `hashicorp/go-metrics` for configuring metrics export instead of `armon/go-metrics`
   * Replace all application imports of `github.com/armon/go-metrics` with `github.com/hashicorp/go-metrics`
   * Instrument your build system to build with the `hashicorpmetrics` tag.
