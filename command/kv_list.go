package command

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

var _ cli.Command = (*KVListCommand)(nil)
var _ cli.CommandAutocomplete = (*KVListCommand)(nil)

type KVListCommand struct {
	*BaseCommand

	flagDepth      int
	flagFilter     string
	flagRecursive  bool
	flagConcurrent uint
}

func (c *KVListCommand) Synopsis() string {
	return "List data or secrets"
}

func (c *KVListCommand) Help() string {
	helpText := `

Usage: vault kv list [options] PATH

  Lists data from Vault's key-value store at the given path.

  List values under the "my-app" folder of the key-value store:

      $ vault kv list secret/my-app/

  Additional flags and more advanced use cases are detailed below.

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *KVListCommand) Flags() *FlagSets {
	set := c.flagSet(FlagSetHTTP | FlagSetOutputFormat)

	f := set.NewFlagSet("Command Options")

	f.BoolVar(&BoolVar{
		Name:    "recursive",
		Target:  &c.flagRecursive,
		Default: false,
		Usage:   "Recursively list data for a given path.",
	})

	f.IntVar(&IntVar{
		Name:    "depth",
		Target:  &c.flagDepth,
		Default: -1,
		Usage:   "Specifies the depth for recursive listing.",
	})

	f.StringVar(&StringVar{
		Name:    "filter",
		Target:  &c.flagFilter,
		Default: `.*`,
		Usage:   "Specifies a regular expression for filtering paths.",
	})

	f.UintVar(&UintVar{
		Name:    "concurrent",
		Target:  &c.flagConcurrent,
		Default: 16,
		Usage:   "Specifies the number of concurrent recursions to run.",
	})

	return set
}

func (c *KVListCommand) AutocompleteArgs() complete.Predictor {
	return c.PredictVaultFolders()
}

func (c *KVListCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *KVListCommand) Run(args []string) int {
	f := c.Flags()

	if err := f.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	args = f.Args()
	switch {
	case len(args) < 1:
		c.UI.Error(fmt.Sprintf("Not enough arguments (expected 1, got %d)", len(args)))
		return 1
	case len(args) > 1:
		c.UI.Error(fmt.Sprintf("Too many arguments (expected 1, got %d)", len(args)))
		return 1
	}

	if c.flagRecursive && c.flagDepth < -1 {
		c.UI.Error(fmt.Sprintf("Invalid recursion depth: %d", c.flagDepth))
		return 1
	}

	if _, e := regexp.Compile(c.flagFilter); c.flagRecursive && e != nil {
		c.UI.Error(fmt.Sprintf(
			"Invalid regular expression: %s", c.flagFilter,
		))
		return 1
	}

	if c.flagRecursive && c.flagConcurrent <= 0 {
		c.UI.Error(fmt.Sprintf(
			"Invalid concurrency value: %d", c.flagConcurrent,
		))
		return 1
	}

	client, err := c.Client()
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}

	path := ensureTrailingSlash(sanitizePath(args[0]))
	mountPath, v2, err := isKVv2(path, client)
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}

	if v2 {
		path = addPrefixToVKVPath(path, mountPath, "metadata")
		if err != nil {
			c.UI.Error(err.Error())
			return 2
		}
	}

	secret, err := client.Logical().List(path)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error listing %s: %s", path, err))
		return 2
	}
	if secret == nil || secret.Data == nil {
		c.UI.Error(fmt.Sprintf("No value found at %s", path))
		return 2
	}

	// If the secret is wrapped, return the wrapped response.
	if secret.WrapInfo != nil && secret.WrapInfo.TTL != 0 {
		return OutputSecret(c.UI, secret)
	}

	s, ok := extractListData(secret)
	if !ok {
		c.UI.Error(fmt.Sprintf("No entries found at %s", path))
		return 2
	}

	// If we have to list keys recursively.
	if c.flagRecursive {
		if c.flagDepth != 0 {
			var (
				i = int32(1)
				r = &kvListRecursiveParams{
					v2:     &v2,
					client: client,
					data:   []*kvData{},
					depth:  int32(c.flagDepth),
					track:  int32(0),
					filter: regexp.MustCompile(c.flagFilter),
					wg:     sync.WaitGroup{},
					sem:    make(chan int32, c.flagConcurrent),
					mux:    sync.Mutex{},
					tck:    time.NewTicker(time.Millisecond * 100),
				}
				e      bool
				opErr  string
				opList = []string{}
			)

			// One of the base cases for `kvListRecursive()'.
			path = ensureTrailingSlash(path)

			// Append the first entry (only for tabular format).
			if Format(c.UI) == "table" {
				if v2 {
					r.data = append(
						r.data,
						&kvData{
							removePrefixFromVKVPath(path, "metadata/"),
							secret,
							nil,
						},
					)
				} else {
					r.data = append(r.data, &kvData{path, secret, nil})
				}
			}

			// Launch the recursive call and wait for it them to finish.
			r.wg.Add(1)
			go kvListRecursive(r, path, s)
			for len(r.sem) > 0 || atomic.LoadInt32(&r.track) < i {
				// For loop termination.
				if atomic.LoadInt32(&r.track) == 0 {
					atomic.AddInt32(&r.track, 1)
				}

				select {
				case x, ok := <-r.sem:
					if ok {
						// For continuing the loop.
						i += x
					} else {
						break
					}
				default:
					<-r.tck.C
				}
			}
			r.wg.Wait()

			// Print the entries.
			for _, d := range r.data {
				if d.path != "" {
					opList = append(opList, d.path)
					if d.err != nil {
						e = true
						opErr += fmt.Sprintf("\n\t%s: %s\n", d.path, d.err)
					}
				}
			}

			sort.Strings(opList)
			OutputList(c.UI, opList)

			if e {
				c.UI.Error(fmt.Sprintf("Errors:%s", opErr))
				return 3
			}
			return 0
		}
		return 0
	}

	return OutputList(c.UI, secret)
}
