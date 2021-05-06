# jwt
[![Go Reference](https://pkg.go.dev/badge/github.com/hashicorp/cap/jwt.svg)](https://pkg.go.dev/github.com/hashicorp/cap/jwt)

Package jwt provides signature verification and claims set validation for JSON Web Tokens (JWT)
of the JSON Web Signature (JWS) form.

Primary types provided by the package:

* `KeySet`: Represents a set of keys that can be used to verify the signatures of JWTs.
  A KeySet is expected to be backed by a set of local or remote keys.
  
* `Validator`: Provides signature verification and claims set validation behavior for JWTs.

* `Expected`: Defines the expected claims values to assert when validating a JWT.

* `Alg`: Represents asymmetric signing algorithms.

### Examples:

Please see [docs_test.go](./docs_test.go) for additional usage examples.
