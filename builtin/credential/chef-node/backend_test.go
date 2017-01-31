package chefnode

import (
	"net/url"
	"testing"

	"strings"

	"github.com/hashicorp/vault/helper/policyutil"
	"github.com/hashicorp/vault/logical"
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
		t.Fatal("Failed to create environment")
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
		t.Fatal("Failed to read environment")
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
	headers, _ := authHeaders(conf, vaultURL, "POST", false)

	sigVer := headers.Get("X-Ops-Sign")
	sig := headers.Get("X-Ops-Authorization")
	ts := headers.Get("X-Ops-Timestamp")
	key, _ := parsePublicKey(pubKey)

	if !authenticate("test_client", ts, sig, sigVer, key, vaultURL.Path) {
		t.Fatal("Couldn't authenticate request")
	}
}
