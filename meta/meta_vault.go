// +build vault

package meta

func AdditionalOptionsUsage() string {
	return `
  -wrap-ttl=""            Indiciates that the response should be wrapped in a
                          cubbyhole token with the requested TTL. The response
                          will live at "cubbyhole/response" in the cubbyhole of
                          the returned token with a key of "response" and can
                          be parsed as a normal API Secret. The backend can
                          also request wrapping; the lesser of the values is
                          used. This is a numeric string with an optional
                          suffix "s", "m", or "h"; if no suffix is specified it
                          will be parsed as seconds. May also be specified via
                          VAULT_WRAP_TTL.
`
}
