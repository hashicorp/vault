# numcpus

[![Go Reference](https://pkg.go.dev/badge/github.com/tklauser/numcpus.svg)](https://pkg.go.dev/github.com/tklauser/numcpus)
[![GitHub Action Status](https://github.com/tklauser/numcpus/workflows/Tests/badge.svg)](https://github.com/tklauser/numcpus/actions?query=workflow%3ATests)
[![Go Report Card](https://goreportcard.com/badge/github.com/tklauser/numcpus)](https://goreportcard.com/report/github.com/tklauser/numcpus)

Package numcpus provides information about the number of CPU.

It gets the number of CPUs (online, offline, present, possible or kernel
maximum) on a Linux, Darwin, FreeBSD, NetBSD, OpenBSD,  DragonflyBSD or
Solaris/Illumos system.

On Linux, the information is retrieved by reading the corresponding CPU
topology files in `/sys/devices/system/cpu`.

Not all functions are supported on Darwin, FreeBSD, NetBSD, OpenBSD,
DragonflyBSD and Solaris/Illumos.

## Usage

```Go
package main

import (
	"fmt"
	"os"

	"github.com/tklauser/numcpus"
)

func main() {
	online, err := numcpus.GetOnline()
	if err != nil {
		fmt.Fprintf(os.Stderr, "GetOnline: %v\n", err)
	}
	fmt.Printf("online CPUs: %v\n", online)

	possible, err := numcpus.GetPossible()
	if err != nil {
		fmt.Fprintf(os.Stderr, "GetPossible: %v\n", err)
	}
	fmt.Printf("possible CPUs: %v\n", possible)
}
```

## References

* [Linux kernel sysfs documenation for CPU attributes](https://www.kernel.org/doc/Documentation/ABI/testing/sysfs-devices-system-cpu)
* [Linux kernel CPU topology documentation](https://www.kernel.org/doc/Documentation/cputopology.txt)
