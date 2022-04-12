<!--THIS FILE IS AUTO GENERATED. This file is generated from JSDoc comments in lib/core/addon/components/search-select.js. To make changes, first edit that file and run "yarn gen-story-md search-select" to re-generate the content.-->

## SearchSelect
The `SearchSelect` is an implementation of the [ember-power-select](https://github.com/cibernox/ember-power-select) used for form elements where options come dynamically from the API.

**Params**

| Param | Type | Default | Description |
| --- | --- | --- | --- |
| id | <code>string</code> |  | The name of the form field |
| models | <code>Array</code> |  | An array of model types to fetch from the API. |
| onChange | <code>function</code> |  | The onchange action for this form field. |
| inputValue | <code>string</code> \| <code>Array</code> |  | A comma-separated string or an array of strings. |
| label | <code>string</code> |  | Label for this form field |
| fallbackComponent | <code>string</code> |  | name of component to be rendered if the API call 403s |
| [backend] | <code>string</code> |  | name of the backend if the query for options needs additional information (eg. secret backend) |
| [disallowNewItems] | <code>boolean</code> | <code>false</code> | Controls whether or not the user can add a new item if none found |
| [helpText] | <code>string</code> |  | Text to be displayed in the info tooltip for this form field |
| [selectLimit] | <code>number</code> |  | A number that sets the limit to how many select options they can choose |
| [subText] | <code>string</code> |  | Text to be displayed below the label |
| [subLabel] | <code>string</code> |  | a smaller label below the main Label |
| [wildcardLabel] | <code>string</code> |  | when you want the searchSelect component to return a count on the model for options returned when using a wildcard you must provide a label of the count e.g. role.  Should be singular. |
| options | <code>Array</code> |  | *Advanced usage* - `options` can be passed directly from the outside to the power-select component. If doing this, `models` should not also be passed as that will overwrite the passed value. |
| search | <code>function</code> |  | *Advanced usage* - Customizes how the power-select component searches for matches - see the power-select docs for more information. |

**Example**
  
```js
<SearchSelect @id="group-policies" @models={{["policies/acl"]}} @onChange={{onChange}} @selectLimit={{2}} @inputValue={{get model valuePath}} @helpText="Policies associated with this group" @label="Policies" @fallbackComponent="string-list" />
```

**See**

- [Uses of SearchSelect](https://github.com/hashicorp/vault/search?l=Handlebars&q=SearchSelect+OR+search-select)
- [SearchSelect Source Code](https://github.com/hashicorp/vault/blob/main/ui/lib/core/addon/components/search-select.js)

---
