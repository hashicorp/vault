
# NavigateInput
&#x60;NavigateInput&#x60; components are used to filter list data.

| Param | Type | Default | Description |
| --- | --- | --- | --- |
| filter | <code>String</code> | <code></code> | The filtered string. |
| [placeholder] | <code>String</code> | <code>&quot;Filter items&quot;</code> | The message inside the input to indicate what the user should enter into the space. |
| [urls] | <code>Object</code> | <code></code> | An object containing list=route url. |
| [filterFocusDidChange] | <code>function</code> | <code></code> | A function called when the focus changes. |
| [filterDidChange] | <code>function</code> | <code></code> | A function called when the filter string changes. |
| [filterMatchesKey] | <code>function</code> | <code></code> | A function used to match to a specific key, such as an Id. |
| [filterPartialMatch] | <code>function</code> | <code></code> | A function used to filter through a partial match. Such as "oo" of "root". |
| [baseKey] | <code>String</code> | <code>&quot;&quot;</code> | A string to transition by Id. |
| [shouldNavigateTree] | <code>Boolean</code> | <code>false</code> | If true, navigate a larger tree, such as when you're navigating leases under access. |
| [mode] | <code>String</code> | <code>&quot;secrets&quot;</code> | Mode which plays into navigation type. |
| [extraNavParams] | <code>String</code> | <code>&quot;&quot;</code> | A string used in route transition when necessary. |

**Example**  
```hbs preview-template
<NavigateInput @filter={@roleFiltered} @placeholder="placeholder text" urls="{{hash list="vault.cluster.secrets.backend.kubernetes.roles"}}"/>
```
