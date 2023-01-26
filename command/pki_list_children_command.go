package command

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/ryanuber/columnize"
)

type PKIListChildrenCommand struct {
	*BaseCommand

	flagConfig          string
	flagReturnIndicator string
	flagDefaultDisabled bool
	flagList            bool

	flagUseNames bool

	flagSignatureMatch    bool
	flagIndirectSignMatch bool
	flagKeyIdMatch        bool
	flagSubjectMatch      bool
	flagPathMatch         bool
}

func (c *PKIListChildrenCommand) Synopsis() string {
	return "Determine Which (of a List) of Certificates Were Issued by A Given Parent Certificate"
}

func (c *PKIListChildrenCommand) Help() string {
	helpText := `
Usage: vault pki list-intermediates PARENT CHILD-LIST
PARENT is the certificate that might be the issuer that everything should be verified against.
CHILD-LIST is a comma-separated list of intermediates to be compared to the PARENT.  It can be omitted, in which case a
list will be constructed from all accessible pki mounts.
This returns a list of issuing certificates, and whether they are a match. 
By default, the type of match required is whether the PARENT has the expected subject, key_id, and could have (directly)
signed this issuer.  This can be updated by changed the corresponding flag.
`
	return strings.TrimSpace(helpText)
}

func (c *PKIListChildrenCommand) Flags() *FlagSets {
	set := c.flagSet(FlagSetHTTP | FlagSetOutputFormat)
	f := set.NewFlagSet("Command Options")

	f.BoolVar(&BoolVar{
		Name:    "subject_match",
		Target:  &c.flagSubjectMatch,
		Default: true,
		EnvVar:  "",
		Usage:   `Whether the subject name of the potential parent cert matches the issuer name of the child cert`,
	})

	f.BoolVar(&BoolVar{
		Name:    "key_id_match",
		Target:  &c.flagKeyIdMatch,
		Default: true,
		EnvVar:  "",
		Usage:   `Whether the subject key_id of the potential parent cert matches the issuing key id of the child cert`,
	})

	f.BoolVar(&BoolVar{
		Name:    "path_match",
		Target:  &c.flagPathMatch,
		Default: false,
		EnvVar:  "",
		Usage:   `Whether the potential parent appears in the certificate chain of the issued cert`,
	})

	f.BoolVar(&BoolVar{
		Name:    "direct_sign",
		Target:  &c.flagSignatureMatch,
		Default: true,
		EnvVar:  "",
		Usage:   `Whether the key of the potential parent signed this issued certificate`,
	})

	f.BoolVar(&BoolVar{
		Name:    "indirect_sign",
		Target:  &c.flagIndirectSignMatch,
		Default: true,
		EnvVar:  "",
		Usage:   `Whether trusting the parent certificate is sufficient to trust the child certificate`,
	})

	f.BoolVar(&BoolVar{
		Name:    "use_names",
		Target:  &c.flagUseNames,
		Default: false,
		EnvVar:  "",
		Usage:   `Whether the list of issuers returned is referred to by name when it exists rather than uuid`,
	})

	return set
}

