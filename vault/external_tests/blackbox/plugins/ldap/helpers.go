// Copyright IBM Corp. 2025, 2026
// SPDX-License-Identifier: BUSL-1.1

package ldap

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/helper/testcluster/blackbox"
)

// LDAPDomainConfig holds configuration for an isolated LDAP domain
type LDAPDomainConfig struct {
	URL      string // LDAP server URL (private IP for Vault operations)
	SetupURL string // LDAP server URL (public IP for setup operations from GitHub runner)
	BindDN   string // Admin bind DN
	BindPass string // Admin bind password
	BaseDN   string // Base DN for this domain (e.g., dc=bbsdk-xxx,dc=test,dc=enos,dc=com)
	UserDN   string // User DN (e.g., ou=users,dc=bbsdk-xxx,dc=test,dc=enos,dc=com)
	GroupDN  string // Group DN (e.g., ou=groups,dc=bbsdk-xxx,dc=test,dc=enos,dc=com)
}

// setupLDAPSecretsEngine configures the LDAP secrets engine with environment variables
// Deprecated: Use SetupLDAPSecretsEngineWithConfig for isolated domain support
func setupLDAPSecretsEngine(t *testing.T, v *blackbox.Session, mount string) {
	// Enable LDAP secrets engine
	v.MustEnableSecretsEngine(mount, &api.MountInput{Type: "ldap"})

	// Configure using environment variables set by ldap.tf
	ldapServer := os.Getenv("LDAP_SERVER")
	ldapPort := os.Getenv("LDAP_PORT")
	ldapBindDN := os.Getenv("LDAP_BIND_DN")
	ldapBindPass := os.Getenv("LDAP_BIND_PASS")
	ldapUsername := os.Getenv("LDAP_USERNAME") // "enos"

	if ldapServer == "" || ldapPort == "" || ldapBindDN == "" || ldapBindPass == "" {
		t.Fatal("Required LDAP environment variables not set")
	}

	v.MustWrite(mount+"/config", map[string]any{
		"binddn":   ldapBindDN,
		"bindpass": ldapBindPass,
		"url":      "ldap://" + ldapServer + ":" + ldapPort,
		"userdn":   "ou=users,dc=" + ldapUsername + ",dc=com",
	})
}

// SetupLDAPSecretsEngineWithConfig configures the LDAP secrets engine with an isolated domain config
func SetupLDAPSecretsEngineWithConfig(t *testing.T, v *blackbox.Session, mount string, config *LDAPDomainConfig) {
	t.Helper()

	// Enable LDAP secrets engine
	v.MustEnableSecretsEngine(mount, &api.MountInput{Type: "ldap"})

	// Configure using isolated domain config
	v.MustWrite(mount+"/config", map[string]any{
		"binddn":   config.BindDN,
		"bindpass": config.BindPass,
		"url":      config.URL,
		"userdn":   config.UserDN,
		"groupdn":  config.GroupDN,
	})

	t.Logf("Configured LDAP secrets engine at %s with isolated domain %s", mount, config.BaseDN)
	t.Logf("DEBUG: LDAP Configuration:")
	t.Logf("  URL: %s", config.URL)
	t.Logf("  BindDN: %s", config.BindDN)
	t.Logf("  UserDN: %s", config.UserDN)
	t.Logf("  GroupDN: %s", config.GroupDN)
	t.Logf("DEBUG: To test LDAP connection manually, run:")
	t.Logf("  ldapsearch -x -H %s -b \"%s\" -D %s -w %s",
		config.URL, config.BaseDN, config.BindDN, config.BindPass)
}

