package openpgp

import (
	"context"
	"testing"

	"github.com/hashicorp/vault/logical"
)

func TestPGP_ShowSessionKey(t *testing.T) {
	storage := &logical.InmemStorage{}
	b := Backend()

	// add a new private key to vault
	req := &logical.Request{
		Storage:   storage,
		Operation: logical.UpdateOperation,
		Path:      "keys/test",
		Data: map[string]interface{}{
			"generate": false,
			"key":      privateSessionDecryptKey,
		},
	}
	_, err := b.HandleRequest(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}

	showSessionKey := func(keyName, ciphertext, format, signerKey, expected string) {
		reqDecrypt := &logical.Request{
			Storage:   storage,
			Operation: logical.UpdateOperation,
			Path:      "show-session-key/" + keyName,
			Data: map[string]interface{}{
				"ciphertext": ciphertext,
				"format":     format,
				"signer_key": signerKey,
			},
		}

		resp, err := b.HandleRequest(context.Background(), reqDecrypt)
		if err != nil {
			t.Fatal(err)
		}
		if resp.IsError() {
			t.Fatalf("not expected error response: %#v", *resp)
		}

		if resp == nil {
			t.Fatalf("no name key found in response data %#v", resp)
		}
		plaintext, ok := resp.Data["session_key"]
		if !ok {
			t.Fatalf("no name key found in response data %#v", resp.Data)
		}
		if plaintext != expected {
			t.Fatalf("expected plaintext %s, got: %s", expected, plaintext)
		}
	}

	showSessionKey("test", encryptedSessionMessageAsciiArmored, "ascii-armor", "", "9:BDF8F7A2A573556C1E7D2FE9ADDCA7188C451C60B5311025F2A900E9FC61809E")
	showSessionKey("test", encryptedSessionMessageBase64Encoded, "base64", "", "9:EC211D19FA4FFC7F88B6AC6A1112C88032910753AB52FEF10C71D850A721151C")
	showSessionKey("test", encryptedAndSignedSessionMessageAsciiArmored, "ascii-armor", publicSessionSignerKey, "9:E82F18EF202C2356CB758B56B79D4F4CCAD6CC21FDFBB6867FCFAC25BC10863C")
	showSessionKey("test", encryptedSessionMessageAsciiArmored[:398], "ascii-armor", "", "9:BDF8F7A2A573556C1E7D2FE9ADDCA7188C451C60B5311025F2A900E9FC61809E")
	showSessionKey("test", encryptedSessionMessageBase64EncodedWithMultipleKeys, "base64", "", "9:F8054D6D0F6E9C89155B829BC71E0613472EA70E32B9DA7893960536B04BB2BD")
}

func TestPGP_ShowSessionKeyError(t *testing.T) {
	storage := &logical.InmemStorage{}
	b := Backend()

	// generate a new key in the vault
	req := &logical.Request{
		Storage:   storage,
		Operation: logical.UpdateOperation,
		Path:      "keys/testGenerated",
		Data: map[string]interface{}{
			"real_name": "Vault PGP test",
		},
	}
	_, err := b.HandleRequest(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}

	// add a new key to the vault
	req = &logical.Request{
		Storage:   storage,
		Operation: logical.UpdateOperation,
		Path:      "keys/test",
		Data: map[string]interface{}{
			"generate": false,
			"key":      privateSessionDecryptKey,
		},
	}
	_, err = b.HandleRequest(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}

	showSessionKeyMustFail := func(keyName, ciphertext, format, signerKey string) {
		reqDecrypt := &logical.Request{
			Storage:   storage,
			Operation: logical.UpdateOperation,
			Path:      "show-session-key/" + keyName,
			Data: map[string]interface{}{
				"ciphertext": ciphertext,
				"format":     format,
				"signer_key": signerKey,
			},
		}

		resp, _ := b.HandleRequest(context.Background(), reqDecrypt)
		if !resp.IsError() {
			t.Fatalf(
				"expected to fail, keyname: %s, format: %s, cipertext: %s, signer key %s",
				keyName, format, ciphertext, signerKey)
		}
	}

	showSessionKeyMustFail("doNotExist", encryptedSessionMessageAsciiArmored, "ascii-armor", "")
	showSessionKeyMustFail("test", encryptedSessionMessageAsciiArmored, "invalidFormat", "")
	showSessionKeyMustFail("test", encryptedSessionMessageAsciiArmored[:128], "ascii-armor", "")
	showSessionKeyMustFail("test", encryptedSessionMessageAsciiArmored[:256], "ascii-armor", "")
	showSessionKeyMustFail("test", encryptedSessionMessageAsciiArmored[:386], "ascii-armor", "")

	// Wrong key for the message
	showSessionKeyMustFail("testGenerated", encryptedSessionMessageAsciiArmored, "ascii-armor", "")

	// Wrongly encoded
	showSessionKeyMustFail("test", "Not ASCII armored", "ascii-armor", "")
	showSessionKeyMustFail("test", "Not base64 encoded", "base64", "")

	// Signer key is not properly ASCII-armored
	showSessionKeyMustFail("test", encryptedSessionMessageAsciiArmored, "ascii-armor", "Signer key is not ASCII armored")

}

