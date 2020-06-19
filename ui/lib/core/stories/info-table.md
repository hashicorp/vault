<!--THIS FILE IS AUTO GENERATED. This file is generated from JSDoc comments in lib/core/addon/components/info-table.js. To make changes, first edit that file and run "yarn gen-story-md info-table" to re-generate the content.-->

## InfoTable
InfoTable components are a table with a single column and header. They are used to render a list of InfoTableRow components.

**Params**

| Param | Type | Default | Description |
| --- | --- | --- | --- |
| [title] | <code>String</code> | <code>Info Table</code> | The title of the table. Used for accessibility purposes. |
| header | <code>String</code> | <code></code> | The column header. |
| items | <code>Array</code> | <code></code> | An array of strings which will be used as the InfoTableRow value. |

**Example**
  
```js
<InfoTable
  @title="Known Primary Cluster Addrs"
  @header="cluster_addr"
  @items={{knownPrimaryClusterAddrs}}
/>
```
        

**See**

- [Uses of InfoTable](https://github.com/hashicorp/vault/search?l=Handlebars&q=InfoTable+OR+info-table)
- [InfoTable Source Code](https://github.com/hashicorp/vault/blob/master/ui/lib/core/addon/components/info-table.js)

---
