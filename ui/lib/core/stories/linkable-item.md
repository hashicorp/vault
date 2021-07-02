<!--THIS FILE IS AUTO GENERATED. This file is generated from JSDoc comments in lib/core/addon/components/linkable-item.js. To make changes, first edit that file and run "yarn gen-story-md linkable-item" to re-generate the content.-->

## LinkableItem
LinkableItem components have two contextual components, a Content component used to show information on the left with a Menu component on the right, all aligned vertically centered. If passed a link, the block will be clickable.

**Params**

| Param | Type | Default | Description |
| --- | --- | --- | --- |
| [link] | <code>object</code> | <code></code> | Link should have route and model |
| [disabled] | <code>boolean</code> | <code>false</code> | If no link then should be given a disabled attribute equal to true |

**Example**
  
```js
<LinkableItem @link={{hash route='vault.backends' model='my-backend-path'}} data-test-row="my-backend-path" as |Li|/>
```

**See**

- [Uses of LinkableItem](https://github.com/hashicorp/vault/search?l=Handlebars&q=LinkableItem+OR+linkable-item)
- [LinkableItem Source Code](https://github.com/hashicorp/vault/blob/master/ui/lib/core/addon/components/linkable-item.js)

---