const privateSessionDecryptKey = `-----BEGIN PGP PRIVATE KEY BLOCK-----

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
UWauhQ==
=2Rz1
-----END PGP PRIVATE KEY BLOCK-----`

const publicSessionSignerKey = `-----BEGIN PGP PUBLIC KEY BLOCK-----

mQENBFmcjO0BCACgAjFrcoCmtIjMHvNc12wfZqOVxs2hvk2FD5wFtupVOTmBrHzU
Jry1wb6J6azrHlGXC3lz5tpC2AczhjNSOXBKMR3VxPfzVJwS3aHC6+DrTsmIJYYD
kDABCiyD1Br7W5wMSa4KCK8/YSsopD7j4ffSgjijZ5ytFmeYeMDzpVP3DO+uo3Mo
TyU4/DsW1jE78IJKJ/10ZvtP0ZLT9UXBPqwo4kGFQTjoTXf++0y3qXQ5eGwe2AbY
K7z6fyMp7q8HrJD3bor1uiV2ookCm0576twI90gTsKXU9SkGT6vKweqbRfCjpaSq
8K+IVm7Hfjm6kL18f+SzJ0hSvJt1FtICZzTXABEBAAG0K1ZhdWx0IFNpZ25lciBL
ZXkgPHZhdWx0LXNpZ25lckBleGFtcGxlLmNvbT6JAVQEEwEKAD4WIQQ1dVo1ZrPi
hI+ct5sP4181AmkrmwUCWZyM7QIbAwUJA8JnAAULCQgHAwUVCgkICwUWAgMBAAIe
AQIXgAAKCRAP4181AmkrmyhLB/4+8ffovwEnlF0EK0lOrn0xYLiKI3pOZuHw9y9F
qb1ViWSH3E8raLnjt970YMnpirrzI12qe1Pexzr6eYw9kuZRAu8xP2cryWviauZe
RiiP3MDRe4MT9f05eZAoaL1wEDbgAbZhbGmlhwc+CBo3Sue03rXVZZc2+pGpe6o5
K1mkDXuvekXsAmrQGo8PzMdJe18xido7hT1YnVntVrA1mWLBTij5sw4yAwVxuUN/
CH1bmOP6Y+zef7/pkgSFS5v+XRibKrvlOUHgvKC4J/ruFc7IOK55mvaRzgliQAg8
qwaTlakao9PnWmcKPXKsosFI01CCmE3RDdEXM62y3rlboAY0uQENBFmcjO0BCAC6
HyHdPQWS0+pGcBhGipm5GJmyyrwQuPiS5/gPecK+F9QysnIX6xFkJM13U/4J2p7t
vUPvSHWs4/Qtpi8w3n3Yddm17qGZUHIb4QJsc8hUasTtjuDHOYpgEd51ujTWJ1k5
O39NDso1cRihXUPgI4GzhL6l1afLcqp1Y5YQyqEFLaz/IZMMlzt0zhAJLV0oyNtQ
Bumm3szcUOyRpkgGvtpFJgcSxb2tzZosUwxykct9ALQP/Q+WLhXqeo9bv5ziJtQm
y9kFGDhTwVaKXDOHYypKSz1R18MMcxHHQZjbBVnSBUaPCXuMx651F5pQF0tsbEaQ
6K1WiVowFpm0w9Tp/t8VABEBAAGJATwEGAEKACYWIQQ1dVo1ZrPihI+ct5sP4181
AmkrmwUCWZyM7QIbDAUJA8JnAAAKCRAP4181Amkrm//pB/wK4FjGPA9kHQ6TJRW/
Sv6Mup9Cw3b7eOOMTGasKkk6w3kal/jCyCsfA0SbceQUQb06pOMTEzht545qiUJH
+T8/7A8/ak7GvLOTYLbig5qFNHWmoxB0ychrsc3siVtqSsrm9FhuY2O4475XODnB
FR1GNuxuvH8wpVI/viJdDluExX1ltWI6L5xEk/1CatcAKhc3yx9ggz/n7qrPh1iV
FWxBhYvEEvQaLU8Y7ALAd/f5+iFioSbW1hoaZDXtWvXuR70piZAsVCCkf+hZc6am
eccF2zPehP6Jv+0eXHpC9NmaZH6A1zAB2cvFY8x6PQ6JOjqUEfmEJ3I2+eA8ru6s
Jt9I
=+eRe
-----END PGP PUBLIC KEY BLOCK-----`

