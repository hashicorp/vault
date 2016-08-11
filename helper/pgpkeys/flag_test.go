package pgpkeys

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/keybase/go-crypto/openpgp"
	"github.com/keybase/go-crypto/openpgp/packet"
)

func TestPubKeyFilesFlag_implements(t *testing.T) {
	var raw interface{}
	raw = new(PubKeyFilesFlag)
	if _, ok := raw.(flag.Value); !ok {
		t.Fatalf("PubKeysFilesFlag should be a Value")
	}
}

func TestPubKeyFilesFlagSetBinary(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "vault-test")
	if err != nil {
		t.Fatalf("Error creating temporary directory: %s", err)
	}
	defer os.RemoveAll(tempDir)

	decoder := base64.StdEncoding
	pub1Bytes, err := decoder.DecodeString(pubKey1)
	if err != nil {
		t.Fatalf("Error decoding bytes for public key 1: %s", err)
	}
	err = ioutil.WriteFile(tempDir+"/pubkey1", pub1Bytes, 0755)
	if err != nil {
		t.Fatalf("Error writing pub key 1 to temp file: %s", err)
	}
	pub2Bytes, err := decoder.DecodeString(pubKey2)
	if err != nil {
		t.Fatalf("Error decoding bytes for public key 2: %s", err)
	}
	err = ioutil.WriteFile(tempDir+"/pubkey2", pub2Bytes, 0755)
	if err != nil {
		t.Fatalf("Error writing pub key 2 to temp file: %s", err)
	}
	pub3Bytes, err := decoder.DecodeString(pubKey3)
	if err != nil {
		t.Fatalf("Error decoding bytes for public key 3: %s", err)
	}
	err = ioutil.WriteFile(tempDir+"/pubkey3", pub3Bytes, 0755)
	if err != nil {
		t.Fatalf("Error writing pub key 3 to temp file: %s", err)
	}

	pkf := new(PubKeyFilesFlag)
	err = pkf.Set(tempDir + "/pubkey1,@" + tempDir + "/pubkey2")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	err = pkf.Set(tempDir + "/pubkey3")
	if err == nil {
		t.Fatalf("err: should not have been able to set a second value")
	}

	expected := []string{strings.Replace(pubKey1, "\n", "", -1), strings.Replace(pubKey2, "\n", "", -1)}
	if !reflect.DeepEqual(pkf.String(), fmt.Sprint(expected)) {
		t.Fatalf("Bad: %#v", pkf)
	}
}

func TestPubKeyFilesFlagSetB64(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "vault-test")
	if err != nil {
		t.Fatalf("Error creating temporary directory: %s", err)
	}
	defer os.RemoveAll(tempDir)

	err = ioutil.WriteFile(tempDir+"/pubkey1", []byte(pubKey1), 0755)
	if err != nil {
		t.Fatalf("Error writing pub key 1 to temp file: %s", err)
	}
	err = ioutil.WriteFile(tempDir+"/pubkey2", []byte(pubKey2), 0755)
	if err != nil {
		t.Fatalf("Error writing pub key 2 to temp file: %s", err)
	}
	err = ioutil.WriteFile(tempDir+"/pubkey3", []byte(pubKey3), 0755)
	if err != nil {
		t.Fatalf("Error writing pub key 3 to temp file: %s", err)
	}

	pkf := new(PubKeyFilesFlag)
	err = pkf.Set(tempDir + "/pubkey1,@" + tempDir + "/pubkey2")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	err = pkf.Set(tempDir + "/pubkey3")
	if err == nil {
		t.Fatalf("err: should not have been able to set a second value")
	}

	expected := []string{pubKey1, pubKey2}
	if !reflect.DeepEqual(pkf.String(), fmt.Sprint(expected)) {
		t.Fatalf("bad: got %s, expected %s", pkf.String(), fmt.Sprint(expected))
	}
}

