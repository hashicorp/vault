<!--THIS FILE IS AUTO GENERATED. This file is generated from JSDoc comments in lib/core/addon/components/stat-text.js. To make changes, first edit that file and run "yarn gen-story-md stat-text" to re-generate the content.-->

## StatText
StatText components are used to display a label and associated statistic below, with the option to add a description.

**Params**

| Param | Type | Default | Description |
| --- | --- | --- | --- |
| label | <code>string</code> | <code>null</code> | the label for the statistic |
| value | <code>string</code> | <code>null</code> | value passed in, usually a number or statistic |
| [size] | <code>string</code> | <code>&quot;m&quot;</code> | size the component as small or large, 's', 'm' or 'l' |
| [subText] | <code>string</code> |  | subText is optional and will display below the label |

**Example**
  
```js
<StatText @label="Active Clients" @stat="4,198" @size="l" @subText="These are the active client counts"/>
```

**See**

- [Uses of StatText](https://github.com/hashicorp/vault/search?l=Handlebars&q=StatText+OR+stat-text)
- [StatText Source Code](https://github.com/hashicorp/vault/blob/master/ui/lib/core/addon/components/stat-text.js)

---
