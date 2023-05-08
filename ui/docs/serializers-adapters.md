## Guidelines

- Prepend internal functions with an underscore to differentiate from Ember methods `_getUrl`
- Consider using the [named-path](../app/adapters/named-path.js) adapter if the model name is part of the request path
- Utilize the serializer to remove sending model attributes that do not correspond to an API parameter. Example in [key serializer](../app/serializers/pki/key.js)

```js
export default class SomeSerializer extends ApplicationSerializer {
  attrs = {
    attrName: { serialize: false },
  };
}
```

> Note: this will remove the attribute when calling `snapshot.serialize()` method regardless of any custom serializer logic.

## Gotchas

- The JSON serializer removes attributes with empty arrays [Example in MFA serializer](https://github.com/hashicorp/vault/blob/e55c18ed1299e0d36b88e603fa9f12adaf8e75dc/ui/app/serializers/mfa-login-enforcement.js#L37-L44)
