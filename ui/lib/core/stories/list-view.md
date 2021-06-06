<!--THIS FILE IS AUTO GENERATED. This file is generated from JSDoc comments in lib/core/addon/components/list-view.js. To make changes, first edit that file and run "yarn gen-story-md list-view" to re-generate the content.-->

## ListView
`ListView` components are used in conjuction with `ListItem` for rendering a list.


| Param | Type | Default | Description |
| --- | --- | --- | --- |
| items | <code>Array</code> | <code></code> | An array of items to render as a list |
| [itemNoun] | <code>String</code> | <code></code> | A noun to use in the empty state of message and title. |
| [message] | <code>String</code> | <code></code> | The message to display within the banner. |

**Example**
  
```js
<ListView @items={{model}} @itemNoun="role" @paginationRouteName="scope.roles" as |list|>
  {{#if list.empty}}
    <list.empty @title="No roles here" />
  {{else}}
    <div>
      {{list.item.id}}
    </div>
  {{/if}}
</ListView>
```

**See**

- [Uses of ListView](https://github.com/hashicorp/vault/search?l=Handlebars&q=ListView+OR+list-view)
- [ListView Source Code](https://github.com/hashicorp/vault/blob/master/ui/lib/core/addon/components/list-view.js)

---
