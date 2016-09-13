package command

import (
	"bytes"
	"encoding/base64"
	"os"
	"reflect"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/vault/helper/pgpkeys"
	"github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/meta"
	"github.com/hashicorp/vault/vault"
	"github.com/keybase/go-crypto/openpgp"
	"github.com/keybase/go-crypto/openpgp/packet"
	"github.com/mitchellh/cli"
)

func TestInit(t *testing.T) {
	ui := new(cli.MockUi)
	c := &InitCommand{
		Meta: meta.Meta{
			Ui: ui,
		},
	}

	core := vault.TestCore(t)
	ln, addr := http.TestServer(t, core)
	defer ln.Close()

	init, err := core.Initialized()
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if init {
		t.Fatal("should not be initialized")
	}

	args := []string{"-address", addr}
	if code := c.Run(args); code != 0 {
		t.Fatalf("bad: %d\n\n%s", code, ui.ErrorWriter.String())
	}

	init, err = core.Initialized()
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if !init {
		t.Fatal("should be initialized")
	}

	sealConf, err := core.SealAccess().BarrierConfig()
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	expected := &vault.SealConfig{
		Type:            "shamir",
		SecretShares:    5,
		SecretThreshold: 3,
	}
	if !reflect.DeepEqual(expected, sealConf) {
		t.Fatalf("expected:\n%#v\ngot:\n%#v\n", expected, sealConf)
	}
}

func TestInit_Check(t *testing.T) {
	ui := new(cli.MockUi)
	c := &InitCommand{
		Meta: meta.Meta{
			Ui: ui,
		},
	}

	core := vault.TestCore(t)
	ln, addr := http.TestServer(t, core)
	defer ln.Close()

	// Should return 2, not initialized
	args := []string{"-address", addr, "-check"}
	if code := c.Run(args); code != 2 {
		t.Fatalf("bad: %d\n\n%s", code, ui.ErrorWriter.String())
	}

	// Now initialize it
	args = []string{"-address", addr}
	if code := c.Run(args); code != 0 {
		t.Fatalf("bad: %d\n\n%s", code, ui.ErrorWriter.String())
	}

	// Should return 0, initialized
	args = []string{"-address", addr, "-check"}
	if code := c.Run(args); code != 0 {
		t.Fatalf("bad: %d\n\n%s", code, ui.ErrorWriter.String())
	}

	init, err := core.Initialized()
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if !init {
		t.Fatal("should be initialized")
	}
}

func TestInit_custom(t *testing.T) {
	ui := new(cli.MockUi)
	c := &InitCommand{
		Meta: meta.Meta{
			Ui: ui,
		},
	}

	core := vault.TestCore(t)
	ln, addr := http.TestServer(t, core)
	defer ln.Close()

	init, err := core.Initialized()
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if init {
		t.Fatal("should not be initialized")
	}

	args := []string{
		"-address", addr,
		"-key-shares", "7",
		"-key-threshold", "3",
	}
	if code := c.Run(args); code != 0 {
		t.Fatalf("bad: %d\n\n%s", code, ui.ErrorWriter.String())
	}

	init, err = core.Initialized()
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if !init {
		t.Fatal("should be initialized")
	}

	sealConf, err := core.SealAccess().BarrierConfig()
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	expected := &vault.SealConfig{
		Type:            "shamir",
		SecretShares:    7,
		SecretThreshold: 3,
	}
	if !reflect.DeepEqual(expected, sealConf) {
		t.Fatalf("expected:\n%#v\ngot:\n%#v\n", expected, sealConf)
	}

	re, err := regexp.Compile("\\s+Initial Root Token:\\s+(.*)")
	if err != nil {
		t.Fatalf("Error compiling regex: %s", err)
	}
	matches := re.FindAllStringSubmatch(ui.OutputWriter.String(), -1)
	if len(matches) != 1 {
		t.Fatalf("Unexpected number of tokens found, got %d", len(matches))
	}

	rootToken := matches[0][1]

	client, err := c.Client()
	if err != nil {
		t.Fatalf("Error fetching client: %v", err)
	}

	client.SetToken(rootToken)

	re, err = regexp.Compile("\\s*Unseal Key \\d+: (.*)")
	if err != nil {
		t.Fatalf("Error compiling regex: %s", err)
	}
	matches = re.FindAllStringSubmatch(ui.OutputWriter.String(), -1)
	if len(matches) != 7 {
		t.Fatalf("Unexpected number of keys returned, got %d, matches was \n\n%#v\n\n, input was \n\n%s\n\n", len(matches), matches, ui.OutputWriter.String())
	}

	var unsealed bool
	for i := 0; i < 3; i++ {
		decodedKey, err := base64.StdEncoding.DecodeString(strings.TrimSpace(matches[i][1]))
		if err != nil {
			t.Fatalf("err decoding key %v: %v", matches[i][1], err)
		}
		unsealed, err = core.Unseal(decodedKey)
		if err != nil {
			t.Fatalf("err during unseal: %v; key was %v", err, matches[i][1])
		}
	}
	if !unsealed {
		t.Fatal("expected to be unsealed")
	}

	tokenInfo, err := client.Auth().Token().LookupSelf()
	if err != nil {
		t.Fatalf("Error looking up root token info: %v", err)
	}

	if tokenInfo.Data["policies"].([]interface{})[0].(string) != "root" {
		t.Fatalf("expected root policy")
	}
}

