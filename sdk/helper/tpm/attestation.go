// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: MPL-2.0

package tpm

import (
	"bytes"
	"context"
	"crypto"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/go-attestation/attest"
	"github.com/hashicorp/vault/api"
	"golang.org/x/crypto/blake2b"
)

// AttestationConfig contains the parameters for the TPM attestation workflow.
type AttestationConfig struct {
	// Namespace is the Vault namespace to use for attestation
	Namespace string

	// MountPath is the path to the TPM auth mount (e.g., "identity/tpm")
	MountPath string

	// RoleName is the role to attest to.
	RoleName string

	// CSRParams contains the certificate subject parameters
	CSRParams *CSRParams

	// TPMConfig identifies the TPM device to use for attestation.
	TPMConfig
}

// CSRParams contains parameters for generating a certificate signing request.
type CSRParams struct {
	CommonName string
	Country    string
	State      string
	Locality   string
	Org        string
	OrgUnit    string
}

// AttestationResult contains the output of a successful TPM attestation.
type AttestationResult struct {
	// CertPEM is the issued certificate in PEM format
	CertPEM []byte

	// CAChainPEM is the CA chain in PEM format
	CAChainPEM []byte

	// KeyMaterial contains the in-memory TPM key data
	KeyMaterial *TPMKeyMaterial

	// NotBefore is the certificate validity start time
	NotBefore time.Time

	// NotAfter is the certificate validity end time
	NotAfter time.Time
}

func (a *AttestationResult) Hash() []byte {
	var buf []byte
	buf = append(buf, a.CertPEM...)
	buf = append(buf, a.CAChainPEM...)
	buf = append(buf, []byte(a.KeyMaterial.PublicKeySHA256)...)
	buf = append(buf, a.KeyMaterial.AKBlob...)
	buf = append(buf, a.KeyMaterial.AppBlob...)

	result := make([]byte, 32)
	hash := blake2b.Sum256(buf)
	copy(result, hash[:])
	return result
}