const encryptedSessionMessageAsciiArmored = `-----BEGIN PGP MESSAGE-----

hQEMA923ECy/uCBhAQf8DLagsnoLuM4AyKiTyvZ7uSQTkmOkwXwn1WWsxoKJkzdI
v2XJ7knQ3UR5nnhI8xVbAnZVZjx8wYaBPUvV2VqhA2sTn36mGlGw43ngDOFB1cKW
1VM9JY0xqxuHaIR3mvYFjb/iuoT2BM7SmCuIEJYgxKEM+/R1o9rkCenj2pOj4+XK
ryXv+iHQAar6Ic2G3g9T7Mu7Uw6+n1xBWr/XzPnJRJf4WB4m7sqd/Wm7NkHnvgde
P9kawh1lHYj32WdLUqZpQB3zQRguDHFfQA8vRVEG4Gyz/o7um5PFc4kDES0JYzNc
p6p64MAF+vMpSOsFU2TaixSmraidaWHVPYcao/w2UNJDAQ43l9lh064yz9bCaH41
UyEQpNH+l1EpqnIbu+iIQb3a02GwBB8lfEW7cFku8121H8XapkgKZDsmXD/7v0eW
e8iwFg==
=+yfj
-----END PGP MESSAGE-----`

const encryptedSessionMessageBase64Encoded = `hQEMA923ECy/uCBhAQf/XPUNCcaIUyTDDQ+rII/sj24VtnBUdXDNntOtBX4pxIHzMWr6oCWGgZZV
WTRzRP4nEclUUhWKHDlEd7/1bG/1z3Px3JWXdnSCHYl3AqdFkS4bW26wpO+gcCTbiixo+JE93QoG
84rb5k6gdNGsVEpioDFK1FLGL9pPvyR+kp4JRg8qD1FpDsvow+zhJqgAak87s4Ly/YnYiVYbGjPl
u0pqEkvJwHnIyKThFW5N6OCYjB2pFpVLER7x6RGjuX6tRRYZayzT4sVKGj0Efp6T32EEVPURiJSn
elpIPEd8+8i/7X0Co6iNFEyucgxhaxN+ujqSxx+6ZIFV4UKC0LFgR2iF99JDAQ6ofxvUtoxMGKON
WVtrVMjN8Db3KXQ5rt/tyKbTVGXQot6ocSZ2Ae+rnSTiq0boGrWDnuYZHawc16iJhbcP68ERgg==`

