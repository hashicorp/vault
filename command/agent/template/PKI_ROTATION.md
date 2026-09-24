# Draft PKI rotation API

This prototype depends on the Consul Template `Runner.ForcePKIRefresh` API.
It is disabled by default and only supports `pkiCert`, not `secret` or
`pkiCertExternalCa`.

Enable it on a listener restricted to trusted operators, for example a Unix
socket with filesystem access controls:

```hcl
listener "unix" {
  address = "/run/vault-agent/agent.sock"
  socket_mode = "0600"
  tls_disable = true
  agent_api {
    enable_pki_rotate = true
  }
}
```

```sh
curl --unix-socket /run/vault-agent/agent.sock \
  -H 'Content-Type: application/json' \
  -d '{"destination":"/etc/service/certificate.rendered"}' \
  http://localhost/agent/v1/templates/rotate
```

The destination must exactly match a configured template destination. No path is
accepted for deletion or arbitrary template execution. The API itself does not
perform Vault token authorization: listener access is the authorization boundary,
like the existing local quit API. Do not expose it on an untrusted network.
Authorized calls can issue certificates and run existing template commands.
When `require_request_header` is enabled, include `X-Vault-Request: true`.

A `202` response with `status: accepted` means the existing PKI watcher was
signaled. It does not mean issuance, rendering, or the configured `exec` command
has finished successfully. Check the resulting certificate serial and the Agent
logs. The request has a one-second acceptance timeout; issuance uses the existing
client timeout and retry configuration. A timed-out HTTP request can still have
been accepted; this prototype has no operation ID or exactly-once guarantee.

The active watcher bypasses the valid certificate cached in the destination and
requests a new certificate with the current token. It then uses the normal
render and command path. It does not delete the current files, replace the
runner, restart Agent, or force authentication. Failed issuance leaves the
existing rendered certificate intact. Rotation intent survives fetch retries
within the same dependency, but is not persisted across runner/process restarts.
Pending signals coalesce; a request arriving during issuance may cause another
issuance. Normal renewal scheduling resumes after the fresh certificate.

Responses: `400` invalid JSON or missing destination; `404` disabled API or
unknown destination; `405` non-POST; `409` one-shot mode or no active `pkiCert`
dependency; `503` unavailable/stopped/not-ready template server. A template may
be configured but have no active dependency before its initial render.

This API does not implement a maintenance-window scheduler. A systemd timer or
other external scheduler can call it. Upstream API review, a released Consul
Template dependency, and product documentation are required before release.