// WriteLibrarySetWithRetry writes an LDAP library set configuration with retry logic
// LDAP library operations can take time as Vault verifies service accounts exist in LDAP
// Uses a 60-second timeout to accommodate LDAP verification delays
func WriteLibrarySetWithRetry(t *testing.T, v *blackbox.Session, path string, data map[string]any) {
	t.Helper()

	// Log debug information about what we're trying to create
	t.Logf("DEBUG: Creating LDAP library set at path: %s", path)
	t.Logf("DEBUG: Library set configuration: %+v", data)

	// Extract service account names for debugging
	if serviceAccounts, ok := data["service_account_names"].([]string); ok {
		t.Logf("DEBUG: Vault will verify these service accounts exist in LDAP: %v", serviceAccounts)
		t.Logf("DEBUG: If this times out, manually verify accounts exist with the ldapsearch commands logged above")
	}

	v.EventuallyWithTimeout(func() error {
		_, err := v.Client.Logical().Write(path, data)
		if err != nil {
			t.Logf("DEBUG: Library set creation attempt failed: %v", err)
		}
		return err
	}, 60*time.Second)

	t.Logf("DEBUG: Successfully created LDAP library set at %s", path)
}

// checkLDAPUserExists verifies if a user exists in the LDAP directory
func checkLDAPUserExists(t *testing.T, username string) bool {
	ldapServer := os.Getenv("LDAP_SERVER")
	ldapPort := os.Getenv("LDAP_PORT")
	ldapUsername := os.Getenv("LDAP_USERNAME")
	ldapAdminPw := os.Getenv("LDAP_ADMIN_PW")

	cmd := exec.Command("ldapsearch",
		"-x",
		"-H", "ldap://"+ldapServer+":"+ldapPort,
		"-b", "ou=users,dc="+ldapUsername+",dc=com",
		"-D", "cn=admin,dc="+ldapUsername+",dc=com",
		"-w", ldapAdminPw,
		"(uid="+username+")",
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Logf("LDAP search command failed: %v", err)
		return false
	}

	return strings.Contains(string(output), "dn:")
}

// getLDIFPath returns the full path to an LDIF file in testdata
func getLDIFPath(filename string) string {
	_, currentFile, _, _ := runtime.Caller(0)
	dir := filepath.Dir(currentFile)
	return filepath.Join(dir, "testdata", filename)
}

// skipIfLDAPNotAvailable skips the test if LDAP configuration is not available
func skipIfLDAPNotAvailable(t *testing.T) {
	if os.Getenv("LDAP_SERVER") == "" {
		t.Skip("LDAP configuration not available - skipping LDAP test")
	}
}

// isCI returns true if running in a CI environment
func isCI() bool {
	return os.Getenv("CI") != "" ||
		os.Getenv("GITHUB_ACTIONS") != "" ||
		os.Getenv("ENOS_VAR_ci") != ""
}

// maskPassword masks a password for logging, showing only first 2 and last 2 characters
func maskPassword(password string) string {
	if len(password) <= 4 {
		return "****"
	}
	return password[:2] + "****" + password[len(password)-2:]
}

// waitForLDAP waits for the LDAP server to be ready by repeatedly attempting to connect
func waitForLDAP(t *testing.T, config *LDAPDomainConfig, timeout time.Duration) {
	t.Helper()

	deadline := time.Now().Add(timeout)
	attempt := 0

	for time.Now().Before(deadline) {
		attempt++
		cmd := exec.Command("ldapsearch",
			"-x",
			"-H", config.SetupURL, // Use public IP for connectivity checks from GitHub runner
			"-D", config.BindDN,
			"-w", config.BindPass,
			"-b", "",
			"-s", "base",
			"objectClass=*",
		)

		if output, err := cmd.CombinedOutput(); err == nil {
			t.Logf("LDAP server ready after %d attempts", attempt)
			return
		} else {
			if attempt == 1 || attempt%5 == 0 {
				t.Logf("LDAP connectivity check attempt %d failed, retrying... (output: %s)", attempt, string(output))
			}
		}

		time.Sleep(1 * time.Second)
	}

	t.Fatalf("LDAP server at %s did not become ready within %v", config.SetupURL, timeout)
}