func TestInit_PGP(t *testing.T) {
	ui := new(cli.MockUi)
	c := &InitCommand{
		Meta: meta.Meta{
			Ui: ui,
		},
	}

	core := vault.TestCore(t)
	ln, addr := http.TestServer(t, core)
	defer ln.Close()

	init, err := core.Initialized()
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if init {
		t.Fatal("should not be initialized")
	}

	tempDir, pubFiles, err := getPubKeyFiles(t)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	args := []string{
		"-address", addr,
		"-key-shares", "2",
		"-pgp-keys", pubFiles[0] + ",@" + pubFiles[1] + "," + pubFiles[2],
		"-key-threshold", "2",
		"-root-token-pgp-key", pubFiles[0],
	}

	// This should fail, as key-shares does not match pgp-keys size
	if code := c.Run(args); code == 0 {
		t.Fatalf("bad (command should have failed): %d\n\n%s", code, ui.ErrorWriter.String())
	}

	args = []string{
		"-address", addr,
		"-key-shares", "4",
		"-pgp-keys", pubFiles[0] + ",@" + pubFiles[1] + "," + pubFiles[2] + "," + pubFiles[3],
		"-key-threshold", "2",
		"-root-token-pgp-key", pubFiles[0],
	}

	ui.OutputWriter.Reset()

	if code := c.Run(args); code != 0 {
		t.Fatalf("bad: %d\n\n%s", code, ui.ErrorWriter.String())
	}

	init, err = core.Initialized()
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if !init {
		t.Fatal("should be initialized")
	}

	sealConf, err := core.SealAccess().BarrierConfig()
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	pgpKeys := []string{}
	for _, pubFile := range pubFiles {
		pub, err := pgpkeys.ReadPGPFile(pubFile)
		if err != nil {
			t.Fatalf("bad: %v", err)
		}
		pgpKeys = append(pgpKeys, pub)
	}

	expected := &vault.SealConfig{
		Type:            "shamir",
		SecretShares:    4,
		SecretThreshold: 2,
		PGPKeys:         pgpKeys,
	}
	if !reflect.DeepEqual(expected, sealConf) {
		t.Fatalf("expected:\n%#v\ngot:\n%#v\n", expected, sealConf)
	}

	re, err := regexp.Compile("\\s+Initial Root Token:\\s+(.*)")
	if err != nil {
		t.Fatalf("Error compiling regex: %s", err)
	}
	matches := re.FindAllStringSubmatch(ui.OutputWriter.String(), -1)
	if len(matches) != 1 {
		t.Fatalf("Unexpected number of tokens found, got %d", len(matches))
	}

	encRootToken := matches[0][1]
	privKeyBytes, err := base64.StdEncoding.DecodeString(pgpkeys.TestPrivKey1)
	if err != nil {
		t.Fatalf("error decoding private key: %v", err)
	}
	ptBuf := bytes.NewBuffer(nil)
	entity, err := openpgp.ReadEntity(packet.NewReader(bytes.NewBuffer(privKeyBytes)))
	if err != nil {
		t.Fatalf("Error parsing private key: %s", err)
	}
	var rootBytes []byte
	rootBytes, err = base64.StdEncoding.DecodeString(encRootToken)
	if err != nil {
		t.Fatalf("Error decoding root token: %s", err)
	}
	entityList := &openpgp.EntityList{entity}
	md, err := openpgp.ReadMessage(bytes.NewBuffer(rootBytes), entityList, nil, nil)
	if err != nil {
		t.Fatalf("Error decrypting root token: %s", err)
	}
	ptBuf.ReadFrom(md.UnverifiedBody)
	rootToken := ptBuf.String()

	parseDecryptAndTestUnsealKeys(t, ui.OutputWriter.String(), rootToken, false, nil, nil, core)

	client, err := c.Client()
	if err != nil {
		t.Fatalf("Error fetching client: %v", err)
	}

	client.SetToken(rootToken)

	tokenInfo, err := client.Auth().Token().LookupSelf()
	if err != nil {
		t.Fatalf("Error looking up root token info: %v", err)
	}

	if tokenInfo.Data["policies"].([]interface{})[0].(string) != "root" {
		t.Fatalf("expected root policy")
	}
}
