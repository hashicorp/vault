// +build vault

package meta

func AdditionalOptionsUsage() string {
	return `
  -wrap-ttl=""            Indicates that the response should be wrapped in a
                          cubbyhole token with the requested TTL. The response
                          can be fetched by calling the "sys/wrapping/unwrap"
                          endpoint, passing in the wrappping token's ID. This
                          is a numeric string with an optional suffix
                          "s", "m", or "h"; if no suffix is specified it will
                          be parsed as seconds. May also be specified via
                          VAULT_WRAP_TTL.
`
}
