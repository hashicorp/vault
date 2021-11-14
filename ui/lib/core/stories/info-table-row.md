<!--THIS FILE IS AUTO GENERATED. This file is generated from JSDoc comments in lib/core/addon/components/info-table-row.js. To make changes, first edit that file and run "yarn gen-story-md info-table-row" to re-generate the content.-->

## InfoTableRow
`InfoTableRow` displays a label and a value in a table-row style manner. The component is responsive so
that the value breaks under the label on smaller viewports.

**Params**

| Param | Type | Default | Description |
| --- | --- | --- | --- |
| value | <code>any</code> | <code></code> | The the data to be displayed - by default the content of the component will only show if there is a value. Also note that special handling is given to boolean values - they will render `Yes` for true and `No` for false. |
| label | <code>string</code> | <code>null</code> | The display name for the value. |
| helperText | <code>string</code> | <code>null</code> | Text to describe the value displayed beneath the label. |
| alwaysRender | <code>Boolean</code> | <code>false</code> | Indicates if the component content should be always be rendered.  When false, the value of `value` will be used to determine if the component should render. |
| [type] | <code>string</code> | <code>&quot;array&quot;</code> | The type of value being passed in.  This is used for when you want to trim an array.  For example, if you have an array value that can equal length 15+ this will trim to show 5 and count how many more are there |
| [isLink] | <code>Boolean</code> | <code>true</code> | Passed through to InfoTableItemArray. Indicates if the item should contain a link-to component.  Only setup for arrays, but this could be changed if needed. |
| [modelType] | <code>string</code> | <code>null</code> | Passed through to InfoTableItemArray. Tells what model you want data for the allOptions to be returned from.  Used in conjunction with the the isLink. |
| [queryParam] | <code>String</code> |  | Passed through to InfoTableItemArray. If you want to specific a tab for the View All XX to display to.  Ex: role |
| [backend] | <code>String</code> |  | Passed through to InfoTableItemArray. To specify secrets backend to point link to  Ex: transformation |
| [viewAll] | <code>String</code> |  | Passed through to InfoTableItemArray. Specify the word at the end of the link View all. |
| [tooltipText] | <code>String</code> |  | Text if a tooltip should display over the value. |
| [isTooltipCopyable] | <code>Boolean</code> |  | Allows tooltip click to copy |
| [defaultShown] | <code>String</code> |  | Text that renders as value if alwaysRender=true. Eg. "Vault default" |

**Example**
  
```js
<InfoTableRow @value={{5}} @label="TTL" @helperText="Some description"/>
```

**See**

- [Uses of InfoTableRow](https://github.com/hashicorp/vault/search?l=Handlebars&q=InfoTableRow+OR+info-table-row)
- [InfoTableRow Source Code](https://github.com/hashicorp/vault/blob/master/ui/lib/core/addon/components/info-table-row.js)

---
