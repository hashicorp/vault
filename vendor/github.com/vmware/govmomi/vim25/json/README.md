# JSON with Discriminators

The source code in this directory was copied from Go 1.17.13's `encoding/json` package in order to add support for JSON discriminators. Please use the following command to review the diff:

```shell
C1="$(git log --pretty=format:'%h' --no-patch --grep='Vendor Go 1.17.13 encoding/json')" && \
C2="$(git log --pretty=format:'%h' --no-patch --grep='JSON Encoding w Discriminator Support')" && \
git diff "${C1}".."${C2}"
```
