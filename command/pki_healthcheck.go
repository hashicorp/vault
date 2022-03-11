package command

import (
	"bytes"
	"context"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/mitchellh/cli"
	"github.com/posener/complete"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
)

const (
	pkiRetOK            int = 0
	pkiRetUsage         int = 1
	pkiRetInternal      int = 2
	pkiRetInformational int = 3
	pkiRetCritical      int = 4
)

const (
	oneDay   = 24 * time.Hour
	oneMonth = 30 * oneDay
	oneYear  = 365 * oneDay

	suggestNoCRL  = 7 * oneDay
	suggestedCRL  = 50 * oneDay
	definitelyCRL = 90 * oneDay
)

var (
	pkiShouldAuditReqKeys  = []string{"csr", "common_name", "alt_names", "other_sans", "ip_sans", "uri_sans", "ttl", "not_after", "certificate", "ou", "organization", "country", "locality", "province", "street_address", "postal_code", "serial_number", "permitted_dns_domains", "key_type", "key_bits", "policy_identifiers", "not_before_duration", "key_usage", "ext_key_usage", "ext_key_usage_oids", "use_csr_common_name", "use_csr_sans", "use_csr_values", "policy_identifiers"}
	pkiShouldAuditRespKeys = []string{"certificate", "issuing_ca", "serial_number", "error"}
)

// Paths ending with a / apply to subpaths (e.g., /roles/:name).
var pkiPrivilegedPolicyPaths = map[string][]string{
	"/config/ca":  {vault.CreateCapability, vault.UpdateCapability},
	"/config/crl": {vault.CreateCapability, vault.UpdateCapability},
	"/config/url": {vault.CreateCapability, vault.UpdateCapability},
	// Rotating CRL isn't privileged.
	"/intermediate/generate/internal": {vault.CreateCapability, vault.UpdateCapability},
	"/intermediate/generate/external": {vault.CreateCapability, vault.UpdateCapability},
	"/intermediate/set-signed":        {vault.CreateCapability, vault.UpdateCapability},
	// Revoke is privileged since it is by serial number and not by proof
	// of possession.
	"/revoke":                 {vault.CreateCapability, vault.UpdateCapability},
	"/roles/":                 {vault.CreateCapability, vault.UpdateCapability, vault.DeleteCapability},
	"/root/generate/internal": {vault.CreateCapability, vault.UpdateCapability},
	"/root/generate/external": {vault.CreateCapability, vault.UpdateCapability},
	"/root/sign-intermediate": {vault.CreateCapability, vault.UpdateCapability},
	"/root/sign-self-issued":  {vault.CreateCapability, vault.UpdateCapability},
	"/sign-verbatim":          {vault.CreateCapability, vault.UpdateCapability},
	"/tidy":                   {vault.CreateCapability, vault.UpdateCapability},
}

var (
	_ cli.Command             = (*PKIHealthCheckCommand)(nil)
	_ cli.CommandAutocomplete = (*PKIHealthCheckCommand)(nil)
)

func maxInt(a int, b int) int {
	if a > b {
		return a
	}

	return b
}

func sortedKeys(checks map[string]PKIHealthCheckHelper) []string {
	var result []string
	for name := range checks {
		result = append(result, name)
	}
	sort.Strings(result)
	return result
}

type PKIHealthCheckCommand struct {
	*BaseCommand

	crlConfig map[string]interface{}
	roles     map[string]map[string]interface{}
	issuer    *x509.Certificate

	mountDefaultTTL        time.Duration
	mountMaximumTTL        time.Duration
	fetchedMountAuditKeys  bool
	mountAuditRequestKeys  []string
	mountAuditResponseKeys []string

	flagChecks          string
	flagExclude         string
	flagMaxFetchedCerts int
}

type PKIHealthCheckHelper func(mount string) int

func (c *PKIHealthCheckCommand) Synopsis() string {
	return "Check PKI Secrets Engine health and operational status"
}