// Attest performs the two-step TPM attestation workflow to obtain a certificate.
// It creates an attestation key (AK) and application key in the TPM, proves
// possession of the AK through credential activation, and obtains a certificate
// for the application key.
func Attest(ctx context.Context, client *api.Client, config *AttestationConfig) (*AttestationResult, error) {
	if client == nil {
		return nil, fmt.Errorf("client cannot be nil")
	}
	if config == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}
	if strings.TrimSpace(config.MountPath) == "" {
		return nil, fmt.Errorf("mount_path must be specified")
	}

	// Clone the client to avoid modifying the original
	attestClient, err := client.CloneWithHeaders()
	if err != nil {
		return nil, fmt.Errorf("failed to clone client: %w", err)
	}

	if config.Namespace != "" {
		attestClient.SetNamespace(config.Namespace)
	}

	tpm, openErr := openTPM(config.DevicePath)
	if openErr != nil {
		return nil, fmt.Errorf("failed to open TPM: %w", openErr)
	}
	defer tpm.Close()

	// Get endorsement keys
	eks, err := tpm.EKs()
	if err != nil {
		return nil, fmt.Errorf("failed to enumerate TPM EKs: %w", err)
	}
	if len(eks) == 0 {
		return nil, fmt.Errorf("no TPM EKs available")
	}

	ekPEM, err := publicKeyToPEM(eks[0].Public)
	if err != nil {
		return nil, fmt.Errorf("failed to encode EK public key: %w", err)
	}

	// Create attestation key
	ak, err := tpm.NewAK(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create TPM AK: %w", err)
	}
	defer ak.Close(tpm)

	akBlob, err := ak.Marshal()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal TPM AK: %w", err)
	}

	akParamsJSON, err := json.Marshal(akParamsFromAttest(ak.AttestationParameters()))
	if err != nil {
		return nil, fmt.Errorf("failed to encode TPM AK parameters: %w", err)
	}

	// Begin phase: send EK and AK parameters
	beginPath := fmt.Sprintf("%s/role/%s/gentpmcert/begin", config.MountPath, config.RoleName)
	beginSecret, err := attestClient.Logical().WriteWithContext(ctx, beginPath, map[string]interface{}{
		"ek_public_key": ekPEM,
		"ak_params":     base64.StdEncoding.EncodeToString(akParamsJSON),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to start TPM certificate issuance: %w", err)
	}
	if beginSecret == nil || beginSecret.Data == nil {
		return nil, fmt.Errorf("empty response from TPM certificate begin phase")
	}

	challengeID, _ := beginSecret.Data["challenge_id"].(string)
	encCredB64, _ := beginSecret.Data["encrypted_credentials"].(string)
	if strings.TrimSpace(challengeID) == "" || strings.TrimSpace(encCredB64) == "" {
		return nil, fmt.Errorf("incomplete response from TPM certificate begin phase")
	}

	// Decode and activate credential
	encCredBytes, err := base64.StdEncoding.DecodeString(encCredB64)
	if err != nil {
		return nil, fmt.Errorf("failed to decode encrypted credentials: %w", err)
	}

	var encCred attest.EncryptedCredential
	if err := json.Unmarshal(encCredBytes, &encCred); err != nil {
		return nil, fmt.Errorf("failed to decode encrypted credentials JSON: %w", err)
	}

	secret, err := ak.ActivateCredential(tpm, encCred)
	if err != nil {
		return nil, fmt.Errorf("failed to activate TPM credential: %w", err)
	}

	// Create application key
	appKey, err := tpm.NewKey(ak, &attest.KeyConfig{
		Algorithm: attest.ECDSA,
		Size:      256,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create TPM application key: %w", err)
	}
	defer appKey.Close()

	appBlob, err := appKey.Marshal()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal TPM application key: %w", err)
	}

	// Get signer for CSR
	appPubKey := appKey.Public()
	appPrivKey, err := appKey.Private(appPubKey)
	if err != nil {
		return nil, fmt.Errorf("failed to load TPM application private key: %w", err)
	}

	appSigner, ok := appPrivKey.(crypto.Signer)
	if !ok {
		return nil, fmt.Errorf("TPM application key does not implement crypto.Signer (got type %T)", appPrivKey)
	}

	// Build CSR subject
	hostname, err := os.Hostname()
	if err != nil {
		return nil, fmt.Errorf("failed to get hostname: %w", err)
	}

	csrSubject := pkix.Name{
		CommonName: hostname,
	}

	if config.CSRParams != nil {
		if strings.TrimSpace(config.CSRParams.CommonName) != "" {
			csrSubject.CommonName = config.CSRParams.CommonName
		}
		if strings.TrimSpace(config.CSRParams.Country) != "" {
			csrSubject.Country = []string{config.CSRParams.Country}
		}
		if strings.TrimSpace(config.CSRParams.State) != "" {
			csrSubject.Province = []string{config.CSRParams.State}
		}
		if strings.TrimSpace(config.CSRParams.Locality) != "" {
			csrSubject.Locality = []string{config.CSRParams.Locality}
		}
		if strings.TrimSpace(config.CSRParams.Org) != "" {
			csrSubject.Organization = []string{config.CSRParams.Org}
		}
		if strings.TrimSpace(config.CSRParams.OrgUnit) != "" {
			csrSubject.OrganizationalUnit = []string{config.CSRParams.OrgUnit}
		}
	}

	// Create CSR
	csrDER, err := x509.CreateCertificateRequest(rand.Reader, &x509.CertificateRequest{
		Subject: csrSubject,
	}, appSigner)
	if err != nil {
		return nil, fmt.Errorf("failed to create TPM CSR: %w", err)
	}

	csrPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE REQUEST",
		Bytes: csrDER,
	})

	keyCertificationJSON, err := json.Marshal(keyCertificationFromAttest(appKey.CertificationParameters()))
	if err != nil {
		return nil, fmt.Errorf("failed to encode TPM key certification: %w", err)
	}

	// Finish phase: send challenge response and CSR
	finishPath := fmt.Sprintf("%s/role/%s/gentpmcert/finish", config.MountPath, config.RoleName)
	finishSecret, err := attestClient.Logical().WriteWithContext(ctx, finishPath, map[string]interface{}{
		"ek_public_key":      ekPEM,
		"challenge_id":       challengeID,
		"challenge_response": base64.StdEncoding.EncodeToString(secret),
		"csr":                string(csrPEM),
		"key_certification":  base64.StdEncoding.EncodeToString(keyCertificationJSON),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to finish TPM certificate issuance: %w", err)
	}
	if finishSecret == nil || finishSecret.Data == nil {
		return nil, fmt.Errorf("empty response from TPM certificate finish phase")
	}

	certPEM, _ := finishSecret.Data["certificate"].(string)
	if strings.TrimSpace(certPEM) == "" {
		return nil, fmt.Errorf("certificate missing from TPM issuance response")
	}

	caChainPEM, err := extractCAChainPEM(finishSecret.Data["ca_chain"])
	if err != nil {
		return nil, fmt.Errorf("failed to decode ca_chain: %w", err)
	}

	notBefore, notAfter, publicKeySHA256, err := parseIssuedCertificate(certPEM)
	if err != nil {
		return nil, err
	}

	return &AttestationResult{
		CertPEM:    []byte(certPEM),
		CAChainPEM: caChainPEM,
		KeyMaterial: &TPMKeyMaterial{
			Type:            "tpm2",
			AKBlob:          akBlob,
			AppBlob:         appBlob,
			PublicKeySHA256: publicKeySHA256,
		},
		NotBefore: notBefore,
		NotAfter:  notAfter,
	}, nil
}

