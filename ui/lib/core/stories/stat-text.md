<!--THIS FILE IS AUTO GENERATED. This file is generated from JSDoc comments in lib/core/addon/components/stat-text.js. To make changes, first edit that file and run "yarn gen-story-md stat-text" to re-generate the content.-->

## StatText
StatText components are used to display a label and associated value beneath, with the option to include a description.

**Params**

| Param | Type | Default | Description |
| --- | --- | --- | --- |
| label | <code>string</code> | <code>null</code> | The label for the statistic |
| value | <code>string</code> | <code>null</code> | Value passed in, usually a number or statistic |
| size | <code>string</code> | <code>null</code> | Sizing changes whether or not there is subtext. If there is subtext 's' and 'l' are valid sizes. If no subtext, then 'm' is also acceptable. |
| [subText] | <code>string</code> |  | SubText is optional and will display below the label |

**Example**
  
```js
<StatText @label="Active Clients" @stat="4,198" @size="l" @subText="These are the active client counts"/>
```

**See**

- [Uses of StatText](https://github.com/hashicorp/vault/search?l=Handlebars&q=StatText+OR+stat-text)
- [StatText Source Code](https://github.com/hashicorp/vault/blob/master/ui/lib/core/addon/components/stat-text.js)

---
