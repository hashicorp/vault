
# KeyValueHeader
KeyValueHeader components show breadcrumbs for secret engines.

| Param | Type | Default | Description |
| --- | --- | --- | --- |
| [mode] | <code>string</code> | <code>null</code> | Used to set the currentPath. |
| [baseKey] | <code>string</code> | <code>null</code> | Used to generate the path backward. |
| [path] | <code>string</code> | <code>null</code> | The fallback path. |
| [root] | <code>string</code> | <code>null</code> | Used to set the secretPath. |
| [showCurrent] | <code>boolean</code> | <code>true</code> | Boolean to show the second part of the breadcrumb, ex: the secret's name. |
| [linkToPaths] | <code>boolean</code> | <code>true</code> | If true link to the path. |

**Example**  
```hbs preview-template
<KeyValueHeader @path="vault.cluster.secrets.backend.show" @mode={{this.mode}} @root={{@root}}/>
```
