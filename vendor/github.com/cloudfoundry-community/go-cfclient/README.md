# go-cfclient
[![Travis-CI](https://travis-ci.org/cloudfoundry-community/go-cfclient.svg)](https://travis-ci.org/cloudfoundry-community/go-cfclient)
[![GoDoc](https://godoc.org/github.com/cloudfoundry-community/go-cfclient?status.svg)](http://godoc.org/github.com/cloudfoundry-community/go-cfclient)
[![Report card](https://goreportcard.com/badge/github.com/cloudfoundry-community/go-cfclient)](https://goreportcard.com/report/github.com/cloudfoundry-community/go-cfclient)

### Overview

`cfclient` is a package to assist you in writing apps that need to interact with [Cloud Foundry](http://cloudfoundry.org).
It provides functions and structures to retrieve and update


### Usage

```
go get github.com/cloudfoundry-community/go-cfclient
```

NOTE: Currently this project is not versioning its releases and so breaking changes might be introduced.
Whilst hopefully notifications of breaking changes are made via commit messages, ideally your project will use a local
vendoring system to lock in a version of `go-cfclient` that is known to work for you.
This will allow you to control the timing and maintenance of upgrades to newer versions of this library.

Some example code:

```go
package main

import (
	"github.com/cloudfoundry-community/go-cfclient"
)

func main() {
  c := &cfclient.Config{
    ApiAddress:   "https://api.10.244.0.34.xip.io",
    Username:     "admin",
    Password:     "admin",
  }
  client, _ := cfclient.NewClient(c)
  apps, _ := client.ListApps()
  fmt.Println(apps)
}
```

### Development

#### Errors

If the Cloud Foundry error definitions change at <https://github.com/cloudfoundry/cloud_controller_ng/blob/master/vendor/errors/v2.yml>
then the error predicate functions in this package need to be regenerated.

To do this, simply use Go to regenerate the code:

```
go generate
```

### Contributing

Pull requests welcome.