// publicKeyToPEM encodes a public key to PEM format.
func publicKeyToPEM(pub crypto.PublicKey) (string, error) {
	pubBytes, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		return "", err
	}

	return string(pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubBytes,
	})), nil
}

// extractCAChainPEM extracts and normalizes the CA chain from the API response.
func extractCAChainPEM(raw any) ([]byte, error) {
	switch chain := raw.(type) {
	case []interface{}:
		var out strings.Builder
		for _, item := range chain {
			s, ok := item.(string)
			if !ok {
				return nil, fmt.Errorf("ca_chain entry has invalid type %T", item)
			}
			out.WriteString(strings.TrimSpace(s))
			out.WriteString("\n")
		}
		return []byte(out.String()), nil
	case string:
		if strings.TrimSpace(chain) == "" {
			return nil, nil
		}
		if strings.HasSuffix(chain, "\n") {
			return []byte(chain), nil
		}
		return []byte(chain + "\n"), nil
	case nil:
		return nil, nil
	default:
		return nil, fmt.Errorf("unsupported ca_chain type %T", raw)
	}
}

// parseIssuedCertificate parses the certificate and extracts validity times and public key fingerprint.
func parseIssuedCertificate(certPEM string) (time.Time, time.Time, string, error) {
	block, _ := pem.Decode([]byte(certPEM))
	if block == nil || block.Type != "CERTIFICATE" {
		return time.Time{}, time.Time{}, "", fmt.Errorf("failed to parse issued certificate PEM")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return time.Time{}, time.Time{}, "", fmt.Errorf("failed to parse issued certificate: %w", err)
	}

	pubDER, err := x509.MarshalPKIXPublicKey(cert.PublicKey)
	if err != nil {
		return time.Time{}, time.Time{}, "", fmt.Errorf("failed to marshal issued certificate public key: %w", err)
	}
	sum := sha256.Sum256(pubDER)

	return cert.NotBefore, cert.NotAfter, hex.EncodeToString(sum[:]), nil
}

// SaveToDirectory writes the attestation result to a directory with fixed filenames.
// The directory must already exist.
func (r *AttestationResult) SaveToDirectory(destPath string, umask os.FileMode) error {
	if r == nil {
		return fmt.Errorf("attestation result is nil")
	}

	info, err := os.Stat(destPath)
	if err != nil {
		return fmt.Errorf("destination path error: %w", err)
	}
	if !info.IsDir() {
		return fmt.Errorf("destination path must be a directory")
	}

	// Create KeyReference for on-disk organization
	keyRef := &KeyReference{
		Type:            r.KeyMaterial.Type,
		Version:         1,
		PublicKeySHA256: r.KeyMaterial.PublicKeySHA256,
	}

	keyRefJSON, err := json.MarshalIndent(keyRef, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to encode key reference: %w", err)
	}
	keyRefJSON = append(keyRefJSON, '\n')

	privatePerm := applyUmask(0o600, umask)
	publicPerm := applyUmask(0o644, umask)

	files := []struct {
		path string
		data []byte
		perm os.FileMode
	}{
		{filepath.Join(destPath, "client-key.json"), keyRefJSON, privatePerm},
		{filepath.Join(destPath, "ak.blob"), r.KeyMaterial.AKBlob, privatePerm},
		{filepath.Join(destPath, "app.blob"), r.KeyMaterial.AppBlob, privatePerm},
		{filepath.Join(destPath, "ca_chain.pem"), r.CAChainPEM, publicPerm},
		{filepath.Join(destPath, "client.crt"), r.CertPEM, publicPerm},
	}

	for _, f := range files {
		if len(f.data) == 0 && strings.Contains(f.path, "ca_chain") {
			continue // Skip empty CA chain
		}
		if err := writeFileAtomic(f.path, f.data, f.perm); err != nil {
			return fmt.Errorf("failed to write %s: %w", filepath.Base(f.path), err)
		}
	}

	return nil
}

