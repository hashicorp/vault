# MessageError

Renders form errors using the Hds::Alert component and extracts errors from a model, passed errorMessage or array of errors and displays each in a separate banner.

| Param          | Type                | Description                                                                                   |
| -------------- | ------------------- | --------------------------------------------------------------------------------------------- |
| [model]        | <code>object</code> | An Ember data model that will be used to bind error states to the internal `errors` property. |
| [errors]       | <code>array</code>  | An array of error strings to show.                                                            |
| [errorMessage] | <code>string</code> | An Error string to display.                                                                   |

**Example**

```hbs preview-template
<MessageError @errorMessage='there is something very wrong' />
```
