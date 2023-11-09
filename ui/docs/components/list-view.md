# ListView

`ListView` components are used in conjunction with `ListItem` for rendering a list.

| Param                 | Type                | Default                               | Description                                                                                                                     |
| --------------------- | ------------------- | ------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------- |
| [items]               | <code>array</code>  | <code></code>                         | An Ember array of items (objects) to render as a list. Because it's an Ember array it has properties like length an meta on it. |
| [itemNoun]            | <code>string</code> | <code>&quot;item&quot;</code>         | A noun to use in the empty state of message and title.                                                                          |
| [message]             | <code>string</code> | <code>null</code>                     | The message to display within the banner.                                                                                       |
| [paginationRouteName] | <code>string</code> | <code>&quot;&#x27;&#x27;&quot;</code> | The link used in the ListPagination component.                                                                                  |

**Example**

```hbs preview-template
<!-- empty state -->
<ListView @itemNoun='role' as |list|>
  {{#if list.empty}}
    <list.empty @title='No roles here' />
  {{/if}}
</ListView>

<!-- with items -->
<ListView @items={{array (hash id='my-role')}} @itemNoun='role' as |list|>
  {{#if list.item}}
    <ListItem @linkPrefix='vault' as |Item|>
      <Item.content>
        <Icon @name='folder' class='has-text-grey-light' />{{list.item.id}}
      </Item.content>
      <Item.menu as |m|>
        <li class='action'>
          <LinkTo @route='vault' class='is-block'>
            Some action
          </LinkTo>
        </li>
      </Item.menu>
    </ListItem>
  {{/if}}
</ListView>
```