const encryptedAndSignedSessionMessageAsciiArmored = `-----BEGIN PGP MESSAGE-----

hQEMA923ECy/uCBhAQf/fWXnkS/aq982Df+9NjqYna9c8aAcQRuhVi0jc0rasRRj
owPqag1s+0PeaMD07e02+RvWRhDzXnd9OdR4Tm+91e2DFhFQ16OnLf5C/EQhaNk4
LVwKtGAP9IiZJVx+NZNJeadGU7QEFDjnZdLLu7Hh6sTqs6RmopjAucNWNIwVJjn0
+ehlDwfSHOf3aV8ET0hBgIasSgWNUQ8gRrzMdDffZuqovHIOPTpRkyeD1uNFhP5j
TASsDvRPA4YXsqgaXkgpIKmegYWpRGlJAhqB7TzeY+E/+pvqhLSUBUNN+CZpfots
8Xx9furYvCh2ajyKrOgdJcMiU1XUwCqJbPlqcWMyo9LAyQFdkjGNt/dZZfUXpph7
MErulIRedVYhgGa12bLosCeG92n8qAhAoc8tVZxSkJBfwaWOoo8i049buCRDGAIG
AwB3LiMVTe2gfhpidBgPibGcXD5RxIZcyWDKQLAhKe2Ig9+pRSgyZh+GuNkzuyhf
mxPM4WpLMXD/OqAT12jlrCxoDsttoLIvLmGItDjDnO/+pPiw/emMTAsfQPk/VxOz
ASACh6kHAdV/WVUwzDouxIbbEu49vyWWed1ls/8MRaXRwTU7AQGDN/B0BLe99bzZ
Ya1/VsAG4nz6eH2duU6nUI03tNuRjXAdeR2GKOdz8pQYVAcr1YTg1XQ4r4X5HAId
HfP1TjFTrCTztKvpk6D7R++D32zlwjTqyqIrYXJvcifev7zd7tCRGk7D+qvIdrF2
Ubm5BjtILbgkwpbSWBghW+lx5POhVt9mFax+Su9fZkUrPj3UGnHH2jeFB4EwHtkI
TSpU+MkEN1+Gdp+peD7lHSgfOxvpfJt4qA8ic89DSWF1YYK8a8CkiiqnMQ==
=Bepf
-----END PGP MESSAGE-----`

const encryptedSessionMessageBase64EncodedWithMultipleKeys = `hQEMA53ITnFABb0rAQf9Ftqb+riRjV5wz5ghTRlUdEnTx+3NDSINJsHqHgBSj0Leb3ICMmCryFk2
ql8eYWlX04YNgUJjoAUXyDlXaW/F4B62Z+vYiVqcGqwK96UYz6e1AEWj8ir7vrq8W1gnwb0Rx99Q
KkPEbD4gf6xBrnBsPcbo/o2Pxkf9jcFc0Fq78AXrpt+dTyvnLjMnvzcO4XOwXzswq7poG4AGC5Pb
lx40+SKeMN0u7v3/nVk2i0GKJeT4rqq2kD/K07zzBggQd9FLbmGAO/oq44aczizNGiBUTaoidcOG
Rw7oiZvFiA4BMRTdtAGd8+TqBxe19jIuQ0kZD5Bk8VLP+DKAPtnU9ZfSHYUBDAPdtxAsv7ggYQEH
/11QmDPPsScTXb/YWu1tcO20QtDckQ3IYP20oQt4IV+LCIoKd7xJlwyMY467WL0s+7qfSSoi/zdP
GyzZd8spC5m/rRp5tSlauUTYGPgkfun94tU8IZhQUwtZNIj0oxuf4cCyqtQXlycRELK7uFsVpWbs
af1Alu1jB06ZyoPvrddODuwONGtBaeX4SasQX0P3bVZubISOcfQ+xYpQ2kZiUQu2ZMAVJRHB/jfR
x31eMBp5QXSnNvt1G67yl+M9kwppTiToYEEuyRuXPZ0YmwESSLM/EVLl8KGukq7g/5ZmkVMaLped
A710RWmTpQ2+4h5zGQfIJQDzvPEcG6aKtXqoyF3SSgFssSU+J5aMtOmNwhmMbdvQIlwaBNYi6sFM
yMlVzPaMuFD3Kk9RN8jiM2m1Mnn4MPhse/PsAaNFTHYulApGLt/QuSjLjQey0HBb`
