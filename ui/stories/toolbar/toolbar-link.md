<!--THIS FILE IS AUTO GENERATED. This file is generated from JSDoc comments in app/components/auth-config-form/config.js. To make changes, first edit that file and run "yarn gen-story-md auth-config-form/config" to re-generate the content.-->

## ToolbarLink
`ToolbarLink` styles links for the Toolbar.

**Properties**

| Name | Type | Default | Description |
| --- | --- | --- | --- |
| params | <code>Array</code> | <code></code> | Array top pass to LinkTo params |
| type | <code>String</code> | <code></code> | Use "add" to change icon |


**Example**

```js
  <ToolbarLink
    @params={{array 'vault.cluster.policies.create'}}
    @type="add"
  >
    Create policy
  </ToolbarLink>
```

**See**

- [Uses of ToolbarLink](https://github.com/hashicorp/vault/search?l=Handlebars&q=ToolbarLink)
- [ToolbarLink Source Code](https://github.com/hashicorp/vault/blob/master/ui/lib/core/addon/components/toolbar-link.js)

---