// PrepareTestLDAPDomain creates an isolated LDAP domain for the test session.
// Uses the session's namespace to create a unique domain within the existing LDAP server.
//
// Parameters:
//   - t: The test instance
//   - session: The blackbox test session (provides unique namespace)
//   - requireIsolation: If true (CI mode), fails test if domain creation fails.
//     If false (dev mode), skips test if domain creation fails.
//
// Returns:
//   - cleanup: Function to delete the domain and all its contents
//   - config: LDAP config with domain-specific settings
//   - error: Non-nil if domain creation failed
func PrepareTestLDAPDomain(
	t *testing.T,
	session *blackbox.Session,
	requireIsolation bool,
) (cleanup func(), config *LDAPDomainConfig, err error) {
	t.Helper()

	// Get LDAP connection info from environment (set by Enos)
	// LDAP_SERVER: private IP for Vault operations (runs on Vault leader)
	// LDAP_SERVER_PUBLIC: public IP for setup operations (runs from GitHub runner)
	ldapServer := os.Getenv("LDAP_SERVER")
	ldapServerPublic := os.Getenv("LDAP_SERVER_PUBLIC")
	ldapPort := os.Getenv("LDAP_PORT")
	ldapBindDN := os.Getenv("LDAP_BIND_DN")
	ldapBindPass := os.Getenv("LDAP_BIND_PASS")

	t.Logf("LDAP Environment: SERVER=%s SERVER_PUBLIC=%s PORT=%s BIND_DN=%s BIND_PASS=%s",
		ldapServer, ldapServerPublic, ldapPort, ldapBindDN, maskPassword(ldapBindPass))

	if ldapServer == "" || ldapServerPublic == "" || ldapPort == "" || ldapBindDN == "" || ldapBindPass == "" {
		err := fmt.Errorf("LDAP environment variables not set")
		if requireIsolation {
			return nil, nil, fmt.Errorf("%w - required in CI", err)
		}
		return nil, nil, fmt.Errorf("%w - skipping in dev", err)
	}

	// Create unique organizational units under the existing dc=enos,dc=com domain
	// This avoids needing to create intermediate domain components
	domainName := session.Namespace // e.g., "bbsdk-a1b2c3d4"
	baseDN := "dc=enos,dc=com"      // Use existing base domain

	config = &LDAPDomainConfig{
		URL:      fmt.Sprintf("ldap://%s:%s", ldapServer, ldapPort), // Private IP for Vault
		BindDN:   ldapBindDN,
		BindPass: ldapBindPass,
		BaseDN:   baseDN,
		UserDN:   fmt.Sprintf("ou=%s-users,%s", domainName, baseDN),
		GroupDN:  fmt.Sprintf("ou=%s-groups,%s", domainName, baseDN),
	}

	// Store public IP for setup operations (ldapadd commands from GitHub runner)
	config.SetupURL = fmt.Sprintf("ldap://%s:%s", ldapServerPublic, ldapPort)

	// Create domain structure
	if err := createDomain(t, session, config); err != nil {
		if requireIsolation {
			return nil, nil, fmt.Errorf("failed to create LDAP domain in CI: %w", err)
		}
		return nil, nil, fmt.Errorf("failed to create LDAP domain: %w", err)
	}

	cleanup = func() {
		deleteDomain(t, config)
	}

	return cleanup, config, nil
}

// createDomain creates isolated organizational units for the test session
// Uses Eventually to retry the operation until LDAP server is ready
func createDomain(t *testing.T, session *blackbox.Session, config *LDAPDomainConfig) error {
	t.Helper()

	// Extract OU names from UserDN and GroupDN
	// e.g., "ou=bbsdk-xxx-users,dc=enos,dc=com" -> "bbsdk-xxx-users"
	userOU := strings.TrimPrefix(strings.Split(config.UserDN, ",")[0], "ou=")
	groupOU := strings.TrimPrefix(strings.Split(config.GroupDN, ",")[0], "ou=")

	// Create LDIF for isolated organizational units only
	// The base domain dc=enos,dc=com already exists
	ldif := fmt.Sprintf(`dn: %s
objectClass: top
objectClass: organizationalUnit
ou: %s

dn: %s
objectClass: top
objectClass: organizationalUnit
ou: %s
`, config.UserDN, userOU, config.GroupDN, groupOU)

	// Log the exact command for debugging (using SetupURL for public IP)
	t.Logf("DEBUG: Creating LDAP domain OUs with command:")
	t.Logf("  echo '%s' | ldapadd -x -H %s -D %s -w %s",
		strings.ReplaceAll(ldif, "\n", "\\n"), config.SetupURL, config.BindDN, config.BindPass)
	t.Logf("DEBUG: To verify OUs exist, run:")
	t.Logf("  ldapsearch -x -H %s -b \"dc=enos,dc=com\" -D %s -w %s \"(objectClass=organizationalUnit)\"",
		config.SetupURL, config.BindDN, config.BindPass)

	// Wait for LDAP server to be ready with extended timeout (30 seconds)
	// LDAP containers can take time to fully initialize, especially in CI environments
	t.Logf("Waiting for LDAP server to be ready at %s", config.SetupURL)

	waitForLDAP(t, config, 30*time.Second)

	// Use Eventually to retry ldapadd until LDAP server is ready
	session.Eventually(func() error {
		cmd := exec.Command("ldapadd",
			"-x",
			"-H", config.SetupURL, // Use public IP for setup operations from GitHub runner
			"-D", config.BindDN,
			"-w", config.BindPass,
		)
		cmd.Stdin = strings.NewReader(ldif)

		var stderr bytes.Buffer
		cmd.Stderr = &stderr

		if err := cmd.Run(); err != nil {
			t.Logf("ldapadd attempt failed: %v, stderr: %s", err, stderr.String())
			return fmt.Errorf("ldapadd failed: %w, stderr: %s", err, stderr.String())
		}
		return nil
	})

	t.Logf("Created isolated LDAP domain: %s", config.BaseDN)
	return nil
}

