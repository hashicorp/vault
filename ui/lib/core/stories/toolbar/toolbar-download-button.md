<!--THIS FILE IS AUTO GENERATED. This file is generated from JSDoc comments in app/components/toolbar-download-button.js. To make changes, first edit that file and run "yarn gen-story-md toolbar-download-button" to re-generate the content.-->

## ToolbarSecretLink
`ToolbarSecretLink` styles SecretLink for the Toolbar.
It should only be used inside of `Toolbar`.

**Example**

```js
<Toolbar>
  <ToolbarActions>
    <ToolbarDownloadButton @actionText="Download policy" @extension={{if (eq policyType "acl") model.format "sentinel"}} @filename={{model.name}} @data={{model.policy}} />
  </ToolbarActions>
</Toolbar>
```

**See**

- [Uses of ToolbarDownloadButton](https://github.com/hashicorp/vault/search?l=Handlebars&q=ToolbarDownloadButton+OR+toolbar-download-button)
- [ToolbarDownloadButton Source Code](https://github.com/hashicorp/vault/blob/master/ui/app/components/toolbar-download-button.js)

---
