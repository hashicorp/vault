package chefnode

import (
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"testing"

	"strings"

	"fmt"

	"bytes"
	"encoding/json"

	"crypto/rsa"

	"github.com/hashicorp/vault/helper/policyutil"
	"github.com/hashicorp/vault/logical"
	logicaltest "github.com/hashicorp/vault/logical/testing"
)

func TestBackend_Config(t *testing.T) {
	config := logical.TestBackendConfig()
	storage := &logical.InmemStorage{}
	config.StorageView = storage

	b, err := Backend().Setup(config)
	if err != nil {
		t.Fatal(err)
	}

	// Valid Case
	data := map[string]interface{}{
		"client_name": "vault",
		"base_url":    "https://api.chef.io/organizations/test",
		"client_key": `-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEA0ZsplirBTD4+895tiXD8HjCbMI643sPRVz+c92bbntT1OneB
KyAX4YmKlCNjZpEiTrrk7sWMXFKtN0PuAH7ZGSdUtKl4nIGs4WS2L0Axhg8mBvvw
mn+tj9ksf09N9TfNipY/l/flKUvddL+bdsC1OIms1Bv10nS8Bg0d1kGfms7rMose
ea4Tk/XO9d5svNI7i/GZtsvYyQgQvdvbzV5uf2mIhhMQHokRdxII+qR/QYV5qm9C
qgVFSwBeU3uem8Nvi+OyDI19w21esUx5+Dakx9k2nuDqEInPE8yNfzkGDOFLpggZ
/h0XO7KZmyzCLBe99qfjIjaFAux3BZkVQ9km6wIDAQABAoIBADeUb0iUecEf2E2O
M3l4bkILHXuYvMjFH+OEyLiJm77YNVaVjbjDv9FcSVTStW7jGTfLMx1lYLyyZ5/5
8UhMWoDi/wEQ1xyY/iCeNfj9iqRDrA+6CqjNJla4faYcf02AyI3xHVfMsgVrSoPE
sxKgMu2VBDESYPK3ZYwtOjYwHIRN0xJyLbE8WoFi1Z6sFeRKhy+oM7Nf2lNxskMr
U/wWJVqmQxdqa8I4ldq2LoUq0hfWJJ7uD4uXvYmk6nSNDtII0ePj0sAyC6cm/gv+
wdiLilQGr8jOzTKDlB4TON/1zUivadD4Mrr1v1EqRUD3IbSAOMbpd64co2BbhnQh
iV5sfhECgYEA6IsyyIo1R94uLP9IB3PxWIgE7ENMai2ltmhO0hxDoPqRBSgvFdBe
WBQjNjSmUjFGswt3tlOdsoNupByug3fnGO1DC0Cc1q2F05HyYOtxJiPLmWtf1v5W
MYx8R0IaPkwIZHlQ6fWuhjhUt0nfRRPlAYzLGeCTCFdmgx+IRjCmVM0CgYEA5r+p
pncWVfO/z25i8hwtd7XE+LRpC4FeUiF/JN1qGEDwtKmOua89qt8tB7KdwYuDVBWS
uGsoH7U8o1vm1gjKOFJabtFKXkuidzD1bTwRYJhuE8QReGchdQXOVKDmXuOSxlkG
6u0c834zvxd2SM3ObTshmhfTm1r1mpLGc+shqpcCgYAcxPvna5Hj7kzwLDURFvsI
5OsW/8x4ZmVWB9mYjP6g7975MFuC62CArR0eG61oBcilZgnNeNLNvwz1KMc+ZJsm
rlPZFIlS1ez0m93Mt9Qrz8nklTAqPRUU1Dib6EWu52EybP/hsg+Bc36nnnAM78Up
R+3oqawHICkCl+gYJvStEQKBgDuzQvVqwkCiu/GzIa56U9kxEjE2nCb55aliOT8U
eiqkQqK7a83m5RGchE4FjINS0TukCT3lm4/4mCO711FxHMDNrdAWHiOfdf1YkWcd
r3FKftBmXg7EwAdC5UtIBdJvFr5ysjN9/YuSD1lVfKkBdnMUZXE00O7U7c58QxQi
tacpAoGBANsD/YDhF5zsq8vuTtT2GkrjI8T7MJDjqOH3jAN4ZclXP5UqcC1k3ot6
SLckarscm0jQSlInbYwd+PzOeQY5bc8UrFWcsX3G0PSHEDcNbC9K64LWidlIs01S
AhXEa2Ie+mNe5fKzBIhREz9cQo6kF+3lYeL3XeKZiWMMgEsFlmFZ
-----END RSA PRIVATE KEY-----
`,
	}

	_, err = b.HandleRequest(&logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "config",
		Data:      data,
		Storage:   storage,
	})
	if err != nil {
		t.Fatal(err)
	}

	resp, err := b.HandleRequest(&logical.Request{
		Operation: logical.ReadOperation,
		Path:      "config",
		Storage:   storage,
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil || resp.IsError() {
		t.Fatal("Couldn't read config data")
	}
	if resp.Data["client_name"].(string) != data["client_name"].(string) {
		t.Fatal("Couldn't read client name from config")
	}
	if resp.Data["base_url"].(string) != data["base_url"].(string) {
		t.Fatal("Couldn't read base url from config")
	}
	if resp.Data["client_key"].(string) != data["client_key"].(string) {
		t.Fatal("Couldn't read client key from config")
	}

	// Bad Key
	data2 := map[string]interface{}{
		"client_name": "vault",
		"base_url":    "https://api.chef.io/organizations/test",
		"client_key": `-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEA0ZsplirBTD4+895tiXD8HjCbMI643sPRVz+c92bbntT1OneB
KyAX4YmKlCNjZpEiTrrk7sWMXFKtN0PuAH7ZGSdUtKl4nIGs4WS2L0Axhg8mBvvw
mn+tj9ksf09N9TfNipY/l/flKUvddL+bdsC1OIms1Bv10nS8Bg0d1kGfms7rMose
ea4Tk/XO9d5svNI7i/GZtsvYyQgQvdvbzV5uf2mIhhMQHokRdxII+qR/QYV5qm9C
qgVFSwBeU3uem8Nvi+OyDI19w21esUx5+Dakx9k2nuDqEInPE8yNfzkGDOFLpggZ
/h0XO7KZmyzCLBe99qfjIjaFAux3BZkVQ9km6wIDAQABAoIBADeUb0iUecEf2E2O
M3l4bkILHXuYvMjFH+OEyLiJm77YNVaVjbjDv9FcSVTStW7jGTfLMx1lYLyyZ5/5
8UhMWoDi/wEQ1xyY/iCeNfj9iqRDrA+6CqjNJla4faYcf02AyI3xHVfMsgVrSoPE
sxKgMu2VBDESYPK3ZYwtOjYwHIRN0xJyLbE8WoFi1Z6sFeRKhy+oM7Nf2lNxskMr
U/wWJVqmQxdqa8I4ldq2LoUq0hfWJJ7uD4uXvYmk6nSNDtII0ePj0sAyC6cm/gv+
wdiLilQGr8jOzTKDlB4TON/1zUivadD4Mrr1v1EqRUD3IbSAOMbpd64co2BbhnQh
iV5sfhECgYEA6IsyyIo1R94uLP9IB3PxWIgE7ENMai2ltmhO0hxDoPqRBSgvFdBe
5OsW/8x4ZmVWB9mYjP6g7975MFuC62CArR0eG61oBcilZgnNeNLNvwz1KMc+ZJsm
rlPZFIlS1ez0m93Mt9Qrz8nklTAqPRUU1Dib6EWu52EybP/hsg+Bc36nnnAM78Up
R+3oqawHICkCl+gYJvStEQKBgDuzQvVqwkCiu/GzIa56U9kxEjE2nCb55aliOT8U
eiqkQqK7a83m5RGchE4FjINS0TukCT3lm4/4mCO711FxHMDNrdAWHiOfdf1YkWcd
r3FKftBmXg7EwAdC5UtIBdJvFr5ysjN9/YuSD1lVfKkBdnMUZXE00O7U7c58QxQi
tacpAoGBANsD/YDhF5zsq8vuTtT2GkrjI8T7MJDjqOH3jAN4ZclXP5UqcC1k3ot6
SLckarscm0jQSlInbYwd+PzOeQY5bc8UrFWcsX3G0PSHEDcNbC9K64LWidlIs01S
AhXEa2Ie+mNe5fKzBIhREz9cQo6kF+3lYeL3XeKZiWMMgEsFlmFZ
-----END RSA PRIVATE KEY-----
`,
	}

	_, err = b.HandleRequest(&logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "config",
		Data:      data2,
		Storage:   storage,
	})
	if err == nil {
		t.Fatal("Config accepted bad key")
	}

	// Bad URL
	data3 := map[string]interface{}{
		"client_name": "vault",
		"base_url":    "a0s9duflnc;asd",
		"client_key": `-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEA0ZsplirBTD4+895tiXD8HjCbMI643sPRVz+c92bbntT1OneB
KyAX4YmKlCNjZpEiTrrk7sWMXFKtN0PuAH7ZGSdUtKl4nIGs4WS2L0Axhg8mBvvw
mn+tj9ksf09N9TfNipY/l/flKUvddL+bdsC1OIms1Bv10nS8Bg0d1kGfms7rMose
ea4Tk/XO9d5svNI7i/GZtsvYyQgQvdvbzV5uf2mIhhMQHokRdxII+qR/QYV5qm9C
qgVFSwBeU3uem8Nvi+OyDI19w21esUx5+Dakx9k2nuDqEInPE8yNfzkGDOFLpggZ
/h0XO7KZmyzCLBe99qfjIjaFAux3BZkVQ9km6wIDAQABAoIBADeUb0iUecEf2E2O
M3l4bkILHXuYvMjFH+OEyLiJm77YNVaVjbjDv9FcSVTStW7jGTfLMx1lYLyyZ5/5
8UhMWoDi/wEQ1xyY/iCeNfj9iqRDrA+6CqjNJla4faYcf02AyI3xHVfMsgVrSoPE
sxKgMu2VBDESYPK3ZYwtOjYwHIRN0xJyLbE8WoFi1Z6sFeRKhy+oM7Nf2lNxskMr
U/wWJVqmQxdqa8I4ldq2LoUq0hfWJJ7uD4uXvYmk6nSNDtII0ePj0sAyC6cm/gv+
wdiLilQGr8jOzTKDlB4TON/1zUivadD4Mrr1v1EqRUD3IbSAOMbpd64co2BbhnQh
iV5sfhECgYEA6IsyyIo1R94uLP9IB3PxWIgE7ENMai2ltmhO0hxDoPqRBSgvFdBe
WBQjNjSmUjFGswt3tlOdsoNupByug3fnGO1DC0Cc1q2F05HyYOtxJiPLmWtf1v5W
MYx8R0IaPkwIZHlQ6fWuhjhUt0nfRRPlAYzLGeCTCFdmgx+IRjCmVM0CgYEA5r+p
pncWVfO/z25i8hwtd7XE+LRpC4FeUiF/JN1qGEDwtKmOua89qt8tB7KdwYuDVBWS
uGsoH7U8o1vm1gjKOFJabtFKXkuidzD1bTwRYJhuE8QReGchdQXOVKDmXuOSxlkG
6u0c834zvxd2SM3ObTshmhfTm1r1mpLGc+shqpcCgYAcxPvna5Hj7kzwLDURFvsI
5OsW/8x4ZmVWB9mYjP6g7975MFuC62CArR0eG61oBcilZgnNeNLNvwz1KMc+ZJsm
rlPZFIlS1ez0m93Mt9Qrz8nklTAqPRUU1Dib6EWu52EybP/hsg+Bc36nnnAM78Up
R+3oqawHICkCl+gYJvStEQKBgDuzQvVqwkCiu/GzIa56U9kxEjE2nCb55aliOT8U
eiqkQqK7a83m5RGchE4FjINS0TukCT3lm4/4mCO711FxHMDNrdAWHiOfdf1YkWcd
r3FKftBmXg7EwAdC5UtIBdJvFr5ysjN9/YuSD1lVfKkBdnMUZXE00O7U7c58QxQi
tacpAoGBANsD/YDhF5zsq8vuTtT2GkrjI8T7MJDjqOH3jAN4ZclXP5UqcC1k3ot6
SLckarscm0jQSlInbYwd+PzOeQY5bc8UrFWcsX3G0PSHEDcNbC9K64LWidlIs01S
AhXEa2Ie+mNe5fKzBIhREz9cQo6kF+3lYeL3XeKZiWMMgEsFlmFZ
-----END RSA PRIVATE KEY-----
`,
	}

	_, err = b.HandleRequest(&logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "config",
		Data:      data3,
		Storage:   storage,
	})
	if err == nil {
		t.Fatal("Config accepted bad URL")
	}
}

