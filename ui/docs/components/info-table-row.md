
# InfoTableRow
&#x60;InfoTableRow&#x60; displays a label and a value in a table-row style manner. The component is responsive so
that the value breaks under the label on smaller viewports.

| Param | Type | Default | Description |
| --- | --- | --- | --- |
| label | <code>string</code> | <code>null</code> | The display name for the value. |
| helperText | <code>string</code> | <code>null</code> | Text to describe the value displayed beneath the label. |
| value | <code>any</code> | <code></code> | The the data to be displayed - by default the content of the component will only show if there is a value. Also note that special handling is given to boolean values - they will render `Yes` for true and `No` for false. Overridden by block if exists |
| [alwaysRender] | <code>boolean</code> | <code>false</code> | Indicates if the component content should be always be rendered.  When false, the value of `value` will be used to determine if the component should render. |
| [truncateValue] | <code>boolean</code> | <code>false</code> | Indicates if the value should be truncated. |
| [defaultShown] | <code>string</code> |  | Text that renders as value if alwaysRender=true. Eg. "Vault default" |
| [tooltipText] | <code>string</code> |  | Text if a tooltip should display over the value. |
| [isTooltipCopyable] | <code>boolean</code> |  | Allows tooltip click to copy |
| [formatDate] | <code>string</code> |  | A string of the desired date format that's passed to the date-format helper to render timestamps (ex. "MMM d yyyy, h:mm:ss aaa", see: https://date-fns.org/v2.30.0/docs/format) |
| [formatTtl] | <code>boolean</code> | <code>false</code> | When true, value is passed to the format-duration helper, useful for TTL values |
| [type] | <code>string</code> | <code>&quot;array&quot;</code> | The type of value being passed in.  This is used for when you want to trim an array.  For example, if you have an array value that can equal length 15+ this will trim to show 5 and count how many more are there * InfoTableItemArray * |
| [isLink] | <code>boolean</code> | <code>true</code> | Passed through to InfoTableItemArray. Indicates if the item should contain a link-to component.  Only setup for arrays, but this could be changed if needed. |
| [modelType] | <code>string</code> | <code>null</code> | Passed through to InfoTableItemArray. Tells what model you want data for the allOptions to be returned from.  Used in conjunction with the the isLink. |
| [queryParam] | <code>string</code> |  | Passed through to InfoTableItemArray. If you want to specific a tab for the View All XX to display to.  Ex= role |
| [backend] | <code>string</code> |  | Passed through to InfoTableItemArray. To specify secrets backend to point link to  Ex= transformation |
| [viewAll] | <code>string</code> |  | Passed through to InfoTableItemArray. Specify the word at the end of the link View all. |

**Example**  
```hbs preview-template
<InfoTableRow @value={{5}} @label="TTL" @helperText="Some description"/>
```