func (c *PKIHealthCheckCommand) Help() string {
	helpText := `
Usage: vault pki health-check [options] MOUNT

  Reports status of the mount against best practices and pending failures,
  including default TTL period of roles. This is an informative command and
  not meant to be automatically parsed.

  To check the pki-root mount:

      $ vault pki health-check /pki-root

  To exclude specific checks:

      $ vault pki health-check -exclude=role-localhost-issuance /pki-root

  Known checks:
      crl-validity            - checks if the CRL is still valid
      direct-issuance         - checks if a root is directly issuing leaf certs
      issuer-validity         - checks if the issuing CA certificate is still valid
      mount-audit-keys        - checks if mount is auditable
      policy-scope            - checks if policies allows problematic role endpoints
      role-glob-wildcard      - checks if the role allows wildcard issuance via globs
      role-localhost-issuance - checks if the role allows localhost issuance
      role-ttl-crl            - checks role TTLs for no_store and CRL usage

  Return codes indicate failure type:

      0 - Everything is good.
      1 - Usage error (check CLI parameters).
      2 - Internal error (check CLI parameters and error message).
      3 - Informational warnings that may impact the PKI mount in the future.
      4 - Critical or impending failure of specified PKI mount.

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *PKIHealthCheckCommand) Flags() *FlagSets {
	set := c.flagSet(FlagSetHTTP)
	f := set.NewFlagSet("Command Options")

	f.StringVar(&StringVar{
		Name:    "checks",
		Target:  &c.flagChecks,
		Default: "all",
		EnvVar:  "",
		Usage:   "Comma separate list of checks to run; use all to execute all known checks.",
	})

	f.StringVar(&StringVar{
		Name:    "exclude",
		Target:  &c.flagExclude,
		Default: "",
		EnvVar:  "",
		Usage:   "Comma separate list of checks to exclude.",
	})

	f.DurationVar(&DurationVar{
		Name:    "ttl",
		Target:  &c.mountDefaultTTL,
		Default: 0,
		EnvVar:  "",
		Usage: "The default lease TTL for this secrets engine, if permission " +
			"to read /sys/internal/ui/mounts and /sys/config/state/sanitized " +
			"is lacked by this token but is otherwise known. Will be assumed " +
			"to be 768 hours (Vault default) if this option is not supplied " +
			"and token lacks required permissions.",
	})

	f.DurationVar(&DurationVar{
		Name:    "max-ttl",
		Target:  &c.mountMaximumTTL,
		Default: 0,
		EnvVar:  "",
		Usage: "The maximum lease TTL for this secrets engine, if permission " +
			"to read /sys/internal/ui/mounts and /sys/config/state/sanitized " +
			"is lacked by this token but is otherwise known. Will be assumed " +
			"to be 768 hours (Vault default) if this option is not supplied " +
			"and token lacks required permissions.",
	})

	f.IntVar(&IntVar{
		Name:    "max-fetched-certs",
		Target:  &c.flagMaxFetchedCerts,
		Default: 100,
		EnvVar:  "",
		Usage: "The maximum number of certificates to fetch when running the " +
			"direct-issuance health check (to ensure root CAs haven't issued any " +
			"non-CA leaf certificates",
	})

	return set
}

func (c *PKIHealthCheckCommand) AutocompleteArgs() complete.Predictor {
	// Return an anything predictor here, similar to `vault write`. We
	// don't know what values are valid for the mount path.
	return complete.PredictAnything
}

func (c *PKIHealthCheckCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *PKIHealthCheckCommand) Run(args []string) int {
	f := c.Flags()

	if err := f.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return pkiRetUsage
	}

	args = f.Args()
	if len(args) < 1 {
		c.UI.Error("Not enough arguments (expected mount path, got nothing)")
		return pkiRetUsage
	} else if len(args) > 1 {
		c.UI.Error(fmt.Sprintf("Too many arguments (expected only mount path, got %d arguments)", len(args)))
		for _, arg := range args {
			if strings.HasPrefix(arg, "-") {
				c.UI.Warn(fmt.Sprintf("Options (%v) must be specified before positional arguments (%v)", arg, args[0]))
				break
			}
		}
		return pkiRetUsage
	}

	mount := sanitizePath(args[0])

	// Don't forget to update the help text.
	allChecks := map[string]PKIHealthCheckHelper{
		"crl-validity":            c.checkCRLValidity,
		"direct-issuance":         c.checkDirectIssuance,
		"issuer-validity":         c.checkIssuerValidity,
		"mount-audit-keys":        c.checkMountAuditKeys,
		"policy-scope":            c.checkPolicyScope,
		"role-glob-wildcard":      c.checkRoleGlobWildcard,
		"role-localhost-issuance": c.checkRoleLocalhost,
		"role-ttl-crl":            c.checkRoleTTLAndCRL,
	}

	thisRet := pkiRetOK
	for _, name := range sortedKeys(allChecks) {
		checker := allChecks[name]
		if c.flagChecks != "all" && !strings.Contains(c.flagChecks, name) {
			continue
		}
		if strings.Contains(c.flagExclude, name) {
			continue
		}

		c.UI.Info(fmt.Sprintf("Health check: %v", name))
		checkRet := checker(mount)
		if checkRet >= pkiRetUsage && checkRet <= pkiRetInternal {
			c.UI.Error(fmt.Sprintf("\tFailed with an internal error (%v); check errors and arguments and try again.", checkRet))
			return thisRet
		}

		if checkRet == pkiRetOK {
			c.UI.Info("\tOK\n")
		}

		thisRet = maxInt(thisRet, thisRet)
	}

	return thisRet
}

func decodePem(cert string) ([]byte, error) {
	block, extra := pem.Decode([]byte(cert))
	if strings.TrimSpace(string(extra)) != "" {
		return nil, errors.New("unexpected trailing data after certificate: " + string(extra))
	}

	return block.Bytes, nil
}

func getCertResponse(data map[string]interface{}) (string, error) {
	pemCert, ok := data["certificate"]
	if !ok || len(pemCert.(string)) == 0 {
		return "", errors.New("missing certificate field in response")
	}

	return pemCert.(string), nil
}

func decodeCertResponse(data map[string]interface{}) (*x509.Certificate, error) {
	pemCert, err := getCertResponse(data)
	if err != nil {
		return nil, err
	}

	pemBytes, err := decodePem(pemCert)
	if err != nil {
		return nil, err
	}

	return x509.ParseCertificate(pemBytes)
}

func decodeCRLResponse(data map[string]interface{}) (*pkix.CertificateList, error) {
	pemCert, err := getCertResponse(data)
	if err != nil {
		return nil, err
	}

	pemBytes, err := decodePem(pemCert)
	if err != nil {
		return nil, err
	}

	return x509.ParseCRL(pemBytes)
}

func (c *PKIHealthCheckCommand) queryIssuer(mount string) error {
	if c.issuer != nil {
		return nil
	}

	value, err := c.queryCert(mount, "ca")
	if err != nil {
		return err
	}

	c.issuer = value
	return nil
}

func (c *PKIHealthCheckCommand) queryCert(mount string, serial string) (*x509.Certificate, error) {
	client, err := c.Client()
	if err != nil {
		return nil, err
	}

	path := mount + "/cert/" + serial
	certData, err := client.Logical().Read(path)
	if err != nil {
		return nil, err
	}

	if certData == nil || certData.Data == nil {
		return nil, fmt.Errorf("Expected certificate (%v) to have data", serial)
	}

	cert, err := decodeCertResponse(certData.Data)
	if err != nil {
		c.UI.Info(fmt.Sprintf("\tAPI response: %v\n", certData.Data))
		return nil, err
	}

	return cert, nil
}

func (c *PKIHealthCheckCommand) checkIssuerValidity(mount string) int {
	err := c.queryIssuer(mount)
	if err != nil {
		c.UI.Error("\t" + err.Error())
		return pkiRetInternal
	}

	now := time.Now()
	oneWeek := now.Add(7 * oneDay)
	threeMonthsFromNow := now.Add(3 * oneMonth)
	oneYearFromNow := now.Add(oneYear)

	thisRet := pkiRetOK

	if c.issuer.NotBefore.After(now) {
		c.UI.Info(fmt.Sprintf("\tIssuer certificiate not yet valid:\n\t\tNotBefore=%v > Now=%v\n", c.issuer.NotBefore.Format(time.RFC1123), now.Format(time.RFC1123)))
		thisRet = maxInt(thisRet, pkiRetInformational)
	}

	if c.issuer.NotAfter.Before(now) {
		c.UI.Error(fmt.Sprintf("\tIssuer certificiate no longer valid:\n\t\tNotAfter=%v < Now=%v\n", c.issuer.NotAfter.Format(time.RFC1123), now.Format(time.RFC1123)))
		thisRet = maxInt(thisRet, pkiRetCritical)
	} else if c.issuer.NotAfter.Before(oneWeek) {
		c.UI.Error(fmt.Sprintf("\tIssuer certificiate expiring within a week:\n\t\tNotAfter=%v < Now+7 Days=%v\n", c.issuer.NotAfter.Format(time.RFC1123), oneWeek.Format(time.RFC1123)))
		thisRet = maxInt(thisRet, pkiRetCritical)
	} else if c.issuer.NotAfter.Before(threeMonthsFromNow) {
		c.UI.Warn(fmt.Sprintf("\tIssuer certificiate expiring within three months:\n\t\tNotAfter=%v < Now+3 Months=%v\n", c.issuer.NotAfter.Format(time.RFC1123), threeMonthsFromNow.Format(time.RFC1123)))
		thisRet = maxInt(thisRet, pkiRetCritical)
	} else if c.issuer.NotAfter.Before(oneYearFromNow) {
		c.UI.Info(fmt.Sprintf("\tIssuer certificiate expiring within a year:\n\t\tNotAfter=%v < Now+1 Year=%v\n", c.issuer.NotAfter.Format(time.RFC1123), oneYearFromNow.Format(time.RFC1123)))
		thisRet = maxInt(thisRet, pkiRetInformational)
	}

	return thisRet
}

func (c *PKIHealthCheckCommand) queryCRLConfig(mount string) error {
	if c.crlConfig != nil {
		return nil
	}

	client, err := c.Client()
	if err != nil {
		return err
	}

	path := mount + "/config/crl"

	configData, err := client.Logical().Read(path)
	if err != nil {
		if strings.Contains(err.Error(), "permission denied") {
			c.UI.Error(fmt.Sprintf("\tunable to read %v; check token ACLs or exclude this check and try again", path))
		}

		return err
	}

	if configData == nil || configData.Data == nil {
		// Vault appears to return empty responses if the config hasn't yet been
		// written. Build our own config struct instead.
		c.crlConfig = make(map[string]interface{})
		c.crlConfig["disable"] = false
		c.crlConfig["expiry"] = "72h"
		return nil
	}

	// Handle explicitly setting the default
	if configData.Data["expiry"].(string) == "" {
		configData.Data["expiry"] = "72"
	}

	c.crlConfig = configData.Data
	return nil
}

func (c *PKIHealthCheckCommand) checkCRLValidity(mount string) int {
	err := c.queryCRLConfig(mount)
	if err != nil {
		c.UI.Error("\t" + err.Error())
		return pkiRetInternal
	}

	if c.crlConfig["disable"].(bool) {
		c.UI.Info("\tCRL disabled\n")
		return pkiRetOK
	}

	client, err := c.Client()
	if err != nil {
		c.UI.Error("\t" + err.Error())
		return pkiRetInternal
	}

	path := mount + "/cert/crl"

	certData, err := client.Logical().Read(path)
	if err != nil {
		c.UI.Error("\t" + err.Error())
		return pkiRetInternal
	}

	if certData == nil || certData.Data == nil {
		c.UI.Error("\tExpected CRL to have data.\n")
		return pkiRetInternal
	}

	crl, err := decodeCRLResponse(certData.Data)
	if err != nil {
		c.UI.Error("\t" + err.Error())
		c.UI.Info(fmt.Sprintf("API response: %v", certData.Data))
		return pkiRetInternal
	}

	now := time.Now()
	if crl.HasExpired(now) {
		rotateCommand := "$ vault read " + mount + "/crl/rotate"
		c.UI.Error(fmt.Sprintf("\tCRL should have been rotated by now but wasn't; to rotate, run:\n\t\t%v\n", rotateCommand))
		return pkiRetCritical
	}

	return pkiRetOK
}

func (c *PKIHealthCheckCommand) queryRoles(mount string) error {
	if c.roles != nil {
		return nil
	}

	client, err := c.Client()
	if err != nil {
		return err
	}

	path := mount + "/roles"
	data, err := client.Logical().List(path)
	if err != nil {
		if strings.Contains(err.Error(), "permission denied") {
			c.UI.Error(fmt.Sprintf("\tunable to list %v; check token ACLs or exclude this check and try again", path))
		}
		return err
	}

	if data == nil || data.Data == nil {
		return errors.New("expected non-empty result from listing roles")
	}

	roles, ok := data.Data["keys"]
	if !ok {
		return errors.New("expected listing roles to return data")
	}

	roleMap := make(map[string]map[string]interface{})
	for _, rawName := range roles.([]interface{}) {
		name := rawName.(string)
		rolePath := mount + "/roles/" + name
		roleData, err := client.Logical().Read(rolePath)
		if err != nil {
			if strings.Contains(err.Error(), "permission denied") {
				c.UI.Error(fmt.Sprintf("\tunable to read %v; check token ACLs or exclude this check and try again", rolePath))
			}
			return err
		}

		if roleData == nil || roleData.Data == nil {
			return fmt.Errorf("got empty result querying role: %v", name)
		}

		roleMap[name] = roleData.Data
	}

	c.roles = roleMap
	return nil
}

func (c *PKIHealthCheckCommand) checkRoleLocalhost(mount string) int {
	err := c.queryRoles(mount)
	if err != nil {
		c.UI.Error("\t" + err.Error())
		return pkiRetInternal
	}

	thisRet := pkiRetOK

	for name, role := range c.roles {
		allow_localhost := role["allow_localhost"].(bool)
		allow_subdomains := role["allow_subdomains"].(bool)
		allow_glob_domains := role["allow_glob_domains"].(bool)
		allow_bare_domains := role["allow_bare_domains"].(bool)
		allowed_domains := role["allowed_domains"].([]interface{})
		baseCommand := "$ vault pki role-update " + mount + "/roles/" + name

		if !allow_localhost {
			continue
		}

		have_explicit_localhost := false
		var present_domains []string

		matched_domain := ""
		for _, rawDomain := range allowed_domains {
			domain := rawDomain.(string)
			present_domains = append(present_domains, domain)

			if strings.HasSuffix(domain, "localhost") || strings.HasSuffix(domain, "localdomain") {
				have_explicit_localhost = true
				matched_domain = domain
			}
		}

		if have_explicit_localhost {
			updateCommand := baseCommand + " allow_localhost=false"
			msg := "\tRole %v has duplicated value allow_localhost=true and a\n\t"
			msg += "localhost domain in allowed_domains (%v).\n\t"
			msg += "Consider setting allow_localhost=false to use\n\t"
			msg += "allowed_domain matching syntax instead. To update the role, run:\n\t\t%v\n"
			c.UI.Warn(fmt.Sprintf(msg, name, matched_domain, updateCommand))
			thisRet = maxInt(thisRet, pkiRetInformational)
		} else if len(allowed_domains) > 0 {
			newDomains := strings.Join(present_domains, ",")
			forbidCommand := baseCommand + " allow_localhost=false"
			if allow_subdomains || (allow_glob_domains && allow_bare_domains) {
				// Here we can construct an alternative command to allow localhost.
				if allow_subdomains {
					newDomains += ",localhost,localdomain"
				} else if allow_glob_domains && allow_bare_domains {
					newDomains += ",*.localhost,*.localdomain,localhost,localdomain"
				}
				allowCommand := baseCommand + " allowed_domains=" + newDomains + " allow_localhost=false"
				msg := "\tRole %v has allow_localhost=true and a non-empty list\n\t"
				msg += "in allowed_domains. Consider setting allow_localhost=false and either\n\t"
				msg += "explicitly add localhost to allowed_domains via:\n\t\t%v\n"
				msg += "\tOr disable localhost issuance entirely on this role:\n\t\t%v\n"
				c.UI.Warn(fmt.Sprintf(msg, name, allowCommand, forbidCommand))
				thisRet = maxInt(thisRet, pkiRetInformational)
			} else {
				msg := "\tRole %v has allow_localhost=true and a non-empty list in allowed_domains.\n\t"
				msg += "Consider setting allow_localhost=false and either explicitly add localhost\n\t"
				msg += "to allowed_domains or disable localhost issuance entirely on this role:\n\t\t%v\n"
				c.UI.Warn(fmt.Sprintf(msg, name, forbidCommand))
				thisRet = maxInt(thisRet, pkiRetInformational)
			}
		}
	}

	return thisRet
}

func (c *PKIHealthCheckCommand) checkRoleGlobWildcard(mount string) int {
	err := c.queryRoles(mount)
	if err != nil {
		c.UI.Error("\t" + err.Error())
		return pkiRetInternal
	}

	thisRet := pkiRetOK

	for name, role := range c.roles {
		raw_allow_wildcards, have_wildcards := role["allow_wildcard_certificates"]
		if !have_wildcards {
			continue
		}

		allow_wildcards := raw_allow_wildcards.(bool)
		allow_glob_domains := role["allow_glob_domains"].(bool)
		allowed_domains := role["allowed_domains"].([]interface{})
		baseCommand := "$ vault pki role-update " + mount + "/roles/" + name

		if !allow_glob_domains || !allow_wildcards {
			continue
		}

		matched_domain := ""
		for _, rawDomain := range allowed_domains {
			domain := rawDomain.(string)

			if strings.HasPrefix(domain, "*") && !strings.HasSuffix(domain, "localhost") && !strings.HasSuffix(domain, "localdomain") {
				matched_domain = domain
			}
		}

		if matched_domain == "" {
			continue
		}

		forbidCommand := baseCommand + " allow_wildcard_certificates=false"
		msg := "\tRole %v has domain with a leading glob pattern (%v),\n\t"
		msg += "potentially matching wildcard certificate requests. Consider setting\n\t"
		msg += "allow_wildcard_certificates=false to forbid issuance of wildcards\n\t"
		msg += "under this role. To update the role, run:\n\t\t%v\n"
		c.UI.Warn(fmt.Sprintf(msg, name, matched_domain, forbidCommand))
		thisRet = maxInt(thisRet, pkiRetInformational)
	}

	return thisRet
}

func (c *PKIHealthCheckCommand) queryMountTTL(mount string) {
	if c.mountDefaultTTL != 0 || c.mountMaximumTTL != 0 {
		return
	}

	// By using goto, we have to predeclare our later values
	var ok bool
	var path string
	var resp *api.Secret
	var secretMounts map[string]interface{}
	var currentMount string
	var rawMountDetails interface{}
	var mountDetails map[string]interface{}
	var config map[string]interface{}
	var rawMountTTLValue int64
	var rawMountMaxTTLValue int64
	var rawConfigTTLValue int64
	var rawConfigMaxTTLValue int64
	var rawAuditValues []interface{}
	var rawAuditValue interface{}

	client, err := c.Client()
	if err != nil {
		goto Error
	}

	path = "/sys/internal/ui/mounts"
	resp, err = client.Logical().Read(path)
	if err != nil {
		if strings.Contains(err.Error(), "permission denied") {
			c.UI.Warn(fmt.Sprintf("\ttoken lacks permission to read %v; specify -ttl and -max-ttl to provide an accurate value of mount's default TTL\n", path))
		}
		goto Error
	}

	if resp == nil || resp.Data == nil {
		goto Error
	}

	secretMounts = resp.Data["secret"].(map[string]interface{})
	for currentMount, rawMountDetails = range secretMounts {
		if currentMount != mount+"/" {
			continue
		}

		mountDetails = rawMountDetails.(map[string]interface{})
		config = mountDetails["config"].(map[string]interface{})
		rawMountTTLValue, _ = config["default_lease_ttl"].(json.Number).Int64()
		rawMountMaxTTLValue, _ = config["max_lease_ttl"].(json.Number).Int64()

		c.mountDefaultTTL = time.Duration(rawMountTTLValue) * time.Second
		c.mountMaximumTTL = time.Duration(rawMountMaxTTLValue) * time.Second

		// Check for HMAC while were here.
		rawAuditValues, ok = config["audit_non_hmac_request_keys"].([]interface{})
		if ok {
			for _, rawAuditValue = range rawAuditValues {
				c.mountAuditRequestKeys = append(c.mountAuditRequestKeys, rawAuditValue.(string))
			}
		}
		rawAuditValues, ok = config["audit_non_hmac_response_keys"].([]interface{})
		if ok {
			for _, rawAuditValue = range rawAuditValues {
				c.mountAuditResponseKeys = append(c.mountAuditResponseKeys, rawAuditValue.(string))
			}
		}

		c.fetchedMountAuditKeys = true

		break
	}

	// Only query config if we haven't gotten one of our TTL values yet.
	if c.mountDefaultTTL != 0 && c.mountMaximumTTL != 0 {
		return
	}

	path = "/sys/config/state/sanitized"
	resp, err = client.Logical().Read(path)
	if err != nil {
		if strings.Contains(err.Error(), "permission denied") {
			c.UI.Warn(fmt.Sprintf("\ttoken lacks permission to read %v; specify -ttl and -max-ttl to provide an accurate value of mount's default TTL\n", path))
		}
		goto Error
	}

	if resp == nil || resp.Data == nil {
		goto Error
	}

	rawConfigTTLValue, _ = resp.Data["default_lease_ttl"].(json.Number).Int64()
	rawConfigMaxTTLValue, _ = resp.Data["max_lease_ttl"].(json.Number).Int64()

	if c.mountDefaultTTL == 0 {
		if rawConfigTTLValue == 0 {
			c.mountDefaultTTL = 768 * time.Hour
		} else {
			c.mountDefaultTTL = time.Duration(rawConfigTTLValue) * time.Second
		}
	}

	if c.mountMaximumTTL == 0 {
		if rawConfigMaxTTLValue == 0 {
			c.mountMaximumTTL = 768 * time.Hour
		} else {
			c.mountMaximumTTL = time.Duration(rawConfigMaxTTLValue) * time.Second
		}
	}

	return

Error:
	c.mountDefaultTTL = 768 * time.Hour
	c.mountMaximumTTL = 768 * time.Hour
}

func (c *PKIHealthCheckCommand) checkRoleTTLAndCRL(mount string) int {
	err := c.queryCRLConfig(mount)
	if err != nil {
		c.UI.Error("\t" + err.Error())
		return pkiRetInternal
	}

	err = c.queryRoles(mount)
	if err != nil {
		c.UI.Error("\t" + err.Error())
		return pkiRetInternal
	}

	c.queryMountTTL(mount)

	crl_disabled := c.crlConfig["disable"].(bool)

	thisRet := pkiRetOK
	warnedCRLDisabled := false

	for name, role := range c.roles {
		rawRoleTTL, _ := role["ttl"].(json.Number).Int64()
		roleTTL := time.Duration(rawRoleTTL) * time.Second
		rawRoleMaxTTL, _ := role["max_ttl"].(json.Number).Int64()
		roleMaxTTL := time.Duration(rawRoleMaxTTL) * time.Second
		no_store := role["no_store"].(bool)

		if roleTTL == 0 {
			roleTTL = c.mountDefaultTTL
		}
		if roleMaxTTL == 0 {
			roleMaxTTL = c.mountMaximumTTL
		}

		maxTTLDays := maxInt(int(roleTTL/oneDay), int(roleMaxTTL/oneDay))

		if !warnedCRLDisabled && crl_disabled && (roleTTL > definitelyCRL || roleMaxTTL > definitelyCRL) {
			enableCommand := "$ vault write " + mount + "/config/crl disable=false"
			msg := "\tRole (%v) or mount/system TTL exceeds %v days (%v days); it is strongly encouraged\n"
			msg += "\tencouraged to re-enable CRL generation:\n\t\t%v\n"
			c.UI.Error(fmt.Sprintf(msg, name, int(definitelyCRL/oneDay), maxTTLDays, enableCommand))
			thisRet = maxInt(thisRet, pkiRetCritical)
			warnedCRLDisabled = true
		} else if !warnedCRLDisabled && crl_disabled && (roleTTL > suggestedCRL || roleMaxTTL > suggestedCRL) {
			enableCommand := "$ vault write " + mount + "/config/crl disable=false"
			msg := "\tRole (%v) or mount/system TTL exceeds %v days (%v days); it is encouraged to\n"
			msg += "\tre-enable CRL generation:\n\t\t%v\n"
			c.UI.Warn(fmt.Sprintf(msg, name, int(definitelyCRL/oneDay), maxTTLDays, enableCommand))
			thisRet = maxInt(thisRet, pkiRetInformational)
			warnedCRLDisabled = true
		}

		if no_store && (roleTTL > suggestedCRL || roleMaxTTL > suggestedCRL) {
			storeCmd := "$ vault pki role-update " + mount + "/roles/" + name + " no_store=false"
			msg := "\tRole (%v) sets no_store=true, preventing certificates generated under\n"
			msg += "\tthis role from ever being added to the CRL. However, because this role\n"
			msg += "\tcan issue long-lived certificates (%v days), it is encouraged to set\n"
			msg += "\tno_store=false via:\n\t\t%v\n"
			c.UI.Error(fmt.Sprintf(msg, name, maxTTLDays, storeCmd))
			thisRet = maxInt(thisRet, pkiRetCritical)
		} else if !no_store && roleTTL < suggestNoCRL && roleMaxTTL < suggestNoCRL {
			noStoreCmd := "$ vault pki role-update " + mount + "/roles/" + name + " no_store=true"
			msg := "\tRole (%v) sets no_store=false, causing certificates generated under this\n"
			msg += "\trole to be tracked by Vault storage. However, due to their short lifetime\n"
			msg += "\t(under %v days), it is suggested not to store these certificates, to avoid\n"
			msg += "\ta high storage burden on Vault. Consider setting no_store=true via:\n\t\t%v\n"
			c.UI.Error(fmt.Sprintf(msg, name, int(suggestNoCRL/oneDay), noStoreCmd))
			thisRet = maxInt(thisRet, pkiRetCritical)
		}

		// This role was explicitly configured with really weird parameters which
		// weren't inherited from the mount or the system.
		if roleTTL < suggestNoCRL && roleMaxTTL > suggestedCRL && rawRoleTTL != 0 && rawRoleMaxTTL != 0 {
			msg := "\tRole (%v) has a short TTL (%v) but a long maximum TTL (%v).\n"
			msg += "\tConsider aligning the values or splitting the role into two to utilize\n"
			msg += "\tno_store=true for short-lived certificates while allowing revocation\n"
			msg += "\tfor long-lived certificates.\n"
			c.UI.Warn(fmt.Sprintf(msg, name, roleTTL, roleMaxTTL))
			thisRet = maxInt(thisRet, pkiRetInformational)
		}
	}

	return thisRet
}

func isCA(cert *x509.Certificate) bool {
	return cert.BasicConstraintsValid && cert.IsCA
}

func isRootCA(cert *x509.Certificate) bool {
	if !isCA(cert) || !bytes.Equal(cert.RawSubject, cert.RawIssuer) {
		return false
	}

	// In order to be sure this is truly a root CA and not a weird configuration
	// (like an intermediate signed by a different set of key material with the
	// same subject) or something weird like that, we need to validate the
	// signature and ensure _this_ key signed that other key.
	pool := x509.NewCertPool()
	pool.AddCert(cert)
	var options x509.VerifyOptions
	options.Roots = pool

	// Ensure the cert is valid at the of verification so that doesn't affect
	// our signature check.
	options.CurrentTime = cert.NotAfter.Add(-1 * time.Second)

	// Since we only care about the signature, ignore KeyUsages.
	options.KeyUsages = append(options.KeyUsages, x509.ExtKeyUsageAny)

	// Verify the cert and check the signature.
	_, err := cert.Verify(options)
	return err == nil
}

func (c *PKIHealthCheckCommand) checkDirectIssuance(mount string) int {
	err := c.queryIssuer(mount)
	if err != nil {
		c.UI.Error("\t" + err.Error())
		return pkiRetInternal
	}

	if !isRootCA(c.issuer) {
		return pkiRetOK
	}

	client, err := c.Client()
	if err != nil {
		c.UI.Error("\t" + err.Error())
		return pkiRetInternal
	}

	path := mount + "/certs"
	data, err := client.Logical().List(path)
	if err != nil {
		if strings.Contains(err.Error(), "permission denied") {
			c.UI.Error(fmt.Sprintf("\tunable to list %v; check token ACLs or exclude this check and try again", path))
		}
		c.UI.Error("\t" + err.Error())
		return pkiRetInternal
	}

	if data == nil || data.Data == nil {
		c.UI.Error("\tExpected non-empty result from listing roles")
		return pkiRetInternal
	}

	certs, ok := data.Data["keys"]
	if !ok {
		c.UI.Error("\tExpected listing roles to return data")
		return pkiRetInternal
	}

	thisRet := pkiRetOK

	var fetched int
	for _, rawSerial := range certs.([]interface{}) {
		if fetched > c.flagMaxFetchedCerts {
			break
		}

		fetched += 1
		serial := rawSerial.(string)
		cert, err := c.queryCert(mount, serial)
		if err != nil {
			continue
		}

		if !isCA(cert) {
			msg := "\tThis mount has a root certificate which has issued a leaf certificate which isn't\n"
			msg += "\ta CA certificate (%v).\n"
			msg += "\tIt is suggested to only use intermediate certificates for issuance\n"
			msg += "\tas they are easier to rotate than root certificates.\n"
			c.UI.Error(fmt.Sprintf(msg, serial))
			thisRet = pkiRetCritical
			break
		}
	}

	return thisRet
}

func sliceHasValue(slice []string, target string) bool {
	for _, value := range slice {
		if value == target {
			return true
		}
	}
	return false
}

func (c *PKIHealthCheckCommand) checkMountAuditKeys(mount string) int {
	c.queryMountTTL(mount)

	if !c.fetchedMountAuditKeys {
		c.UI.Warn("\tUnable to verify mount audit keys, likely due to lack of read permission on\n\t/sys/internal/ui/mounts")
		return pkiRetInformational
	}

	var missingReqKeys []string
	for _, key := range pkiShouldAuditReqKeys {
		if !sliceHasValue(c.mountAuditRequestKeys, key) {
			missingReqKeys = append(missingReqKeys, key)
		}
	}

	var missingRespKeys []string
	for _, key := range pkiShouldAuditRespKeys {
		if !sliceHasValue(c.mountAuditResponseKeys, key) {
			missingRespKeys = append(missingRespKeys, key)
		}
	}

	thisRet := pkiRetOK
	if len(missingReqKeys) > 0 {
		auditCommand := "$ vault secrets tune -audit-non-hmac-request-keys=%v %v"
		auditCommand = fmt.Sprintf(auditCommand, strings.Join(pkiShouldAuditReqKeys, " -audit-non-hmac-request-keys="), mount)
		msg := "\tMount is missing useful auditing information because of HMAC'd request keys.\n"
		msg += "\tConsider un-HMACing these request parameters via:\n\t\t%v\n"
		c.UI.Warn(fmt.Sprintf(msg, auditCommand))
		thisRet = pkiRetInformational
	}

	if len(missingRespKeys) > 0 {
		auditCommand := "$ vault secrets tune -audit-non-hmac-response-keys=%v %v"
		auditCommand = fmt.Sprintf(auditCommand, strings.Join(pkiShouldAuditRespKeys, " -audit-non-hmac-response-keys="), mount)
		msg := "\tMount is missing useful auditing information because of HMAC'd response keys.\n"
		msg += "\tConsider un-HMACing these response parameters via:\n\t\t%v\n"
		c.UI.Warn(fmt.Sprintf(msg, auditCommand))
		thisRet = pkiRetInformational
	}

	return thisRet
}

func (c *PKIHealthCheckCommand) checkPolicyScope(mount string) int {
	client, err := c.Client()
	if err != nil {
		c.UI.Error("\t" + err.Error())
		return pkiRetInternal
	}

	path := "/sys/policy"
	policyList, err := client.Logical().List(path)
	if err != nil {
		if strings.Contains(err.Error(), "permission denied") {
			c.UI.Warn(fmt.Sprintf("\tunable to list %v; check token ACLs or exclude this check and try again", path))
		} else {
			c.UI.Error("\t" + err.Error())
		}
		return pkiRetInternal
	}

	if policyList == nil || policyList.Data == nil {
		c.UI.Warn("\tGot empty list of policies")
		return pkiRetInternal
	}

	thisRet := pkiRetOK

	// XXX: Handle namespaces correctly?
	ctx := namespace.ContextWithNamespace(context.Background(), namespace.RootNamespace)
	rawKeys := policyList.Data["keys"].([]interface{})
	for _, rawPolicyName := range rawKeys {
		policyName := rawPolicyName.(string)
		path = "/sys/policy/" + policyName

		policyData, err := client.Logical().Read(path)
		if err != nil {
			if strings.Contains(err.Error(), "permission denied") {
				c.UI.Warn(fmt.Sprintf("\tunable to read %v; check token ACLs or exclude this check and try again", path))
			} else {
				c.UI.Error("\t" + err.Error())
			}
			return pkiRetInternal
		}

		if policyData == nil || policyData.Data == nil {
			c.UI.Warn("\tGot empty policy description")
			return pkiRetInternal
		}

		rules, ok := policyData.Data["rules"].(string)
		if !ok {
			c.UI.Warn("\tPolicy " + policyName + " lacks rules data")
			continue
		}

		policy, err := vault.ParseACLPolicy(namespace.RootNamespace, rules)
		if err != nil {
			c.UI.Error(fmt.Sprintf("\tFailed to parse policy %v: %v", policyName, err))
			return pkiRetInternal
		}
		policy.Name = policyName

		acl, err := vault.NewACL(ctx, []*vault.Policy{policy})
		if err != nil {
			c.UI.Error(fmt.Sprintf("\tFailed to create ACL from policy: %v: %v", policyName, err))
		}

		isPrivilegedPolicy := false
		for path, capabilities := range pkiPrivilegedPolicyPaths {
			if isPrivilegedPolicy {
				break
			}

			for _, capability := range capabilities {
				mockRequest := &logical.Request{
					Operation: logical.Operation(capability),
					Path:      mount + path,
				}

				ret := acl.AllowOperation(ctx, mockRequest, false)
				if ret.Allowed {
					isPrivilegedPolicy = true
					c.UI.Error(fmt.Sprintf("\tPolicy %v allows privileged PKI endpoint %v with %v capabilities.\n", policyName, path, capability))
					thisRet = maxInt(thisRet, pkiRetInformational)
					break
				}
			}
		}
	}

	return thisRet
}