func TestBackend_Environment(t *testing.T) {
	config := logical.TestBackendConfig()
	storage := &logical.InmemStorage{}
	config.StorageView = storage
	b, err := Backend().Setup(config)
	if err != nil {
		t.Fatal(err)
	}

	data := map[string]interface{}{
		"policies": "pol1,pol2",
	}

	resp, err := b.HandleRequest(&logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "environment/env1",
		Data:      data,
		Storage:   storage,
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp != nil && resp.IsError() {
		t.Fatal("Failed to create environment")
	}

	resp, err = b.HandleRequest(&logical.Request{
		Operation: logical.ReadOperation,
		Path:      "environment/env1",
		Storage:   storage,
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil || resp.IsError() {
		t.Fatal("Failed to read environment")
	}
	if !policyutil.EquivalentPolicies(strings.Split(data["policies"].(string), ","), resp.Data["policies"].([]string)) {
		t.Fatalf("policies didn't match: expected: %#v\ngot: %#v\n", data, resp.Data)
	}

}

func TestBackend_Tag(t *testing.T) {
	config := logical.TestBackendConfig()
	storage := &logical.InmemStorage{}
	config.StorageView = storage
	b, err := Backend().Setup(config)
	if err != nil {
		t.Fatal(err)
	}

	data := map[string]interface{}{
		"policies": "pol1,pol2",
	}

	resp, err := b.HandleRequest(&logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "tag/tag1",
		Data:      data,
		Storage:   storage,
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp != nil && resp.IsError() {
		t.Fatal("Failed to create tag")
	}

	resp, err = b.HandleRequest(&logical.Request{
		Operation: logical.ReadOperation,
		Path:      "tag/tag1",
		Storage:   storage,
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil || resp.IsError() {
		t.Fatal("Failed to read tag")
	}
	if !policyutil.EquivalentPolicies(strings.Split(data["policies"].(string), ","), resp.Data["policies"].([]string)) {
		t.Fatalf("policies didn't match: expected: %#v\ngot: %#v\n", data, resp.Data)
	}

}

func TestBackend_Role(t *testing.T) {
	config := logical.TestBackendConfig()
	storage := &logical.InmemStorage{}
	config.StorageView = storage
	b, err := Backend().Setup(config)
	if err != nil {
		t.Fatal(err)
	}

	data := map[string]interface{}{
		"policies": "pol1,pol2",
	}

	resp, err := b.HandleRequest(&logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "role/role1",
		Data:      data,
		Storage:   storage,
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp != nil && resp.IsError() {
		t.Fatal("Failed to create role")
	}

	resp, err = b.HandleRequest(&logical.Request{
		Operation: logical.ReadOperation,
		Path:      "role/role1",
		Storage:   storage,
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil || resp.IsError() {
		t.Fatal("Failed to read role")
	}
	if !policyutil.EquivalentPolicies(strings.Split(data["policies"].(string), ","), resp.Data["policies"].([]string)) {
		t.Fatalf("policies didn't match: expected: %#v\ngot: %#v\n", data, resp.Data)
	}

}

func TestBackend_Authenticate(t *testing.T) {
	privKey := `-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEAulLSQ7sdr06pntGFWloywCFfaZGcxZ5HoRtpyAlpwmlzkK4R
YoUrWbTc4fOy81dBpsiHwK3ZzefaPhMckhF/oset0YOIhDRqcM3FflmrFofqlGgu
dTiJ0suH8s6IcSEtzo/wzq5T/KIfJHwVEhHtZGAnwjc2YQzsMz+KTeA5RSt8etNX
L+mWuXLY5HB5A2EIXIiST6DaBPSH49TQn1pjcHjXXeglmkwtjpSd+x6biA6YU0WY
lgH2/aND7e+Pgtox99NXCIp6cw4ne+wdBJOyeOtCYnBQvvF5n0+Jbgcc9Ox9D4fv
TGZzRz83+bqkLRjgAD9CfSh6Ah1hIG+tmCEQywIDAQABAoIBAA2jVEqq3pBfZKEA
Ww9y/LX9e1th0iTQ4hNTy1ld/wTA7TmQ1Cru7m5hg61yRg3zvBV2JiGfWArvRpU2
lufGKh6DGSD1zL9IiuX42dTWwWQjzLLSMVxZKBVq2meWYHxPXmf5NzoZnoImZ7sm
7e/lqgen1iEsI2nVJVDW/MuYdviuBzCO01ZwMmQ0vqMaty+Ed2mK3vSaBXOvZJtp
Vfuw66TVX/g7f3q3l0BgGf90TMRd00eO3WAB+IpePlcBmTvIEihmS/ETrmODflsk
Nnna8MA5UKaCes6i7qMQqFJ/5qQzdjNg40BmfpHYG9izID0LQubV3fLgCxPTsZwW
WSBd1SECgYEA4WC/9tU5b9+SO7Ib4bqgfvrIImx62SA+uBZscTwaubWlZxOazIjE
RzuaEdRH5QnQOxcXtpscCZ+UkNLITKp2bD36D6M46SlbMLLzd9NHP5YlomEyeMsP
7EtafvXj7rjoRmpkvvxnEL2v1tIBvWona5ZtL/XhntoE4vwSjLhbjBsCgYEA06On
eXnTYgNjsOWfqZOXPA2Ty8Jvl+IoKTO29e8FbA0QIIcrJ78Fr8kByREaQ2YkxF9X
Fdv6+lFCzRUSphZHYnBWrCHUc3+2Yuzo+dQZIYvViPW1GeheJhBFw5cGWuWcO9SF
uS7DglNo+NVhAoDLkTno3xjBOnn47gKjZqO8eRECgYAOlaU2guvZmn0rEcaOH/ac
4PushpqYjGaioQjZdws/s0qF1hXxYHRbK7c3qiYQ40avXDoznev9j28cxBckJu/M
52HUOzrGk9+L0jjBK1H0AnJjBKkweeuI3gN4Lc9XNm4JiH8GgOzmf2/uld549HKi
mrRsIxw7nF4uliNZKeD6uwKBgDDmxKC99IjWLafHNwAw2SYIIRlYwP5ARHVYvLLQ
2tjfn9VURjV13vOCJ4Z1DDN8m4xAV1f2r2Q9eIj4kImN5kqpmG1Hl9ZkMRlkkmR/
jJsCu4Fc/M6SsYZsBiKud8py+YmdjpR+aLBpY3zzmOnCJsdUsSkziBph6pHcTDNA
LCFxAoGBAKT+JjoeGKYbNuzgqGnQnIh4+3WUmKbGTEe44qmclXhq5G/ktsRojVBi
YgOxAbk487VpUDB9ptRskG+SgrrQ8uMIIh75Bo7V1tnjHRlLcWQvzN3EOtITGfTq
rAfYJRvFQeYNy3CiuxZpXq72n3EPxn+Chmn5lVuZN8igYhDyK7P/
-----END RSA PRIVATE KEY-----`
	pubKey := `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAulLSQ7sdr06pntGFWloy
wCFfaZGcxZ5HoRtpyAlpwmlzkK4RYoUrWbTc4fOy81dBpsiHwK3ZzefaPhMckhF/
oset0YOIhDRqcM3FflmrFofqlGgudTiJ0suH8s6IcSEtzo/wzq5T/KIfJHwVEhHt
ZGAnwjc2YQzsMz+KTeA5RSt8etNXL+mWuXLY5HB5A2EIXIiST6DaBPSH49TQn1pj
cHjXXeglmkwtjpSd+x6biA6YU0WYlgH2/aND7e+Pgtox99NXCIp6cw4ne+wdBJOy
eOtCYnBQvvF5n0+Jbgcc9Ox9D4fvTGZzRz83+bqkLRjgAD9CfSh6Ah1hIG+tmCEQ
ywIDAQAB
-----END PUBLIC KEY-----`
	conf := &config{
		ClientKey:  privKey,
		ClientName: "test_client",
	}

	vaultURL, _ := url.Parse("http://localhost/v1/chef-node/login")
	headers, _ := authHeaders(conf, vaultURL, "POST", nil, false)

	sigVer := headers.Get("X-Ops-Sign")
	sig := headers.Get("X-Ops-Authorization")
	ts := headers.Get("X-Ops-Timestamp")
	key, _ := parsePublicKey(pubKey)
	keys := []*rsa.PublicKey{key}
	if !authenticate("test_client", ts, sig, sigVer, keys, vaultURL.Path) {
		t.Fatal("Couldn't authenticate request")
	}
}

// This is an acceptance test.
// Requires the following env vars:
// VAULT_CLIENT_NAME - name of the client vault should connect to server as
// VAULT_CLIENT_KEYFILE - path to key for vault client
// VAULT_ADMIN_NAME - name of admin user for object creation
// VAULT_ADMIN_KEYFILE - path to admin's keyfile
// VAULT_CHEF_URL - Chef api endpoint
//
// The test requires that the admin user and the vault client already exist in the
// Chef server.
//
// It also requires that the ACLs on the chef server set the read permissions for
// the client used to connect on any newly created clients.  See the documentation
// for the backend to see how that might be done.
func TestBackendAcc_Login(t *testing.T) {
	if os.Getenv(logicaltest.TestEnvVar) == "" {
		t.Skip(fmt.Sprintf("Acceptance tests skipped unless env '%s' set", logicaltest.TestEnvVar))
		return
	}

	clientName := os.Getenv("VAULT_CLIENT_NAME")
	if clientName == "" {
		t.Fatalf("env var VAULT_CLIENT_NAME not set")
	}
	clientKeyFile := os.Getenv("VAULT_CLIENT_KEYFILE")
	if clientKeyFile == "" {
		t.Fatalf("env var VAULT_CLIENT_KEYFILE not set")
	}
	adminName := os.Getenv("VAULT_ADMIN_NAME")
	if adminName == "" {
		t.Fatalf("env var VAULT_ADMIN_NAME not set")
	}
	adminKeyFile := os.Getenv("VAULT_ADMIN_KEYFILE")
	if adminKeyFile == "" {
		t.Fatalf("env var VAULT_ADMIN_KEYFILE not set")
	}
	chefURL := os.Getenv("VAULT_CHEF_URL")
	if chefURL == "" {
		t.Fatalf("env var VAULT_CHEF_URL not set")
	}

	env := randString()
	err := setupTestEnv(env)
	if err != nil {
		t.Fatalf("Couldn't setup test environment %s", env)
	}
	defer teardownTestEnv(env)

	role1 := randString()
	role2 := randString()
	err = setupTestRole(role1)
	if err != nil {
		t.Fatalf("Couldn't setup test role %s", role1)
	}
	err = setupTestRole(role2)
	if err != nil {
		t.Fatalf("Couldn't setup test role %s", role2)
	}
	defer teardownTestRole(role1)
	defer teardownTestRole(role2)

	nodeName := randString()
	tag1 := randString()
	tag2 := randString()
	tagList := []string{tag1, tag2}
	roleList := []string{role1, role2}
	nodeKey, err := setupTestNode(nodeName, env, roleList, tagList)
	secondaryKey, err := addClientKey(nodeName)

	if err != nil {
		t.Fatalf("Couldn't setup test node %s", nodeName)
	}
	defer teardownTestNode(nodeName)

	storage := &logical.InmemStorage{}
	bconfig := logical.TestBackendConfig()
	bconfig.StorageView = storage
	b, err := Backend().Setup(bconfig)
	if err != nil {
		t.Fatal(err)
	}
	vaultKey, err := ioutil.ReadFile(clientKeyFile)
	if err != nil {
		t.Fatal(err)
	}
	vConfig := map[string]interface{}{
		"client_name": clientName,
		"client_key":  string(vaultKey),
		"base_url":    chefURL,
	}
	_, err = b.HandleRequest(&logical.Request{
		Operation: logical.UpdateOperation,
		Storage:   storage,
		Path:      "config",
		Data:      vConfig,
	})

	cpData := map[string]interface{}{
		"policies": "cp",
	}

	cpResp, err := b.HandleRequest(&logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "client/" + nodeName,
		Data:      cpData,
		Storage:   storage,
	})
	if err != nil {
		t.Fatal(err)
	}
	if cpResp != nil && cpResp.IsError() {
		t.Fatalf("Failed to create client")
	}

	epData := map[string]interface{}{
		"policies": "ep",
	}

	epResp, err := b.HandleRequest(&logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "environment/" + env,
		Data:      epData,
		Storage:   storage,
	})
	if err != nil {
		t.Fatal(err)
	}
	if epResp != nil && epResp.IsError() {
		t.Fatalf("Failed to create environment")
	}

	rpData := map[string]interface{}{
		"policies": "rp1",
	}
	rpResp, err := b.HandleRequest(&logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "role/" + role1,
		Data:      rpData,
		Storage:   storage,
	})
	if err != nil {
		t.Fatal(err)
	}
	if rpResp != nil && rpResp.IsError() {
		t.Fatalf("Failed to create first role")
	}
	rpData = map[string]interface{}{
		"policies": "rp2",
	}
	rpResp, err = b.HandleRequest(&logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "role/" + role2,
		Data:      rpData,
		Storage:   storage,
	})
	if err != nil {
		t.Fatal(err)
	}
	if rpResp != nil && rpResp.IsError() {
		t.Fatalf("Failed to create second role")
	}

	tpData := map[string]interface{}{
		"policies": "tp1",
	}
	tagResp, err := b.HandleRequest(&logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "tag/" + tag1,
		Data:      tpData,
		Storage:   storage,
	})
	if err != nil {
		t.Fatal(err)
	}
	if tagResp != nil && tagResp.IsError() {
		t.Fatal("Failed to create tag")
	}
	tpData = map[string]interface{}{
		"policies": "tp2",
	}
	tagResp, err = b.HandleRequest(&logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "tag/" + tag2,
		Data:      tpData,
		Storage:   storage,
	})
	if err != nil {
		t.Fatal(err)
	}
	if tagResp != nil && tagResp.IsError() {
		t.Fatal("Failed to create tag")
	}

	conf := &config{
		ClientName: nodeName,
		ClientKey:  string(nodeKey),
	}

	testURL, err := url.Parse("/v1/login")
	if err != nil {
		t.Fatal(err)
	}
	h, err := authHeaders(conf, testURL, "POST", nil, false)

	loginInput := map[string]interface{}{
		"signature_version": h.Get("X-Ops-Sign"),
		"client_name":       nodeName,
		"signature":         h.Get("X-Ops-Authorization"),
		"timestamp":         h.Get("X-Ops-Timestamp"),
	}

	loginRequest := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "login",
		Storage:   storage,
		Data:      loginInput,
	}

	resp, err := b.HandleRequest(loginRequest)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil || resp.Auth == nil || resp.IsError() {
		t.Fatalf("login attempt failed")
	}

	exPols := []string{"default", "cp", "ep", "rp1", "rp2", "tp1", "tp2"}
	if !policyutil.EquivalentPolicies(exPols, resp.Auth.Policies) {
		t.Fatalf("policies didn't match:\nexpected: %#v\ngot: %#v\n", exPols, resp.Auth.Policies)
	}

	conf2 := &config{
		ClientName: nodeName,
		ClientKey:  string(secondaryKey),
	}

	h2, err := authHeaders(conf2, testURL, "POST", nil, false)

	loginInput = map[string]interface{}{
		"signature_version": h2.Get("X-Ops-Sign"),
		"client_name":       nodeName,
		"signature":         h2.Get("X-Ops-Authorization"),
		"timestamp":         h2.Get("X-Ops-Timestamp"),
	}

	loginRequest = &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "login",
		Storage:   storage,
		Data:      loginInput,
	}

	resp, err = b.HandleRequest(loginRequest)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil || resp.Auth == nil || resp.IsError() {
		t.Fatalf("login attempt with secondary key failed")
	}

	if !policyutil.EquivalentPolicies(exPols, resp.Auth.Policies) {
		t.Fatalf("policies didn't match:\nexpected: %#v\ngot: %#v\n", exPols, resp.Auth.Policies)
	}
}

func randString() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, 20)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func chefRequest(object string, method string, body []byte) (*http.Response, error) {
	adminName := os.Getenv("VAULT_ADMIN_NAME")
	adminKeyFile := os.Getenv("VAULT_ADMIN_KEYFILE")
	chefURL := os.Getenv("VAULT_CHEF_URL")
	key, _ := ioutil.ReadFile(adminKeyFile)
	conf := &config{
		ClientName: adminName,
		ClientKey:  string(key),
	}
	url, _ := url.Parse(chefURL + "/" + object)
	bodyBuf := bytes.NewBuffer(body)
	headerBuf := bytes.NewBuffer(body)

	headers, err := authHeaders(conf, url, method, headerBuf, true)
	if err != nil {
		return nil, err
	}
	headers.Add("Content-Type", "application/json")
	req, err := http.NewRequest(method, url.String(), bodyBuf)
	req.Header = headers
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func setupTestEnv(env string) error {
	envData := map[string]string{
		"name": env,
	}
	envJSON, err := json.Marshal(envData)
	if err != nil {
		return err
	}
	_, err = chefRequest("environments", "POST", envJSON)
	return err
}

func setupTestRole(role string) error {
	roleData := map[string]string{
		"name": role,
	}
	roleJSON, err := json.Marshal(roleData)
	if err != nil {
		return err
	}
	_, err = chefRequest("roles", "POST", roleJSON)
	return err
}

func setupTestNode(name string, env string, roles []string, tags []string) (string, error) {
	clientData := struct {
		Name        string `json:"name"`
		GenerateKey bool   `json:"create_key"`
	}{
		name,
		true,
	}
	clientJSON, err := json.Marshal(clientData)
	if err != nil {
		return "", err
	}

	clientResp, err := chefRequest("clients", "POST", clientJSON)
	if err != nil {
		return "", err
	}
	defer clientResp.Body.Close()
	cBody, err := ioutil.ReadAll(clientResp.Body)
	if err != nil {
		return "", err
	}
	var clientRespStruct struct {
		PrivateKey string `json:"private_key"`
	}
	err = json.Unmarshal(cBody, &clientRespStruct)
	if err != nil {
		return "", err
	}

	nodeData := struct {
		Name        string `json:"name"`
		ChefEnv     string `json:"chef_environment"`
		NormalAttrs struct {
			Tags []string `json:"tags"`
		} `json:"normal"`
		AutoAttrs struct {
			Roles []string `json:"roles"`
		} `json:"automatic"`
	}{
		name,
		env,
		struct {
			Tags []string `json:"tags"`
		}{
			tags,
		},
		struct {
			Roles []string `json:"roles"`
		}{
			roles,
		},
	}
	nodeJSON, err := json.Marshal(nodeData)
	if err != nil {
		return "", err
	}
	_, err = chefRequest("nodes", "POST", nodeJSON)
	if err != nil {
		return "", err
	}

	return clientRespStruct.PrivateKey, nil
}

func addClientKey(name string) (string, error) {
	keyReq := struct {
		Name   string `json:"name"`
		Exp    string `json:"expiration_date"`
		Create bool   `json:"create_key"`
	}{
		name + "_2",
		"infinity",
		true,
	}
	keyJSON, err := json.Marshal(keyReq)
	if err != nil {
		return "", err
	}

	keyResp, err := chefRequest("clients/"+name+"/keys", "POST", keyJSON)
	if err != nil {
		return "", err
	}
	defer keyResp.Body.Close()
	kBody, err := ioutil.ReadAll(keyResp.Body)
	if err != nil {
		return "", err
	}

	var KeyRespStruct struct {
		PrivateKey string `json:"private_key"`
	}
	err = json.Unmarshal(kBody, &KeyRespStruct)
	if err != nil {
		return "", err
	}

	return KeyRespStruct.PrivateKey, nil
}

func teardownTestNode(name string) error {
	_, err := chefRequest("nodes/"+name, "DELETE", []byte(""))
	if err != nil {
		return err
	}
	_, err = chefRequest("clients/"+name, "DELETE", []byte(""))
	return err
}

func teardownTestRole(name string) error {
	_, err := chefRequest("roles/"+name, "DELETE", []byte(""))
	return err
}

func teardownTestEnv(name string) error {
	_, err := chefRequest("environments/"+name, "DELETE", []byte(""))
	return err
}
