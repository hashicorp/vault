<!--THIS FILE IS AUTO GENERATED. This file is generated from JSDoc comments in app/components/pagination-controls.js. To make changes, first edit that file and run "yarn gen-story-md pagination-controls" to re-generate the content.-->

## PaginationControls
PaginationControls components are used to paginate through item lists

**Params**

| Param | Type | Default | Description |
| --- | --- | --- | --- |
| total | <code>number</code> |  | total number of items |
| [startPage] | <code>number</code> | <code>1</code> | initial page number to select |
| [size] | <code>number</code> | <code>15</code> | number of items to display per page |
| onChange | <code>function</code> |  | callback fired on page change |

**Example**
  
```js
<PaginationControls @startPage={{1}} @total={{100}} @size={{15}} @onChange={{this.onPageChange}} />
```

**See**

- [Uses of PaginationControls](https://github.com/hashicorp/vault/search?l=Handlebars&q=PaginationControls+OR+pagination-controls)
- [PaginationControls Source Code](https://github.com/hashicorp/vault/blob/master/ui/app/components/pagination-controls.js)

---
