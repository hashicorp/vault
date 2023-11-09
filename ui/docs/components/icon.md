# Icon

&#x60;Icon&#x60; components are glyphs used to indicate important information.

Flight icon documentation at https://flight-hashicorp.vercel.app/

| Param  | Type                | Default           | Description                           |
| ------ | ------------------- | ----------------- | ------------------------------------- |
| name   | <code>string</code> | <code>null</code> | The name of the SVG to render inline. |
| [size] | <code>string</code> | <code>16</code>   | size for flight icon, can be 16 or 24 |

**Example**

```hbs preview-template
<Icon @name='cancel-square-outline' @size='24' />
```
