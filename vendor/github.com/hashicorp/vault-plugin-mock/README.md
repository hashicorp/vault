# vault-plugin-mock

This is a mock Vault plugin for testing. 

It *may* be imported by the real Vault, so while it is functionally a NOOP, changes can break Vault. In particular:

- Don't rename the package
- Don't rename or remove the exported constant
- Ensure it always compiles