// deleteDomain recursively deletes the isolated organizational units
func deleteDomain(t *testing.T, config *LDAPDomainConfig) {
	t.Helper()

	// Delete both user and group OUs
	for _, dn := range []string{config.UserDN, config.GroupDN} {
		cmd := exec.Command("ldapdelete",
			"-x",
			"-r",                  // recursive
			"-H", config.SetupURL, // Use public IP for cleanup operations from GitHub runner
			"-D", config.BindDN,
			"-w", config.BindPass,
			dn,
		)

		var stderr bytes.Buffer
		cmd.Stderr = &stderr

		if err := cmd.Run(); err != nil {
			t.Logf("Warning: Failed to delete LDAP OU %s: %v, stderr: %s", dn, err, stderr.String())
		} else {
			t.Logf("Deleted isolated LDAP OU: %s", dn)
		}
	}
}

// CreateLDAPUser creates a user in the isolated domain
func CreateLDAPUser(t *testing.T, config *LDAPDomainConfig, username, password string) error {
	t.Helper()

	// Create LDIF for user
	userDN := fmt.Sprintf("uid=%s,%s", username, config.UserDN)
	ldif := fmt.Sprintf(`dn: %s
objectClass: top
objectClass: person
objectClass: organizationalPerson
objectClass: inetOrgPerson
uid: %s
cn: %s
sn: %s
userPassword: %s
`, userDN, username, username, username, password)

	// Log the exact command for debugging (using SetupURL for public IP)
	t.Logf("DEBUG: Creating LDAP user with command:")
	t.Logf("  echo '%s' | ldapadd -x -H %s -D %s -w %s",
		strings.ReplaceAll(ldif, "\n", "\\n"), config.SetupURL, config.BindDN, config.BindPass)
	t.Logf("DEBUG: To verify user exists, run:")
	t.Logf("  ldapsearch -x -H %s -b \"%s\" -D %s -w %s \"(uid=%s)\"",
		config.SetupURL, config.UserDN, config.BindDN, config.BindPass, username)
	t.Logf("DEBUG: To test bind as user, run:")
	t.Logf("  ldapsearch -x -H %s -b \"dc=enos,dc=com\" -D \"%s\" -w \"%s\" -s base",
		config.SetupURL, userDN, password)

	cmd := exec.Command("ldapadd",
		"-x",
		"-H", config.SetupURL, // Use public IP for setup operations from GitHub runner
		"-D", config.BindDN,
		"-w", config.BindPass,
	)
	cmd.Stdin = strings.NewReader(ldif)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create user %s: %w, stderr: %s", username, err, stderr.String())
	}

	t.Logf("Created LDAP user: %s", userDN)
	return nil
}

