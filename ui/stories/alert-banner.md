<a name="AlertBanner
`AlertBanner` components are used to inform users of important messages.module_"></a>

## AlertBanner
`AlertBanner` components are used to inform users of important messages.

**Properties**

| Name | Type | Default | Description |
| --- | --- | --- | --- |
| [AlertBanner.type] | <code>String</code> | <code></code> | The banner type. Should either be `info`, `warning`, `success`, or `danger`. |
| [AlertBanner.message] | <code>String</code> | <code></code> | The message to display within the banner. |

**Example** 

```js
<AlertBanner @type="danger" @message="{{model.keyId}} is not a valid lease ID"/>
```