<!--THIS FILE IS AUTO GENERATED. This file is generated from JSDoc comments in app/components/auth-config-form/config.js. To make changes, first edit that file and run "yarn gen-story-md auth-config-form/config" to re-generate the content.-->

## ToolbarSecretLink
`ToolbarSecretLink` styles SecretLink for the Toolbar.

**Properties**

ToolbarSecretLink takes the same properties as SecretLink, but allows you to set the type for the icon.

| Name | Type | Default | Description |
| --- | --- | --- | --- |
| type | <code>String</code> | <code></code> | Use "add" to change icon |


**Example**

```js
  <ToolbarSecretLink
    @params={{array 'vault.cluster.policies.create'}}
    @type="add"
  >
    Create policy
  </ToolbarSecretLink>
```

**See**

- [Uses of ToolbarSecretLink](https://github.com/hashicorp/vault/search?l=Handlebars&q=ToolbarSecretLink)
- [ToolbarSecretLink Source Code](https://github.com/hashicorp/vault/blob/master/ui/app/components/toolbar-secret-link.js)

---
