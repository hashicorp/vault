# ToggleButton

`ToggleButtons` are used to expand and collapse content with a toggle.

## Properties
| Property | Required | Default value | Type | Description | Example |
|---|---|---|---|---|
| toggleAttr | [x] | `null` | string | The attribute upon which to toggle. | `'showOptions'` |
| toggleTarget | [x] | `null` | element | The target upon which the event handler should be added. | `this` |
| openLabel || `'Hide options'` | string | The message to display when the toggle is open. ||
| closedLabel || `'More options'` | string | The message to display when the toggle is closed. ||

## Usage

```javascript
  <ToggleButton
    @openLabel="Encrypt Output with PGP"
    @closedLabel="Encrypt Output with PGP"
    @toggleTarget={{this}}
    @toggleAttr="showOptions"
  />

  {{#if showOptions}}
    <div>
      <p>
        I will be toggled!
      </p>
    </div>
  {{/if}}
```
https://github.com/hashicorp/vault/search?l=Handlebars&q=ToggleButton

## Source
https://github.com/hashicorp/vault/blob/master/ui/app/components/toggle-button.js
