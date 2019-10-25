<!--THIS FILE IS AUTO GENERATED. This file is generated from JSDoc comments in lib/core/addon/components/search-select.js. To make changes, first edit that file and run "yarn gen-story-md search-select" to re-generate the content.-->

## SearchSelect
The `SearchSelect` is an implementation of the [ember-power-select-with-create](https://github.com/poteto/ember-cli-flash) used for form elements where options come dynamically from the API.

**Params**

| Param | Type | Description |
| --- | --- | --- |
| id | <code>String</code> | The name of the form field |
| models | <code>String</code> | An array of model types to fetch from the API. |
| onChange | <code>Func</code> | The onchange action for this form field. |
| inputValue | <code>String</code> \| <code>Array</code> | A comma-separated string or an array of strings. |
| [helpText] | <code>String</code> | Text to be displayed in the info tooltip for this form field |
| label | <code>String</code> | Label for this form field |
| fallbackComponent | <code>String</code> | name of component to be rendered if the API call 403s |
| options | <code>Array</code> | *Advanced usage* - `options` can be passed directly from the outside to the power-select component. If doing this, `models` should not also be passed as that will overwrite the passed value. |
| search | <code>Func</code> | *Advanced usage* - Customizes how the power-select component searches for matches - see the power-select docs for more information. |

**Example**
  
```js
<SearchSelect @id="group-policies" @models={{["policies/acl"]}} @onChange={{onChange}} @inputValue={{get model valuePath}} @helpText="Policies associated with this group" @label="Policies" @fallbackComponent="string-list" />
```

**See**

- [Uses of SearchSelect](https://github.com/hashicorp/vault/search?l=Handlebars&q=SearchSelect+OR+search-select)
- [SearchSelect Source Code](https://github.com/hashicorp/vault/blob/master/ui/lib/core/addon/components/search-select.js)

---
