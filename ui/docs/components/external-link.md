
# ExternalLinkComponent
&#x60;ExternalLink&#x60; components are used to render anchor links to non-cluster links. Automatically opens in a new tab with noopener noreferrer.
 To link to developer.hashicorp.com, use DocLink .

| Param | Type | Default | Description |
| --- | --- | --- | --- |
| href | <code>string</code> | <code>&quot;\&quot;https://example.com/\&quot;&quot;</code> | The full href with protocol |
| [sameTab] | <code>boolean</code> | <code>false</code> | by default, these links open in new tab. To override, pass @sameTab={{true}} |

**Example**  
```hbs preview-template
<ExternalLink @href="https://hashicorp.com">Arbitrary Link</ExternalLink>
```
