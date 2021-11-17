JOSE
============
[![Build Status](https://travis-ci.org/SermoDigital/jose.svg?branch=master)](https://travis-ci.org/SermoDigital/jose)
[![GoDoc](https://godoc.org/github.com/SermoDigital/jose?status.svg)](https://godoc.org/github.com/SermoDigital/jose)

JOSE is a comprehensive set of JWT, JWS, and JWE libraries.

## Why

The only other JWS/JWE/JWT implementations are specific to JWT, and none
were particularly pleasant to work with.

These libraries should provide an easy, straightforward way to securely
create, parse, and validate JWS, JWE, and JWTs.

## Notes:
JWE is currently unimplemented.

## Version 0.9:

## Documentation

The docs can be found at [godoc.org] [docs], as usual.

A gopkg.in mirror can be found at https://gopkg.in/jose.v1, thanks to
@zia-newversion. (For context, see issue #30.) 

### [JWS RFC][jws]
### [JWE RFC][jwe]
### [JWT RFC][jwt]

## License

[MIT] [license].

[docs]:    https://godoc.org/github.com/SermoDigital/jose
[license]: https://github.com/SermoDigital/jose/blob/master/LICENSE.md
[jws]: https://tools.ietf.org/html/rfc7515
[jwe]: https://tools.ietf.org/html/rfc7516
[jwt]: https://tools.ietf.org/html/rfc7519
