package main // import "layeh.com/radius/cmd/radtest"

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"

	"layeh.com/radius"
)

const usage = `
Sends an Access-Request RADIUS packet to a server and prints the result.
`

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %s [flags] <user> <password> <radius-server>[:port] <nas-port-number> <secret>\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Fprint(os.Stderr, usage)
	}
	timeout := flag.Duration("timeout", time.Second*10, "timeout for the request to finish")
	flag.Parse()
	if flag.NArg() != 5 {
		flag.Usage()
		os.Exit(1)
	}

	host, port, err := net.SplitHostPort(flag.Arg(2))
	if err != nil {
		host = flag.Arg(2)
		port = "1812"
	}
	hostport := net.JoinHostPort(host, port)

	packet := radius.New(radius.CodeAccessRequest, []byte(flag.Arg(4)))
	packet.Add("User-Name", flag.Arg(0))
	packet.Add("User-Password", flag.Arg(1))
	nasPort, _ := strconv.Atoi(flag.Arg(3))
	packet.Add("NAS-Port", uint32(nasPort))

	client := radius.Client{
		DialTimeout: *timeout,
		ReadTimeout: *timeout,
	}
	received, err := client.Exchange(packet, hostport)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var status string
	if received.Code == radius.CodeAccessAccept {
		status = "Accept"
	} else {
		status = "Reject"
	}
	if msg, ok := received.Value("Reply-Message").(string); ok {
		status += " (" + msg + ")"
	}

	fmt.Println(status)

	if received.Code != radius.CodeAccessAccept {
		os.Exit(2)
	}
}
