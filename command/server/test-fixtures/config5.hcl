disable_cache = true
 disable_mlock = true

 ui = true

 listener "tcp" {
     address = "127.0.0.1:443"
 	 allow_stuff = true
 }

 backend "consul" {
     foo = "bar"
     advertise_addr = "foo"
 }

 ha_backend "consul" {
     bar = "baz"
     advertise_addr = "snafu"
     disable_clustering = "true"
 }

 service_registration "consul" {
     foo = "bar"
 }

 telemetry {
     statsd_address = "bar"
     usage_gauge_period = "5m"
     maximum_gauge_cardinality = 100

     statsite_address = "foo"
     dogstatsd_addr = "127.0.0.1:7254"
     dogstatsd_tags = ["tag_1:val_1", "tag_2:val_2"]
     metrics_prefix = "myprefix"

     lease_metrics_epsilon = "1h"
     num_lease_metrics_buckets = 2
     add_lease_metrics_namespace_labels = true 
 }

 sentinel {
     additional_enabled_modules = []
 }

 max_lease_ttl = "10h"
 default_lease_ttl = "10h"
 cluster_name = "testcluster"
 pid_file = "./pidfile"
 raw_storage_endpoint = true
 disable_sealwrap = true
 disable_printable_check = true