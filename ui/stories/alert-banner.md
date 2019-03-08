# AlertBanner

`AlertBanners` are used to inform users of important messages.

## Properties
| Property | Required | Default value | Type | Description | Example |
|---|---|---|---|---|
| type | [x] | `null` | string | The banner type. Should either be `info`, `warning`, `success`, or `danger`. | `info` |
| message ||| string | The message to display within the banner. | `Hello!` |

## Usage

```javascript
<AlertBanner
    @type="danger"
    @message="{{model.keyId}} is not a valid lease ID"
  />
```
https://github.com/hashicorp/vault/search?l=Handlebars&q=AlertBanner

## Source
https://github.com/hashicorp/vault/blob/master/ui/app/components/alert-banner.js
