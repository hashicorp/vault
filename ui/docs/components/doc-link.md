# DocLink

&#x60;DocLink&#x60; components are used to render anchor links to relevant Vault documentation at developer.hashicorp.com.

| Param | Type                | Default                                  | Description                                                                             |
| ----- | ------------------- | ---------------------------------------- | --------------------------------------------------------------------------------------- |
| path  | <code>string</code> | <code>&quot;\&quot;/\&quot;&quot;</code> | The path to documentation on developer.hashicorp.com that the component should link to. |

**Example**

```hbs preview-template
<DocLink @path='/vault/docs/secrets/kv/kv-v2.html'>Learn about KV v2</DocLink>
```
