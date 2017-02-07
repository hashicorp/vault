package main // import "layeh.com/radius/cmd/radserver"

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"unicode"

	"layeh.com/radius"
)

var secret = flag.String("secret", "", "shared RADIUS secret between clients and server")
var command string
var arguments []string

func handler(w radius.ResponseWriter, p *radius.Packet) {
	username, password, ok := p.PAP()
	if !ok {
		w.AccessReject()
		return
	}
	log.Printf("%s requesting access (%s #%d)\n", username, w.RemoteAddr(), p.Identifier)

	cmd := exec.Command(command, arguments...)

	cmd.Env = os.Environ()
	for _, attr := range p.Attributes {
		name, ok := p.Dictionary.Name(attr.Type)
		if !ok {
			continue
		}
		name = strings.Map(func(r rune) rune {
			if unicode.IsDigit(r) {
				return r
			}
			if unicode.IsLetter(r) {
				if unicode.IsUpper(r) {
					return r
				}
				return unicode.ToUpper(r)
			}
			return '_'
		}, name)
		value := fmt.Sprint(attr.Value)
		cmd.Env = append(cmd.Env, "RADIUS_"+name+"="+value)
	}

	cmd.Env = append(cmd.Env, "RADIUS_USERNAME="+username, "RADIUS_PASSWORD="+password)

	output, err := cmd.Output()
	if err != nil {
		log.Printf("handler error: %s\n", err)
	}

	var attributes []*radius.Attribute
	if len(output) > 0 {
		attributes = []*radius.Attribute{
			p.Dictionary.MustAttr("Reply-Message", string(output)),
		}
	}

	if cmd.ProcessState != nil && cmd.ProcessState.Success() {
		log.Printf("%s accepted (%s #%d)\n", username, w.RemoteAddr(), p.Identifier)
		w.AccessAccept(attributes...)
	} else {
		log.Printf("%s rejected (%s #%d)\n", username, w.RemoteAddr(), p.Identifier)
		w.AccessReject(attributes...)
	}
}

const usage = `
program is executed when an Access-Request RADIUS packet is received. If
program exits sucessfully, an Access-Accept response is sent, otherwise, an
Access-Reject is sent. If standard out is non-empty, it is included as an
Reply-Message attribute in the response.

Any known RADIUS attribute will be added to the process's environment. The
attribute name undergoes the following conversion before being set:
 - it is prefixed with RADIUS_
 - all letters of the attribute name are changed to uppercase
 - any non-digit and non-letter character is replaced with underscore (_)
For example, the NAS-IP-Address attribute will be named RADIUS_NAS_IP_ADDRESS.

Two special environment variables are also include: RADIUS_USERNAME and
RADIUS_PASSWORD, which hold the username and password, respectively.
`

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %s [flags] <program> [program arguments...]\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Fprint(os.Stderr, usage)
	}
	flag.Parse()

	if *secret == "" || flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
	}

	command = flag.Arg(0)
	arguments = flag.Args()[1:]

	log.Println("radserver starting")

	server := radius.Server{
		Handler:    radius.HandlerFunc(handler),
		Secret:     []byte(*secret),
		Dictionary: radius.Builtin,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
