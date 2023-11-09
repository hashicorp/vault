
# AlertInline
AlertInline components are used to inform users of important messages.

| Param | Type | Default | Description |
| --- | --- | --- | --- |
| type | <code>string</code> | <code>null</code> | The alert type passed to the message-types helper. |
| [message] | <code>string</code> | <code>null</code> | The message to display within the alert. |
| [paddingTop] | <code>boolean</code> | <code>false</code> | Whether or not to add padding above component. |
| [isMarginless] | <code>boolean</code> | <code>false</code> | Whether or not to remove margin bottom below component. |
| [sizeSmall] | <code>boolean</code> | <code>false</code> | Whether or not to display a small font with padding below of alert message. |
| [mimicRefresh] | <code>boolean</code> | <code>false</code> | If true will display a loading icon when attributes change (e.g. when a form submits and the alert message changes). |

**Example**  
```hbs preview-template
<AlertInline @type="danger" @message="{{model.keyId}} is not a valid lease ID" />
```
