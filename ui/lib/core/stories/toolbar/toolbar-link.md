<!--THIS FILE IS AUTO GENERATED. This file is generated from JSDoc comments in lib/core/addon/components/toolbar-link.js. To make changes, first edit that file and run "yarn gen-story-md toolbar-link" to re-generate the content.-->

## ToolbarLink
`ToolbarLink` components style links and buttons for the Toolbar
It should only be used inside of `Toolbar`.

**Params**

| Param | Type | Description |
| --- | --- | --- |
| params | <code>array</code> | Array to pass to LinkTo |
| type | <code>string</code> | Use "add" to change icon |
| disabled | <code>boolean</code> | pass true to disable link |
| disabledTooltip | <code>string</code> | tooltip to display on hover when disabled |

**Example**
  
```js
<Toolbar>
  <ToolbarActions>
    <ToolbarLink @params={{array 'vault.cluster.policies.create'}} @type="add" @disabled={{true}} @disabledTooltip="This link is disabled">
      Create policy
    </ToolbarLink>
  </ToolbarActions>
</Toolbar>
```

**See**

- [Uses of ToolbarLink](https://github.com/hashicorp/vault/search?l=Handlebars&q=ToolbarLink+OR+toolbar-link)
- [ToolbarLink Source Code](https://github.com/hashicorp/vault/blob/master/ui/lib/core/addon/components/toolbar-link.js)

---
