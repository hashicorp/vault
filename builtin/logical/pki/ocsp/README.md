# PKI OCSP Package

This package provides Online Certificate Status Protocol (OCSP) request and response parsing, creation, and signature validation for Vault PKI.

## Provenance

- **Upstream Source:** [`golang.org/x/crypto/ocsp`](https://pkg.go.dev/golang.org/x/crypto/ocsp)
- **Upstream Version:** `v0.58.0`
- **Upstream SHA:** `8fdf07b3114278fb3b78c1cf1814e08bd851eb4c`
- **Original License:** BSD 3-Clause (The Go Authors)

## Motivation & Modifications

The upstream `golang.org/x/crypto/ocsp` package only supports classical signature algorithms (RSA, ECDSA, DSA). This Vault-owned copy extends the implementation to support modern and post-quantum signature schemes:

1. **ML-DSA (NIST FIPS 204):**
   - Added OID definitions and mappings for ML-DSA-44 (`2.16.840.1.101.3.4.3.17`), ML-DSA-65 (`2.16.840.1.101.3.4.3.18`), and ML-DSA-87 (`2.16.840.1.101.3.4.3.19`).
   - Integrated with standard library `crypto/mldsa` and `crypto/x509` signature algorithms (`x509.MLDSA44`, `x509.MLDSA65`, `x509.MLDSA87`).
   - Enabled direct raw message signing (`crypto.Hash(0)`) in `CreateResponse`.

2. **Ed25519 (RFC 8410 / RFC 8032):**
   - Added OID definition (`1.3.101.112`) and mapping for `x509.PureEd25519`.
   - Enabled direct raw message signing (`crypto.Hash(0)`).
