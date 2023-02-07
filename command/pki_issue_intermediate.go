// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package command

import (
	"context"
	"fmt"
	"io"
	"os"
	paths "path"
	"strings"

	"github.com/hashicorp/vault/api"
	"github.com/posener/complete"
)

type PKIIssueCACommand struct {
	*BaseCommand

	flagConfig          string
	flagReturnIndicator string
	flagDefaultDisabled bool
	flagList            bool

	flagKeyStorageSource string
	flagNewIssuerName    string
}

func (c *PKIIssueCACommand) Synopsis() string {
	return "Given a Parent Certificate, and a List of Generation Parameters, Creates an Issue on a Specified Mount"
}

func (c *PKIIssueCACommand) Help() string {
	helpText := `
Usage: vault pki issue PARENT CHILD_MOUNT options

PARENT is the fully qualified path of the Certificate Authority in vault which will issue the new intermediate certificate.

CHILD_MOUNT is the path of the mount in vault where the new issuer is saved.

options are the superset of the options passed to generate/intermediate and sign-intermediate commands.  At least one option must be set.

This command creates a intermediate certificate authority certificate signed by the parent in the CHILD_MOUNT.

` + c.Flags().Help()
	return strings.TrimSpace(helpText)
}

func (c *PKIIssueCACommand) Flags() *FlagSets {
	set := c.flagSet(FlagSetHTTP | FlagSetOutputFormat)
	f := set.NewFlagSet("Command Options")

	f.StringVar(&StringVar{
		Name:       "type",
		Target:     &c.flagKeyStorageSource,
		Default:    "internal",
		EnvVar:     "",
		Usage:      `Options are "existing" - to use an existing key inside vault, "internal" - to generate a new key inside vault, or "kms" - to link to an external key.  Exported keys are not available through this API.`,
		Completion: complete.PredictSet("internal", "existing", "kms"),
	})

	f.StringVar(&StringVar{
		Name:    "issuer_name",
		Target:  &c.flagNewIssuerName,
		Default: "",
		EnvVar:  "",
		Usage:   `If present, the newly created issuer will be given this name.`,
	})

	return set
}