func TestPubKeyFilesFlagSetKeybase(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "vault-test")
	if err != nil {
		t.Fatalf("Error creating temporary directory: %s", err)
	}
	defer os.RemoveAll(tempDir)

	err = ioutil.WriteFile(tempDir+"/pubkey2", []byte(pubKey2), 0755)
	if err != nil {
		t.Fatalf("Error writing pub key 2 to temp file: %s", err)
	}

	pkf := new(PubKeyFilesFlag)
	err = pkf.Set("keybase:jefferai,@" + tempDir + "/pubkey2" + ",keybase:hashicorp")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	fingerprints := []string{}
	for _, pubkey := range []string(*pkf) {
		keyBytes, err := base64.StdEncoding.DecodeString(pubkey)
		if err != nil {
			t.Fatalf("bad: %v", err)
		}
		pubKeyBuf := bytes.NewBuffer(keyBytes)
		reader := packet.NewReader(pubKeyBuf)
		entity, err := openpgp.ReadEntity(reader)
		if err != nil {
			t.Fatalf("bad: %v", err)
		}
		if entity == nil {
			t.Fatalf("nil entity encountered")
		}
		fingerprints = append(fingerprints, hex.EncodeToString(entity.PrimaryKey.Fingerprint[:]))
	}

	exp := []string{
		"0f801f518ec853daff611e836528efcac6caa3db",
		"cf3d4694c9f57b28cb4092c2eb832c67eb5e8957",
		"91a6e7f85d05c65630bef18951852d87348ffc4c",
	}

	if !reflect.DeepEqual(fingerprints, exp) {
		t.Fatalf("bad: got \n%#v\nexpected\n%#v\n", fingerprints, exp)
	}
}

