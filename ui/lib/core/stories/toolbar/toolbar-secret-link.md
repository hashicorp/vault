<!--THIS FILE IS AUTO GENERATED. This file is generated from JSDoc comments in app/components/toolbar-secret-link.js. To make changes, first edit that file and run "yarn gen-story-md toolbar-secret-link" to re-generate the content.-->

## ToolbarSecretLink
`ToolbarSecretLink` styles SecretLink for the Toolbar.
It should only be used inside of `Toolbar`.


| Param | Type | Default | Description |
| --- | --- | --- | --- |
| type | <code>String</code> | <code>&quot;&quot;</code> | Use "add" to change icon |

**Example**
  
```js
<Toolbar>
  <ToolbarActions>
    <ToolbarSecretLink @params={{array 'vault.cluster.policies.create'}} @type="add">
      Create policy
    </ToolbarSecretLink>
  </ToolbarActions>
</Toolbar>
```

**See**

- [Uses of ToolbarSecretLink](https://github.com/hashicorp/vault/search?l=Handlebars&q=ToolbarSecretLink+OR+toolbar-secret-link)
- [ToolbarSecretLink Source Code](https://github.com/hashicorp/vault/blob/master/ui/app/components/toolbar-secret-link.js)

---
