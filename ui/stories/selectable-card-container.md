<!--THIS FILE IS AUTO GENERATED. This file is generated from JSDoc comments in app/components/selectable-card-container.js. To make changes, first edit that file and run "yarn gen-story-md selectable-card-container" to re-generate the content.-->

## SelectableCardContainer
SelectableCardContainer components are used to hold SelectableCard components.  They act as a CSS grid container, and change grid configurations based on the boolean of @gridContainer.

**Params**

| Param | Type | Default | Description |
| --- | --- | --- | --- |
| counters | <code>Object</code> | <code></code> | Counters is an object that returns total entities, tokens, and an array of objects with the total https request per month. |
| gridContainer | <code>Boolean</code> | <code>false</code> | gridContainer is optional.  If true, it's telling the container it will have a nested CSS grid. |

**Example**
  
```js
<SelectableCardContainer @counters={{model}} @gridContainer="true" />
```

**See**

- [Uses of SelectableCardContainer](https://github.com/hashicorp/vault/search?l=Handlebars&q=SelectableCardContainer+OR+selectable-card-container)
- [SelectableCardContainer Source Code](https://github.com/hashicorp/vault/blob/master/ui/app/components/selectable-card-container.js)

---