func (c *PKIIssueCACommand) Run(args []string) int {
	// Parse Args
	f := c.Flags()
	if err := f.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}
	args = f.Args()

	if len(args) < 3 {
		c.UI.Error("Not enough arguments expected parent issuer and child-mount location and some key_value argument")
		return 1
	}

	stdin := (io.Reader)(os.Stdin)
	data, err := parseArgsData(stdin, args[2:])
	if err != nil {
		c.UI.Error(fmt.Sprintf("Failed to parse K=V data: %s", err))
		return 1
	}

	// Check We Have a Client
	client, err := c.Client()
	if err != nil {
		c.UI.Error(fmt.Sprintf("Failed to obtain client: %v", err))
		return 1
	}

	// Sanity Check the Parent Issuer
	parentMountIssuer := sanitizePath(args[0]) // /pki/issuer/default
	_, parentIssuerName := paths.Split(parentMountIssuer)
	if !strings.Contains(parentMountIssuer, "/issuer/") {
		c.UI.Error(fmt.Sprintf("Parent Issuer %v is Not a PKI Issuer Path of the format /mount/issuer/issuer-ref", parentMountIssuer))
	}
	_, err = client.Logical().Read(parentMountIssuer + "/json")
	if err != nil {
		c.UI.Error(fmt.Sprintf("Unable to access parent issuer %v: %v", parentMountIssuer, err))
	}

	// Set-up Failure State (Immediately Before First Write Call)
	intermediateMount := sanitizePath(args[1])
	failureState := inCaseOfFailure{
		intermediateMount: intermediateMount,
		parentMount:       strings.Split(parentMountIssuer, "/issuer/")[0],
		parentIssuer:      parentMountIssuer,
		newName:           c.flagNewIssuerName,
	}

	// Generate Certificate Signing Request
	csrResp, err := client.Logical().Write(intermediateMount+"/intermediate/generate/"+c.flagKeyStorageSource, data)
	if err != nil {
		if strings.Contains(err.Error(), "no handler for route") { // Mount Given Does Not Exist
			c.UI.Error(fmt.Sprintf("Given Intermediate Mount %v Does Not Exist: %v", intermediateMount, err))
		} else if strings.Contains(err.Error(), "unsupported path") { // Expected if Not a PKI Mount
			c.UI.Error(fmt.Sprintf("Given Intermeidate Mount %v Is Not a PKI Mount: %v", intermediateMount, err))
		} else {
			c.UI.Error(fmt.Sprintf("Failled to Generate Intermediate CSR on %v: %v", intermediateMount, err))
		}
		return 1
	}
	// Parse CSR Response, Also Verifies that this is a PKI Mount
	// (eg. calling the above call on cubbyhole/ won't return an error response)
	csrPemRaw, present := csrResp.Data["csr"]
	if !present {
		c.UI.Error(fmt.Sprintf("Failed to Generate Intermediate CSR on %v, got response: %v", intermediateMount, csrResp))
		return 1
	}
	keyIdRaw, present := csrResp.Data["key_id"]
	if !present && c.flagKeyStorageSource == "internal" {
		c.UI.Error(fmt.Sprintf("Failed to Generate Key on %v, got response: %v", intermediateMount, csrResp))
		return 1
	}

	// If that all Parses, then we've successfully generated a CSR!  Save It (and the Key-ID)
	failureState.csrGenerated = true
	if c.flagKeyStorageSource == "internal" {
		failureState.createdKeyId = keyIdRaw.(string)
	}
	csr := csrPemRaw.(string)
	failureState.csr = csr
	data["csr"] = csr

	// Next, Sign the CSR
	rootResp, err := client.Logical().Write(parentMountIssuer+"/sign-intermediate", data)
	if err != nil {
		c.UI.Error(failureState.generateFailureMessage())
		c.UI.Error(fmt.Sprintf("Error Signing Intermiate On %v", err))
		return 1
	}
	// Success!  Save Our Progress (and Parse the Response)
	failureState.csrSigned = true
	serialNumber := rootResp.Data["serial_number"].(string)
	failureState.certSerialNumber = serialNumber

	caChain := rootResp.Data["ca_chain"].([]interface{})
	caChainPemBundle := ""
	for _, cert := range caChain {
		caChainPemBundle += cert.(string) + "\n"
	}
	failureState.caChain = caChainPemBundle

	// Next Import Certificate
	certificate := rootResp.Data["certificate"].(string)
	issuerId, err := importIssuerWithName(client, intermediateMount, certificate, c.flagNewIssuerName)
	failureState.certIssuerId = issuerId
	if err != nil {
		if strings.Contains(err.Error(), "error naming issuer") {
			failureState.certImported = true
			c.UI.Error(failureState.generateFailureMessage())
			c.UI.Error(fmt.Sprintf("Error Naming Newly Imported Issuer: %v", err))
			return 1
		} else {
			c.UI.Error(failureState.generateFailureMessage())
			c.UI.Error(fmt.Sprintf("Error Importing Into %v Newly Created Issuer %v: %v", intermediateMount, certificate, err))
			return 1
		}
	}
	failureState.certImported = true

	// Then Import Issuing Certificate
	issuingCa := rootResp.Data["issuing_ca"].(string)
	_, err = importIssuerWithName(client, intermediateMount, issuingCa, parentIssuerName)
	if err != nil {
		if strings.Contains(err.Error(), "error naming issuer") {
			c.UI.Warn(fmt.Sprintf("Unable to Set Name on Parent Cert from %v Imported Into %v with serial %v, err: %v", parentIssuerName, intermediateMount, serialNumber, err))
		} else {
			c.UI.Error(failureState.generateFailureMessage())
			c.UI.Error(fmt.Sprintf("Error Importing Into %v Newly Created Issuer %v: %v", intermediateMount, certificate, err))
			return 1
		}
	}

	// Finally Import CA_Chain (just in case there's more information)
	if len(caChain) > 2 { // We've already imported parent cert and newly issued cert above
		importData := map[string]interface{}{
			"pem_bundle": caChainPemBundle,
		}
		_, err := client.Logical().Write(intermediateMount+"/issuers/import/cert", importData)
		if err != nil {
			c.UI.Error(failureState.generateFailureMessage())
			c.UI.Error(fmt.Sprintf("Error Importing CaChain into %v: %v", intermediateMount, err))
			return 1
		}
	}
	failureState.caChainImported = true

	// Finally we read our newly issued certificate in order to tell our caller about it
	c.readAndOutputNewCertificate(client, intermediateMount, issuerId)

	return 0
}

func (c *PKIIssueCACommand) readAndOutputNewCertificate(client *api.Client, intermediateMount string, issuerId string) {
	resp, err := client.Logical().Read(sanitizePath(intermediateMount + "/issuer/" + issuerId))
	if err != nil || resp == nil {
		c.UI.Error(fmt.Sprintf("Error Reading Fully Imported Certificate from %v : %v",
			intermediateMount+"/issuer/"+issuerId, err))
	}

	OutputSecret(c.UI, resp)
}

