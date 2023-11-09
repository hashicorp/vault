
# TtlPicker
TtlPicker components are used to enable and select duration values such as TTL.
This component renders a toggle by default, and passes all relevant attributes
to TtlForm. Please see that component for additional arguments
- allows TTL to be enabled or disabled
- recalculates the time when the unit is changed by the user (eg 60s -&gt; 1m)

| Param | Type | Default | Description |
| --- | --- | --- | --- |
| onChange | <code>function</code> |  | This function will be passed a TTL object, which includes enabled{bool}, seconds{number}, timeString{string}, goSafeTimeString{string}. |
| initialEnabled | <code>Boolean</code> | <code>false</code> | Set this value if you want the toggle on when component is mounted |
| label | <code>String</code> | <code>&quot;Time</code> | to live (TTL)"  - Label is the main label that lives next to the toggle. Yielded values will replace the label |
| labelDisabled |  | <code>Label</code> | to display when TTL is toggled off |
| helperTextEnabled | <code>String</code> | <code>&quot;&quot;</code> | This helper text is shown under the label when the toggle is switched on |
| helperTextDisabled | <code>String</code> | <code>&quot;&quot;</code> | This helper text is shown under the label when the toggle is switched off |
| initialValue | <code>string</code> | <code>null</code> | InitialValue is the duration value which will be shown when the component is loaded. If it can't be parsed, will default to 0. |
| changeOnInit | <code>boolean</code> | <code>false</code> | if true, calls the onChange hook when component is initialized |
| hideToggle | <code>Boolean</code> | <code>false</code> | set this value if you'd like to hide the toggle and just leverage the input field |

**Example**  
```hbs preview-template
<TtlPicker @onChange={{this.handleChange}} @initialEnabled={{@model.myAttribute}} @initialValue={{@model.myAttribute}}/>
```
