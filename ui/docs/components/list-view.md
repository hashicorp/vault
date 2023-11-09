
# ListView
&#x60;ListView&#x60; components are used in conjunction with &#x60;ListItem&#x60; for rendering a list.

| Param | Type | Default | Description |
| --- | --- | --- | --- |
| [items] | <code>array</code> | <code></code> | An Ember array of items (objects) to render as a list. Because it's an Ember array it has properties like length an meta on it. |
| [itemNoun] | <code>string</code> | <code>&quot;item&quot;</code> | A noun to use in the empty state of message and title. |
| [message] | <code>string</code> | <code>null</code> | The message to display within the banner. |
| [paginationRouteName] | <code>string</code> | <code>&quot;&#x27;&#x27;&quot;</code> | The link used in the ListPagination component. |

**Example**  
```hbs preview-template
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
