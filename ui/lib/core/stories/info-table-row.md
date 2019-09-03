<!--THIS FILE IS AUTO GENERATED. This file is generated from JSDoc comments in lib/core/addon/components/info-table-row.js. To make changes, first edit that file and run "yarn gen-story-md info-table-row" to re-generate the content.-->

## InfoTableRow
`InfoTableRow` displays a label and a value in a table-row style manner. The component is responsive so
that the value breaks under the label on smaller viewports.


| Param | Type | Default | Description |
| --- | --- | --- | --- |
| value | <code>any</code> | <code></code> | The the data to be displayed - by default the content of the component will only show if there is a value. Also note that special handling is given to boolean values - they will render `Yes` for true and `No` for false. |
| label | <code>string</code> | <code>null</code> | The display name for the value. |
| alwaysRender | <code>Boolean</code> | <code>false</code> | Indicates if the component content should be always be rendered.  When false, the value of `value` will be used to determine if the component should render. |

**Example**
  
```js
<InfoTableRow @value={{5}} @label="TTL" />
```

**See**

- [Uses of InfoTableRow](https://github.com/hashicorp/vault/search?l=Handlebars&q=InfoTableRow+OR+info-table-row)
- [InfoTableRow Source Code](https://github.com/hashicorp/vault/blob/master/ui/app/components/info-table-row.js)

---