const pubKey1 = `mQENBFXbjPUBCADjNjCUQwfxKL+RR2GA6pv/1K+zJZ8UWIF9S0lk7cVIEfJiprzzwiMwBS5cD0da
rGin1FHvIWOZxujA7oW0O2TUuatqI3aAYDTfRYurh6iKLC+VS+F7H+/mhfFvKmgr0Y5kDCF1j0T/
063QZ84IRGucR/X43IY7kAtmxGXH0dYOCzOe5UBX1fTn3mXGe2ImCDWBH7gOViynXmb6XNvXkP0f
sF5St9jhO7mbZU9EFkv9O3t3EaURfHopsCVDOlCkFCw5ArY+DUORHRzoMX0PnkyQb5OzibkChzpg
8hQssKeVGpuskTdz5Q7PtdW71jXd4fFVzoNH8fYwRpziD2xNvi6HABEBAAG0EFZhdWx0IFRlc3Qg
S2V5IDGJATgEEwECACIFAlXbjPUCGy8GCwkIBwMCBhUIAgkKCwQWAgMBAh4BAheAAAoJEOfLr44B
HbeTo+sH/i7bapIgPnZsJ81hmxPj4W12uvunksGJiC7d4hIHsG7kmJRTJfjECi+AuTGeDwBy84TD
cRaOB6e79fj65Fg6HgSahDUtKJbGxj/lWzmaBuTzlN3CEe8cMwIPqPT2kajJVdOyrvkyuFOdPFOE
A7bdCH0MqgIdM2SdF8t40k/ATfuD2K1ZmumJ508I3gF39jgTnPzD4C8quswrMQ3bzfvKC3klXRlB
C0yoArn+0QA3cf2B9T4zJ2qnvgotVbeK/b1OJRNj6Poeo+SsWNc/A5mw7lGScnDgL3yfwCm1gQXa
QKfOt5x+7GqhWDw10q+bJpJlI10FfzAnhMF9etSqSeURBRW5AQ0EVduM9QEIAL53hJ5bZJ7oEDCn
aY+SCzt9QsAfnFTAnZJQrvkvusJzrTQ088eUQmAjvxkfRqnv981fFwGnh2+I1Ktm698UAZS9Jt8y
jak9wWUICKQO5QUt5k8cHwldQXNXVXFa+TpQWQR5yW1a9okjh5o/3d4cBt1yZPUJJyLKY43Wvptb
6EuEsScO2DnRkh5wSMDQ7dTooddJCmaq3LTjOleRFQbu9ij386Do6jzK69mJU56TfdcydkxkWF5N
ZLGnED3lq+hQNbe+8UI5tD2oP/3r5tXKgMy1R/XPvR/zbfwvx4FAKFOP01awLq4P3d/2xOkMu4Lu
9p315E87DOleYwxk+FoTqXEAEQEAAYkCPgQYAQIACQUCVduM9QIbLgEpCRDny6+OAR23k8BdIAQZ
AQIABgUCVduM9QAKCRAID0JGyHtSGmqYB/4m4rJbbWa7dBJ8VqRU7ZKnNRDR9CVhEGipBmpDGRYu
lEimOPzLUX/ZXZmTZzgemeXLBaJJlWnopVUWuAsyjQuZAfdd8nHkGRHG0/DGum0l4sKTta3OPGHN
C1z1dAcQ1RCr9bTD3PxjLBczdGqhzw71trkQRBRdtPiUchltPMIyjUHqVJ0xmg0hPqFic0fICsr0
YwKoz3h9+QEcZHvsjSZjgydKvfLYcm+4DDMCCqcHuJrbXJKUWmJcXR0y/+HQONGrGJ5xWdO+6eJi
oPn2jVMnXCm4EKc7fcLFrz/LKmJ8seXhxjM3EdFtylBGCrx3xdK0f+JDNQaC/rhUb5V2XuX6VwoH
/AtY+XsKVYRfNIupLOUcf/srsm3IXT4SXWVomOc9hjGQiJ3rraIbADsc+6bCAr4XNZS7moViAAcI
PXFv3m3WfUlnG/om78UjQqyVACRZqqAGmuPq+TSkRUCpt9h+A39LQWkojHqyob3cyLgy6z9Q557O
9uK3lQozbw2gH9zC0RqnePl+rsWIUU/ga16fH6pWc1uJiEBt8UZGypQ/E56/343epmYAe0a87sHx
8iDV+dNtDVKfPRENiLOOc19MmS+phmUyrbHqI91c0pmysYcJZCD3a502X1gpjFbPZcRtiTmGnUKd
OIu60YPNE4+h7u2CfYyFPu3AlUaGNMBlvy6PEpU=`
const pubKey2 = `mQENBFXbkJEBCADKb1ZvlT14XrJa2rTOe5924LQr2PTZlRv+651TXy33yEhelZ+V4sMrELN8fKEG
Zy1kNixmbq3MCF/671k3LigHA7VrOaH9iiQgr6IIq2MeIkUYKZ27C992vQkYLjbYUG8+zl5h69S4
0Ixm0yL0M54XOJ0gm+maEK1ZESKTUlDNkIS7l0jLZSYwfUeGXSEt6FWs8OgbyRTaHw4PDHrDEE9e
Q67K6IZ3YMhPOL4fVk4Jwrp5R/RwiklT+lNozWEyFVwPFH4MeQMs9nMbt+fWlTzEA7tI4acI9yDk
Cm1yN2R9rmY0UjODRiJw6z6sLV2T+Pf32n3MNSUOYczOjZa4VBwjABEBAAG0EFZhdWx0IFRlc3Qg
S2V5IDKJATgEEwECACIFAlXbkJECGy8GCwkIBwMCBhUIAgkKCwQWAgMBAh4BAheAAAoJEOuDLGfr
XolXqz4H/28IuoRxGKoJ064YHjPkkpoddW6zdzzNfHipZnNfEUiTEls4qF1IB81M2xqfiXIFRIdO
2kaLkRPFhO0hRxbtI6VuZYLgG3QCaXhxW6GyFa5zKABqhb5ojexdnAYRswaHV201ZCclj9rnJN1P
Ag0Rz6MdX/w1euEWktQxWKo42oZKyx8oT9p6lrv5KRmGkdrg8K8ARmRILjmwuBAgJM0eXBZHNGWX
elk4YmOgnAAcZA6ZAo1G+8Pg6pKKP61ETewuCg3/u7N0vDttB+ZXqF88W9jAYlvdgbTtajNF5IDY
DjTzWfeCaIB18F9gOzXq15SwWeDDI+CU9Nmq358IzXlxk4e5AQ0EVduQkQEIAOjZV5tbpfIh5Qef
pIp2dpGMVfpgPj4RNc15CyFnb8y6dhCrdybkY9GveXJe4F3GNYnSfB42cgxrfhizX3LakmZQ/SAg
+YO5KxfCIN7Q9LPNeTgPsZZT6h8lVuXUxOFKXfRaR3/tGF5xE3e5CoZRsHV/c92h3t1LdJNOnC5m
UKIPO4zDxiw/C2T2q3rP1kmIMaOH724kEH5A+xcp1cBHyt0tdHtIWuQv6joTJzujqViRhlCwQYzQ
SKpSBxwhBsorPvyinZI/ZXA4XXZc5RoMqV9rikedrb1rENO8JOuPu6tMS+znFu67skq2gFFZwCQW
IjdHm+2ukE+PE580WAWudyMAEQEAAYkCPgQYAQIACQUCVduQkQIbLgEpCRDrgyxn616JV8BdIAQZ
AQIABgUCVduQkQAKCRArYtevdF38xtzgB/4zVzozBpVOnagRkA7FDsHo36xX60Lik+ew0m28ueDD
hnV3bXQsCvn/6wiCVWqLOTDeYCPlyTTpEMyk8zwdCICW6MgSkVHWcEDOrRqIrqm86rirjTGjJSgQ
e3l4CqJvkn6jybShYoBk1OZZV6vVv9hPTXXv9E6dLKoEW5YZBrrF+VC0w1iOIvaAQ+QXph20eV4K
BIrp/bhG6PdnigKxuBZ79cdqDnXIzT9UiIa6LYpR0rbeg+7BmuZTTPS8t+41hIiKS+UZFdKa67eY
ENtyOmEMWOFCLLRJGxkleukchiMJ70rknloZXsvJIweXBzSZ6m7mJQBgaig/L/dXyjv6+j2pNB4H
/1trYUtJjXQKHmqlgCmpCkHt3g7JoxWvglnDNmE6q3hIWuVIYQpnzZy1g05+X9Egwc1WVpBB02H7
PkUZTfpaP/L6DLneMmSKPhZE3I+lPIPjwrxqh6xy5uQezcWkJTNKvPWF4FJzrVvx7XTPjfGvOB0U
PEnjvtZTp5yOhTeZK7DgIEtb/Wcrqs+iRArQKboM930ORSZhwvGK3F9V/gMDpIrvge5vDFsTEYQd
w/2epIewH0L/FUb/6jBRcVEpGo9Ayg+Jnhq14GOGcd1y9oMZ48kYVLVBTA9tQ+82WE8Bch7uFPj4
MFOMVRn1dc3qdXlg3mimA+iK7tABQfG0RJ9YzWs=`
const pubKey3 = `mQENBFXbkiMBCACiHW4/VI2JkfvSEINddS7vE6wEu5e1leNQDaLUh6PrATQZS2a4Q6kRE6WlJumj
6wCeN753Cm93UGQl2Bi3USIEeArIZnPTcocrckOVXxtoLBNKXgqKvEsDXgfw8A+doSfXoDm/3Js4
Wy3WsYKNR9LaPuJZHnpjsFAJhvRVyhH4UFD+1RTSSefq1mozPfDdMoZeZNEpfhwt3DuTJs7RqcTH
CgR2CqhEHnOOE5jJUljHKYLCglE2+8dth1bZlQi4xly/VHZzP3Bn7wKeolK/ROP6VZz/e0xq/BKy
resmxvlBWZ1zWwqGIrV9b0uwYvGrh2hOd5C5+5oGaA2MGcjxwaLBABEBAAG0EFZhdWx0IFRlc3Qg
S2V5IDOJATgEEwECACIFAlXbkiMCGy8GCwkIBwMCBhUIAgkKCwQWAgMBAh4BAheAAAoJEPR5S1b8
LcbdWjEH/2mhqC9a0Vk1IzOgxEoVxYVqVdvaxI0nTZOTfmcFYn4HQlQ+SLEoyNWe5jtkhx4k5uHi
pxwKHzOv02YM14NWC6bvKw2CQETLDPG4Cv8YMUmpho5tnMDdttIzp8HjyJRtHazU1uTes2/yuqh6
LHCejVJI0uST3RibquwdG3QjPP8Umxu+YC9+FOW2Kit/AQ8JluFDJdq3/wSX8VfYZrGdgmreE7KY
MolhCkzGSPj7oFygw8LqKoJvt9tCuBKhZMBuMv1sB5CoJIWdPoqOZc4U7L1XdqfKvFZR/RhuXgN1
lkI9MqrnLDpikL3Lk+ctLxWOjUCW8roqKoHZYBF7XPqdAfm5AQ0EVduSIwEIAOPcjd4QgbLlqIk3
s6BPRRyVzglTgUdf+I0rUDybaDJfJobZd8U6e4hkPvRoQ8tJefnz/qnD/63watAbJYcVTme40I3V
KDOmVGcyaDxiKP1disKqcEJd7XQiI72oAiXmEH0y+5UwnOMks/lwaAGDMGVRjHEXI6fiRPFsfTr8
7qvMJ3pW1OiOXVSezuBNTlmyJC7srQ1/nwxL337ev6D1zQZd3JuhcxLkHrUELLNwzhvcZ70vg645
jAmz8EdmvvoqEPPoHqKgP5AeHACOsTm953KHhgx3NYuGPU/RoIvugKt4Iq5nw7TWFTjPHGVF3GTQ
ry5CZ/AzXiL57hVEhDvmuT8AEQEAAYkCPgQYAQIACQUCVduSIwIbLgEpCRD0eUtW/C3G3cBdIAQZ
AQIABgUCVduSIwAKCRAFI/9Nx3K5IPOFCACsZ/Z4s2LcEoA51TW+T5w+YevlIuq+332JtqNIpuGI
WpGxUxyDyPT0YQWr0SObBORYNr7RP8d/I2rbaFKyaDaKvRofYr+TwXy92phBo7pdEUamBpfrm/sr
+2BgAB2x3HWXp+IMdeVVhqQe8t4cnFm3c1fIdxADyiJuV5ge2Ml5gK5yNwqCQPh7U2RqC+lmVlMJ
GvWInIRn2mf6A7phDYNZfOz6dkar4yyh5r9rRgrZw88r/yIlrq/c6KRUIgnPMrFYEauggOceZ827
+jKkbKWFEuHtiCxW7kRAN25UfnGsPaF+NSGM2q1vCG4HiFydx6lMoXM0Shf8+ZwyrV/5BzAqpWwI
AJ37tEwC58Fboynly6OGOzgPS0xKnzkXMOtquTo0qEH/1bEUsBknn795BmZOTf4oBC5blN6qRv7c
GSP00i+sxql1NhTjJcjJtfOPUzgfW+Af/+HR648z4c7c6MCjDFKnk8ZkoGLRU7ISjenkNFzvu2bj
lxJkil0uJDlLPbbX80ojzV1GS9g+ZxVPR+68N1QLl2FU6zsfg34upmLLHG8VG4vExzgyNkOwfTYv
dgyRNTjnuPue6H12fZZ9uCNeG52v7lR3eoQcCxBOniwgipB8UJ52RWXblwxzCtGtDi/EWB3zLTUn
puKcgucA0LotbihSMxhDylaARfVO1QV6csabM/g=`
