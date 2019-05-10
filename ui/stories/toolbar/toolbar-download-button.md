<!--THIS FILE IS AUTO GENERATED. This file is generated from JSDoc comments in app/components/auth-config-form/config.js. To make changes, first edit that file and run "yarn gen-story-md auth-config-form/config" to re-generate the content.-->

## ToolbarDownloadButton
`ToolbarDownloadButton` styles DownloadButton for the Toolbar.

**Example**

```js
  <ToolbarDownloadButton
    @actionText="Download policy"
    @extension={{if (eq policyType "acl") model.format "sentinel"}}
    @filename={{model.name}}
    @data={{model.policy}}
  />
```

**See**

- [Uses of ToolbarDownloadButton](https://github.com/hashicorp/vault/search?l=Handlebars&q=ToolbarDownloadButton)
- [ToolbarDownloadButton Source Code](https://github.com/hashicorp/vault/blob/master/ui/app/components/toolbar-download-button.js)

---
