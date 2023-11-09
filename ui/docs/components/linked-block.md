
# LinkedBlock
LinkedBlock components are linkable divs that yield any content nested within them. They are often used in list views such as when listing the secret engines.

| Param | Type | Default | Description |
| --- | --- | --- | --- |
| params | <code>Array</code> | <code></code> | These are values sent to the router's transitionTo method.  First item is route, second is the optional path. |
| [queryParams] | <code>Object</code> | <code></code> | queryParams can be passed via this property. It needs to be an object. |
| [linkPrefix] | <code>String</code> | <code></code> | Overwrite the params with custom route.  Needed for use in engines (KMIP and PKI). ex: vault.cluster.secrets.backend.kmip |
| [encode] | <code>Boolean</code> | <code>false</code> | Encode the path. |
| [disabled] | <code>boolean</code> |  | disable the link -- prevents on click and removes linked-block hover styling |

**Example**  
```hbs preview-template
<LinkedBlock @params={{array 'vault.cluster.secrets.backend.show' 'my-secret-path'}} @queryParams={{hash version=1}} @class="list-item-row">
 Yield
</LinkedBlock>
```
