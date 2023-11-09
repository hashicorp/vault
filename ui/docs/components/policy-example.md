# PolicyExample

The PolicyExample component receives a policy type (acl, rgp, or egp) and renders a copyable policy example of
that type using the `<JsonEditor>` component. Inside a modal, the PolicyExample component must be wrapped in a conditional
(example below), otherwise the `<JsonEditor>` value wont render until its focused.

| Param      | Type                | Description                                                                                    |
| ---------- | ------------------- | ---------------------------------------------------------------------------------------------- |
| policyType | <code>string</code> | policy type to decide which template to render; can either be "acl" or "rgp"                   |
| container  | <code>string</code> | selector for the container the example renders inside, passed to the copy button in JsonEditor |

**Example**

```hbs preview-template
<PolicyExample @policyType='acl' @container='#search-select-modal' />
```
