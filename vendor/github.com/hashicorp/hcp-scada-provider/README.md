# HCP SCADA Provider ![build badge](https://github.com/hashicorp/hcp-scada-provider/actions/workflows/test.yml/badge.svg?branch=main)

SCADA is an internal component of the infrastructure of the [Hashicorp Cloud Platform](https://cloud.hashicorp.com/) control plane and it stands for Supervisory Control And Data Acquisition. It provides HCP with access to functions and data on HCP-managed clusters.

The provider package establishes and maintains a long-lived connection to allow incoming requests to be served by components running in relative network isolation.

SCADA is a variation on a [NAT traversal](https://en.wikipedia.org/wiki/NAT_traversal) technique.

## Who uses it

It's in use internally at Hashicorp.

## License

[Mozilla Public License v2.0](LICENSE)