// CreateLDAPGroup creates a group in the isolated domain
func CreateLDAPGroup(t *testing.T, config *LDAPDomainConfig, groupName string, members []string) error {
	t.Helper()

	groupDN := fmt.Sprintf("cn=%s,%s", groupName, config.GroupDN)

	// Build member DNs
	var memberDNs []string
	for _, member := range members {
		memberDN := fmt.Sprintf("uid=%s,%s", member, config.UserDN)
		memberDNs = append(memberDNs, fmt.Sprintf("member: %s", memberDN))
	}

	// Create LDIF for group
	ldif := fmt.Sprintf(`dn: %s
objectClass: top
objectClass: groupOfNames
cn: %s
%s
`, groupDN, groupName, strings.Join(memberDNs, "\n"))

	cmd := exec.Command("ldapadd",
		"-x",
		"-H", config.URL,
		"-D", config.BindDN,
		"-w", config.BindPass,
	)
	cmd.Stdin = strings.NewReader(ldif)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create group %s: %w, stderr: %s", groupName, err, stderr.String())
	}

	t.Logf("Created LDAP group: %s with %d members", groupDN, len(members))
	return nil
}

// AddUserToGroup adds a user to a group in the isolated domain
func AddUserToGroup(t *testing.T, config *LDAPDomainConfig, username, groupName string) error {
	t.Helper()

	groupDN := fmt.Sprintf("cn=%s,%s", groupName, config.GroupDN)
	userDN := fmt.Sprintf("uid=%s,%s", username, config.UserDN)

	// Create LDIF for adding member
	ldif := fmt.Sprintf(`dn: %s
changetype: modify
add: member
member: %s
`, groupDN, userDN)

	cmd := exec.Command("ldapmodify",
		"-x",
		"-H", config.URL,
		"-D", config.BindDN,
		"-w", config.BindPass,
	)
	cmd.Stdin = strings.NewReader(ldif)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to add user %s to group %s: %w, stderr: %s", username, groupName, err, stderr.String())
	}

	t.Logf("Added user %s to group %s", username, groupName)
	return nil
}

// RemoveUserFromGroup removes a user from a group in the isolated domain
func RemoveUserFromGroup(t *testing.T, config *LDAPDomainConfig, username, groupName string) error {
	t.Helper()

	groupDN := fmt.Sprintf("cn=%s,%s", groupName, config.GroupDN)
	userDN := fmt.Sprintf("uid=%s,%s", username, config.UserDN)

	// Create LDIF for removing member
	ldif := fmt.Sprintf(`dn: %s
changetype: modify
delete: member
member: %s
`, groupDN, userDN)

	cmd := exec.Command("ldapmodify",
		"-x",
		"-H", config.URL,
		"-D", config.BindDN,
		"-w", config.BindPass,
	)
	cmd.Stdin = strings.NewReader(ldif)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to remove user %s from group %s: %w, stderr: %s", username, groupName, err, stderr.String())
	}

	t.Logf("Removed user %s from group %s", username, groupName)
	return nil
}

// CheckLDAPUserExistsInDomain verifies if a user exists in the isolated domain
func CheckLDAPUserExistsInDomain(t *testing.T, config *LDAPDomainConfig, username string) bool {
	t.Helper()

	cmd := exec.Command("ldapsearch",
		"-x",
		"-H", config.URL,
		"-b", config.UserDN,
		"-D", config.BindDN,
		"-w", config.BindPass,
		fmt.Sprintf("(uid=%s)", username),
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Logf("LDAP search command failed: %v", err)
		return false
	}

	return strings.Contains(string(output), "dn:")
}

// createDynamicRole creates a dynamic role with standard configuration
func createDynamicRole(v *blackbox.Session, mount, roleName string, config map[string]any) {
	defaultConfig := map[string]any{
		"username_template": "v-test-{{random 8}}",
		"default_ttl":       "1h",
		"max_ttl":           "24h",
	}

	// Merge provided config with defaults
	for k, v := range config {
		defaultConfig[k] = v
	}

	v.MustWrite(mount+"/role/"+roleName, defaultConfig)
}

// TODO: Add more helper functions as needed for common LDAP operations
// - verifyRoleExists(v, mount, roleName)
// - deleteRole(v, mount, roleName)
// - generateCredentials(v, mount, roleName)
// - parseAuditLog(logPath)
// - waitForLDAPOperation(timeout)