// LoadFromDirectory loads an attestation result from a directory.
func LoadFromDirectory(destPath string) (*AttestationResult, error) {
	certPath := filepath.Join(destPath, "client.crt")
	keyRefPath := filepath.Join(destPath, "client-key.json")
	akPath := filepath.Join(destPath, "ak.blob")
	appPath := filepath.Join(destPath, "app.blob")
	caChainPath := filepath.Join(destPath, "ca_chain.pem")

	certPEM, err := os.ReadFile(certPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read certificate: %w", err)
	}

	keyRefBytes, err := os.ReadFile(keyRefPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read key reference: %w", err)
	}

	var keyRef KeyReference
	if err := json.Unmarshal(keyRefBytes, &keyRef); err != nil {
		return nil, fmt.Errorf("failed to decode TPM key reference %q: %w", keyRefPath, err)
	}

	if keyRef.Type != "tpm2" {
		return nil, fmt.Errorf("unsupported TPM key reference type %q", keyRef.Type)
	}

	akBlob, err := os.ReadFile(akPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read AK blob from %q: %w", akPath, err)
	}

	appBlob, err := os.ReadFile(appPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read app blob from %q: %w", appPath, err)
	}

	keyMaterial := &TPMKeyMaterial{
		Type:            keyRef.Type,
		AKBlob:          akBlob,
		AppBlob:         appBlob,
		PublicKeySHA256: keyRef.PublicKeySHA256,
	}

	caChainPEM, err := os.ReadFile(caChainPath)
	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("failed to read CA chain: %w", err)
	}

	notBefore, notAfter, _, err := parseIssuedCertificate(string(certPEM))
	if err != nil {
		return nil, err
	}

	return &AttestationResult{
		CertPEM:     certPEM,
		CAChainPEM:  caChainPEM,
		KeyMaterial: keyMaterial,
		NotBefore:   notBefore,
		NotAfter:    notAfter,
	}, nil
}

// applyUmask applies a umask to a base permission.
func applyUmask(basePerm, umask os.FileMode) os.FileMode {
	return basePerm &^ umask
}

// writeFileAtomic writes data to a file atomically using a temp file and rename.
func writeFileAtomic(targetPath string, data []byte, perm os.FileMode) error {
	dir := filepath.Dir(targetPath)
	tmpFile, err := os.CreateTemp(dir, ".tmp-tpm-*")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	tmpPath := tmpFile.Name()

	defer func() {
		if tmpFile != nil {
			tmpFile.Close()
			os.Remove(tmpPath)
		}
	}()

	if _, err := tmpFile.Write(data); err != nil {
		return fmt.Errorf("failed to write to temp file: %w", err)
	}
	if err := tmpFile.Sync(); err != nil {
		return fmt.Errorf("failed to sync temp file: %w", err)
	}
	if err := tmpFile.Close(); err != nil {
		return fmt.Errorf("failed to close temp file: %w", err)
	}
	tmpFile = nil

	if err := os.Chmod(tmpPath, perm); err != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("failed to set permissions on temp file: %w", err)
	}
	if err := os.Rename(tmpPath, targetPath); err != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("failed to rename temp file to target: %w", err)
	}

	return nil
}

func (result *AttestationResult) Validate(config TPMConfig) error {
	tpm, err := openTPM(config.DevicePath)
	if err != nil {
		return fmt.Errorf("failed to open TPM: %w", err)
	}
	defer tpm.Close()

	// Load the AK and app key to verify they're valid and match the certificate
	ak, err := tpm.LoadAK(result.KeyMaterial.AKBlob)
	if err != nil {
		return fmt.Errorf("failed to load persisted AK: %w", err)
	}
	defer ak.Close(tpm)

	appKey, err := tpm.LoadKey(result.KeyMaterial.AppBlob)
	if err != nil {
		return fmt.Errorf("failed to load persisted app key: %w", err)
	}
	defer appKey.Close()

	// Compute the SHA256 of the TPM app key's public key
	appPubDER, err := x509.MarshalPKIXPublicKey(appKey.Public())
	if err != nil {
		return fmt.Errorf("failed to marshal TPM application public key: %w", err)
	}
	appPubSum := sha256.Sum256(appPubDER)
	appPubSHA := hex.EncodeToString(appPubSum[:])

	// Validate the key material's fingerprint matches (if present)
	if result.KeyMaterial.PublicKeySHA256 != "" && !strings.EqualFold(result.KeyMaterial.PublicKeySHA256, appPubSHA) {
		return fmt.Errorf("persisted TPM application key fingerprint mismatch")
	}

	// Parse the certificate to get its public key
	block, _ := pem.Decode(result.CertPEM)
	if block == nil || block.Type != "CERTIFICATE" {
		return fmt.Errorf("failed to parse persisted certificate PEM")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return fmt.Errorf("failed to parse persisted certificate: %w", err)
	}

	// Validate the certificate's public key matches the TPM app key
	certPubDER, err := x509.MarshalPKIXPublicKey(cert.PublicKey)
	if err != nil {
		return fmt.Errorf("failed to marshal certificate public key: %w", err)
	}
	if !bytes.Equal(appPubDER, certPubDER) {
		return fmt.Errorf("persisted certificate public key does not match TPM application key")
	}

	return nil
}

// EndorsementKeys returns the endorsement keys for the TPM identified by cfg.
func EndorsementKeys(cfg TPMConfig) ([]attest.EK, error) {
	tpm, err := openTPM(cfg.DevicePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open TPM: %w", err)
	}
	defer tpm.Close()
	return tpm.EKs()
}
