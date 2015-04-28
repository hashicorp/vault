## 0.1.1 (unreleased)

IMPROVEMENTS:

  * command/server: Add configuration option to disable mlock
  * command/server: Disable mlock for dev mode so it works on more systems

BUG FIXES:

  * core: if token helper isn't absolute, prepend with path to Vault
      executable, not "vault" (which requires PATH) [GH-60]

## 0.1.0 (April 28, 2015)

  * Initial release