func (c *PKIListChildrenCommand) Run(args []string) int {
	f := c.Flags()
	if err := f.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	args = f.Args()

	if len(args) < 1 {
		c.UI.Error("Not enough arguments (expected potential parent, got nothing)")
		return 1
	} else if len(args) > 2 {
		c.UI.Error(fmt.Sprintf("Too many arguments (expected only potential issuer and issued, got %d arguments)", len(args)))
		for _, arg := range args {
			if strings.HasPrefix(arg, "-") {
				c.UI.Warn(fmt.Sprintf("Options (%v) must be specified before positional arguments (%v)", arg, args[0]))
				break
			}
		}
		return 1
	}

	client, err := c.Client()
	if err != nil {
		c.UI.Error(fmt.Sprintf("Failed to obtain client: %s", err))
		return 1
	}

	issuer := sanitizePath(args[0])
	issued := ""
	if len(args) > 1 {
		issued = sanitizePath(args[1])
	} else {
		mountListRaw, err := client.Logical().Read("/sys/mounts/")
		if err != nil {
			c.UI.Error(fmt.Sprintf("Failed to Read List of Mounts With Potential Issuers: %v", err))
			return 1
		}
		for path, rawValueMap := range mountListRaw.Data {
			valueMap := rawValueMap.(map[string]interface{})
			if valueMap["type"].(string) == "pki" {
				issuerListEndpoint := sanitizePath(path) + "/issuers"
				rawIssuersResp, err := client.Logical().List(issuerListEndpoint)
				if err != nil {
					c.UI.Error(fmt.Sprintf("Failed to Read List of Issuers within Mount %v: %v", path, err))
					return 1
				}
				if rawIssuersResp == nil {
					continue // TODO: Empty Issuers Response this throws an error
				}
				issuersMap := rawIssuersResp.Data["keys"]
				certList := issuersMap.([]interface{})
				for _, certId := range certList {
					identifier := certId.(string)
					if c.flagUseNames {
						issuerReadResp, err := client.Logical().Read(sanitizePath(path) + "/issuer/" + identifier)
						if err != nil {
							c.UI.Warn(fmt.Sprintf("Unable to Fetch Issuer to Recover Name at: %v", sanitizePath(path)+"/issuer/"+identifier))
						}
						if issuerReadResp != nil {
							issuerName := issuerReadResp.Data["issuer_name"].(string)
							if issuerName != "" {
								identifier = issuerName
							}
						}
					}
					if len(issued) == 0 {
						issued = sanitizePath(path) + "/issuer/" + identifier
					} else {
						issued = issued + "," + sanitizePath(path) + "/issuer/" + identifier
					}
				}
			}
		}
	}

	childrenMatches := make(map[string]bool)

	constraintMap := map[string]bool{
		// This comparison isn't strictly correct, despite a standard ordering these are sets
		"subject_match":   c.flagSubjectMatch,
		"path_match":      c.flagPathMatch,
		"trust_match":     c.flagIndirectSignMatch,
		"key_id_match":    c.flagKeyIdMatch,
		"signature_match": c.flagSignatureMatch,
	}

	for _, child := range strings.Split(issued, ",") {
		path := sanitizePath(child)
		if path != "" {
			err, verifyResults := verifySignBetween(client, issuer, path)
			if err != nil {
				c.UI.Error(fmt.Sprintf("Failed to run verification on path %v: %v", path, err))
				return 1
			}
			childrenMatches[path] = checkIfResultsMatchFilters(verifyResults, constraintMap)
		}
	}

	c.outputResults(childrenMatches)

	return 0
}

func checkIfResultsMatchFilters(verifyResults, constraintMap map[string]bool) bool {
	for key, required := range constraintMap {
		if required == true {
			if verifyResults[key] == false {
				return false
			}
		}
	}
	return true
}

func (c *PKIListChildrenCommand) outputResults(results map[string]bool) error {
	switch Format(c.UI) {
	case "", "table":
		return c.outputResultsTable(results)
	case "json":
		return c.outputResultsJSON(results)
	case "yaml":
		return c.outputResultsYAML(results)
	default:
		return fmt.Errorf("unknown output format: %v", Format(c.UI))
	}
}

func (c *PKIListChildrenCommand) outputResultsTable(results map[string]bool) error {
	data := []string{"intermediate" + hopeDelim + "match?"}
	for field, finding := range results {
		row := field + hopeDelim + strconv.FormatBool(finding)
		data = append(data, row)
	}
	c.UI.Output(tableOutput(data, &columnize.Config{
		Delim: hopeDelim,
	}))
	c.UI.Output("\n")

	return nil
}

func (c *PKIListChildrenCommand) outputResultsJSON(results map[string]bool) error {
	bytes, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return err
	}

	c.UI.Output(string(bytes))
	return nil
}

func (c *PKIListChildrenCommand) outputResultsYAML(results map[string]bool) error {
	bytes, err := yaml.Marshal(results)
	if err != nil {
		return err
	}

	c.UI.Output(string(bytes))
	return nil
}
