<!--THIS FILE IS AUTO GENERATED. This file is generated from JSDoc comments in app/components/selectable-card.js. To make changes, first edit that file and run "yarn gen-story-md selectable-card" to re-generate the content.-->

## SelectableCard
SelectableCard components are card-like components that display a title, total, subtotal, and anything after the yield.
They are designed to be used in containers that act as flexbox or css grid containers.

**Params**

| Param | Type | Default | Description |
| --- | --- | --- | --- |
| cardTitle | <code>String</code> | <code></code> | cardTitle displays the card title |
| total | <code>Number</code> | <code>0</code> | the Total number displays like a title, it's the largest text in the component |
| subText | <code>String</code> | <code></code> | subText describes the total |

**Example**
  
```js
<SelectableCard @cardTitle="Tokens" @total={{totalHttpRequests}} @subText="Total"/>
```

**See**

- [Uses of SelectableCard](https://github.com/hashicorp/vault/search?l=Handlebars&q=SelectableCard+OR+selectable-card)
- [SelectableCard Source Code](https://github.com/hashicorp/vault/blob/master/ui/app/components/selectable-card.js)

---
