## 0.1.1 (May 2, 2015)

IMPROVEMENTS:

  * core: Very verbose error if mlock fails [GH-59]
  * command/*: On error with TLS oversized record, show more human-friendly
      error message. [GH-123]
  * command/read: `lease_renewable` is now outputed along with the secret
      to show whether it is renewable or not
  * command/server: Add configuration option to disable mlock
  * command/server: Disable mlock for dev mode so it works on more systems

BUG FIXES:

  * core: if token helper isn't absolute, prepend with path to Vault
      executable, not "vault" (which requires PATH) [GH-60]
  * core: Any "mapping" routes allow hyphens in keys [GH-119]
  * core: Validate `advertise_addr` is a valid URL with scheme [GH-106]
  * command/auth: Using an invalid token won't crash [GH-75]
  * credential/app-id: app and user IDs can have hyphens in keys [GH-119]
  * helper/password: import proper DLL for Windows to ask password [GH-83]
  * physical/file: create the storge with 0600 permissions [GH-102]

## 0.1.0 (April 28, 2015)

  * Initial release
