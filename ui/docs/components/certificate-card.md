# CertificateCard

The CertificateCard component receives data and optionally receives a boolean declaring if that data is meant to be in PEM
Format. It renders using the `<HDS::Card::Container>`. To the left there is a certificate icon. In the center there is a label
which says which format (PEM or DER) the data is in. Below the label is the truncated data. To the right there is a copy
button to copy the data.

| Param | Type                 | Description                                                                                                   |
| ----- | -------------------- | ------------------------------------------------------------------------------------------------------------- |
| data  | <code>string</code>  | the data to be displayed in the component (usually in PEM or DER format)                                      |
| isPem | <code>boolean</code> | optional argument for if the data is required to be in PEM format (and should thus have the PEM Format label) |

**Example**

```hbs preview-template
<CertificateCard
  @data='-----BEGIN CERTIFICATE-----\nMIIE7TCCA9WgAwIBAgIULcrWXSz3/kG81EgBo0A4Zt'
  @isPem={{true}}
/>
```