func importIssuerWithName(client *api.Client, mount string, bundle string, name string) (issuerUUID string, err error) {
	importData := map[string]interface{}{
		"pem_bundle": bundle,
	}
	writeResp, err := client.Logical().Write(mount+"/issuers/import/cert", importData)
	if err != nil {
		return "", err
	}
	mapping := writeResp.Data["mapping"].(map[string]interface{})
	if len(mapping) > 1 {
		return "", fmt.Errorf("multiple issuers returned, while expected one, got %v", writeResp)
	}
	for issuerId := range mapping {
		issuerUUID = issuerId
	}
	if name != "" && name != "default" {
		nameReq := map[string]interface{}{
			"issuer_name": name,
		}
		ctx := context.Background()
		_, err = client.Logical().JSONMergePatch(ctx, mount+"/issuer/"+issuerUUID, nameReq)
		if err != nil {
			return issuerUUID, fmt.Errorf("error naming issuer %v to %v: %v", issuerUUID, name, err)
		}
	}
	return issuerUUID, nil
}

type inCaseOfFailure struct {
	csrGenerated    bool
	csrSigned       bool
	certImported    bool
	certNamed       bool
	caChainImported bool

	intermediateMount string
	createdKeyId      string
	csr               string
	caChain           string
	parentMount       string
	parentIssuer      string
	certSerialNumber  string
	certIssuerId      string
	newName           string
}

func (state inCaseOfFailure) generateFailureMessage() string {
	message := "A failure has occurred"

	if state.csrGenerated {
		message += fmt.Sprintf(" after \n a Certificate Signing Request was successfully generated on mount %v", state.intermediateMount)
	}
	if state.csrSigned {
		message += fmt.Sprintf(" and after \n that Certificate Signing Request was successfully signed by mount %v", state.parentMount)
	}
	if state.certImported {
		message += fmt.Sprintf(" and after \n the signed certificate was reimported into mount %v , with issuerID %v", state.intermediateMount, state.certIssuerId)
	}

	if state.csrGenerated {
		message += "\n\nTO CONTINUE: \n" + state.toContinue()
	}
	if state.csrGenerated && !state.certImported {
		message += "\n\nTO ABORT: \n" + state.toAbort()
	}

	message += "\n"

	return message
}

func (state inCaseOfFailure) toContinue() string {
	message := ""
	if !state.csrSigned {
		message += fmt.Sprintf("You can continue to work with this Certificate Signing Request CSR PEM, by saving"+
			" it as `pki_int.csr`: %v \n Then call `vault write %v/sign-intermediate csr=@pki_int.csr ...` adding the "+
			"same key-value arguements as to `pki issue` (except key_type and issuer_name) to generate the certificate "+
			"and ca_chain", state.csr, state.parentIssuer)
	}
	if !state.certImported {
		if state.caChain != "" {
			message += fmt.Sprintf("The certificate chain, signed by %v, for this new certificate is: %v", state.parentIssuer, state.caChain)
		}
		message += fmt.Sprintf("You can continue to work with this Certificate (and chain) by saving it as "+
			"chain.pem and importing it as `vault write %v/issuers/import/cert pem_bundle=@chain.pem`",
			state.intermediateMount)
	}
	if !state.certNamed {
		issuerId := state.certIssuerId
		if issuerId == "" {
			message += fmt.Sprintf("The issuer_id is returned as the key in a key_value map from importing the " +
				"certificate chain.")
			issuerId = "<issuer-uuid>"
		}
		message += fmt.Sprintf("You can name the newly imported issuer by calling `vault patch %v/issuer/%v "+
			"issuer_name=%v`", state.intermediateMount, issuerId, state.newName)
	}
	return message
}

func (state inCaseOfFailure) toAbort() string {
	if !state.csrGenerated || (!state.csrSigned && state.createdKeyId == "") {
		return "No state was created by running this command.  Try rerunning this command after resolving the error."
	}
	message := ""
	if state.csrGenerated && state.createdKeyId != "" {
		message += fmt.Sprintf(" A key, with key ID %v was created on mount %v as part of this command."+
			"  If you do not with to use this key and corresponding CSR/cert, you can delete that information by calling"+
			" `vault delete %v/key/%v`", state.createdKeyId, state.intermediateMount, state.intermediateMount, state.createdKeyId)
	}
	if state.csrSigned {
		message += fmt.Sprintf("A certificate with serial number %v was signed by mount %v as part of this command."+
			" If you do not want to use this certificate, consider revoking it by calling `vault write %v/revoke/%v`",
			state.certSerialNumber, state.parentMount, state.parentMount, state.certSerialNumber)
	}
	//if state.certImported {
	//	message += fmt.Sprintf("An issuer with UUID %v was created on mount %v as part of this command.  " +
	//		"If you do not wish to use this issuer, consider deleting it by calling `vault delete %v/issuer/%v`",
	//		state.certIssuerId, state.intermediateMount, state.intermediateMount, state.certIssuerId)
	//}

	return message
}
