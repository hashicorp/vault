<!--THIS FILE IS AUTO GENERATED. This file is generated from JSDoc comments in lib/core/addon/components/toggle-button.js. To make changes, first edit that file and run "yarn gen-story-md toggle-button" to re-generate the content.-->

## onClickCallback : <code>function</code>

**Params**

| Param | Type | Default | Description |
| --- | --- | --- | --- |
| isOpen | <code>boolean</code> |  | determines whether to show open or closed label |
| onClick | [<code>onClickCallback</code>](#onClickCallback) |  | fired when button is clicked |
| [openLabel] | <code>string</code> | <code>&quot;\&quot;Hide options\&quot;&quot;</code> | The message to display when the toggle is open. |
| [closedLabel] | <code>string</code> | <code>&quot;\&quot;More options\&quot;&quot;</code> | The message to display when the toggle is closed. |

**Example**
  
```js
  <ToggleButton @isOpen={{this.showOptions}} @openLabel="Encrypt Output with PGP" @closedLabel="Encrypt Output with PGP" @onClick={{fn (mut this.showOptions}} />
 {{#if showOptions}}
    <div>
      <p>
        I will be toggled!
      </p>
    </div>
  {{/if}}
```

**See**

- [Uses of ToggleButton](https://github.com/hashicorp/vault/search?l=Handlebars&q=ToggleButton+OR+toggle-button)
- [ToggleButton Source Code](https://github.com/hashicorp/vault/blob/main/ui/lib/core/addon/components/toggle-button.js)

---
