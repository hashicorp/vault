# AlertInline

`AlertInline`s are used to inform users of important messages.

## Properties
| Property | Required | Default value | Type | Description | Example |
|---|---|---|---|---|
| type | [x] | `null` | string | The alert type. Should either be `info`, `warning`, `success`, or `danger`. | `info` |
| message ||| string | The message to display within the alert. | `Hello!` |

## Usage

```javascript
<AlertInline
  @type="danger"
  @message="Demoting this DR primary cluster would result in a DR secondary."
/>
```
https://github.com/hashicorp/vault/search?l=Handlebars&q=AlertInline

## Source
https://github.com/hashicorp/vault/blob/master/ui/app/components/alert-inline.js
