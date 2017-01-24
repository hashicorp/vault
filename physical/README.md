## Physical backends

### etcd Backend

The Valut etcd physical backend supports both v2 and v3 APIs. To explicitly specify the API version, add `etcd_api`
field to the backend section in the configuration file.

```
backend "etcd" {
  address = "http://127.0.0.1:2379"
  path = "vault"
  etcd_api = "3"
}
```

The default `etcd_api` version is auto-detected based on the version of the etcd cluster. If the etcd cluster version is 3.1+ and there is no previous data in v2 API, the auto-detected default is v3 API.

### etcd v3 backend

etcd v3 backend is maintained by [etcd team](https://github.com/coreos/etcd/blob/master/MAINTAINERS#L1-L4). It supports all backend features including HA.

etcd version 3.1+ is required to enable v3 backend.

### etcd v2 backend

etcd v2 backend has known issues with HA support.