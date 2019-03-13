<a name="ToggleButton
`ToggleButton` components are used to expand and collapse content with a toggle.module_"></a>

## ToggleButton
`ToggleButton` components are used to expand and collapse content with a toggle.

**See**

- [Uses of ToggleButton](https://github.com/hashicorp/vault/search?l=Handlebars&q=ToggleButton)
- [ToggleButton Source Code](https://github.com/hashicorp/vault/blob/master/ui/app/components/toggle-button.js)

**Properties**

| Name | Type | Default | Description |
| --- | --- | --- | --- |
| toggleAttr | <code>String</code> | <code></code> | The attribute upon which to toggle. |
| attrTarget | <code>Object</code> | <code></code> | The target upon which the event handler should be added. |
| [openLabel] | <code>String</code> | <code>Hide options</code> | The message to display when the toggle is open. |
| [closedLabel] | <code>String</code> | <code>More options</code> | The message to display when the toggle is closed. |

**Example**  

```hbs
  <ToggleButton @openLabel="Encrypt Output with PGP" @closedLabel="Encrypt Output with PGP" @toggleTarget={{this}} @toggleAttr="showOptions"/>
 {{#if showOptions}}
    <div>
      <p>
        I will be toggled!
      </p>
    </div>
  {{/if}}```