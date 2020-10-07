<!--THIS FILE IS AUTO GENERATED. This file is generated from JSDoc comments in lib/core/addon/components/box-radio.js. To make changes, first edit that file and run "yarn gen-story-md box-radio" to re-generate the content.-->

## BoxRadio
BoxRadio components are used to display options for a radio selection.

**Params**

| Param | Type | Default | Description |
| --- | --- | --- | --- |
| displayName | <code>string</code> |  | This is the string that will show on the box radio option. |
| type | <code>string</code> |  | type is the key that the radio input will be identified by. Please use a value without spaces. |
| glyph | <code>string</code> |  | glyph is the name of the icon that will be used in the box |
| groupValue | <code>string</code> |  | The key of the radio option that is currently selected for this radio group |
| groupName | <code>string</code> |  | The name (key) of the group that this radio option belongs to |
| onRadioChange | <code>function</code> |  | This callback will trigger when the radio option is selected (if enabled) |
| [disabled] | <code>boolean</code> | <code>false</code> | This parameter controls whether the radio option is selectable. If not, it will be grayed out and show a tooltip. |
| [tooltipMessage] | <code>string</code> | <code>&quot;default&quot;</code> | The message that shows in the tooltip if the radio option is disabled |

**Example**
  
```js
<BoxRadio @displayName="Catahoula Leopard" @type="catahoula" @glyph="dog" @groupValue="labrador" @groupName="my-favorite-dog" @onRadioChange={{handleRadioChange}} />
```

**See**

- [Uses of BoxRadio](https://github.com/hashicorp/vault/search?l=Handlebars&q=BoxRadio+OR+box-radio)
- [BoxRadio Source Code](https://github.com/hashicorp/vault/blob/master/ui/lib/core/addon/components/box-radio.js)

---
