Migration Guides
================

## Migrating to Version 2

### Path

Previously it was not possible to specify a dot path where a key itself contains a dot. In v2 it is now possible with the escape sequence `~1`. For example, given the JSON doc `{"foo":{"bar.baz":10}}`, the path `foo.bar~1baz` would return `10`. This escape sequence means the character `~` is also a special case, therefore it must also be escaped to the sequence `~0`.

### Consume

Calls to `Consume(root interface{}) (*Container, error)` should be replaced with `Wrap(root interface{}) *Container`.

The error response was removed in order to avoid unnecessary duplicate type checks on `root`. This also allows shorthand chained queries like `gabs.Wrap(foo).S("bar","baz").Data()`.

### Search Across Arrays

All query functions (`Search`, `Path`, `Set`, `SetP`, etc) now attempt to resolve a specific index when they encounter an array. This means path queries must specify an integer index at the level of arrays within the content.

For example, given the sample document:

``` json
{
  "foo": [
    {
      "bar": {
        "baz": 45
      }
    }
  ]
}
```

In v1 the query `Search("foo", "bar", "baz")` would propagate the array in the result giving us `[45]`. In v2 we can access the field directly with `Search("foo", "0", "bar", "baz")`. The index is _required_, otherwise the query fails.

In query functions that do not set a value it is possible to specify `*` instead of an index in order to obtain all elements of the array, this produces the equivalent result as the behaviour from v1. For example, in v2 the query `Search("foo", "*", "bar", "baz")` would return `[45]`.

### Children and ChildrenMap

The `Children` and `ChildrenMap` methods no longer return errors. Instead, in the event of the underlying value being invalid (not an array or object), a `nil` slice and empty map are returned respectively. If explicit type checking is required the recommended approach would be casting on the value, e.g. `foo, ok := obj.Data().([]interface)`.

### Serialising Invalid Types

In v1 attempting to serialise with `Bytes`, `String`, etc, with an invalid structure would result in an empty object `{}`. This behaviour was unintuitive and in v2 `null` will be returned instead. If explicit marshalling is required with proper error propagation it is still recommended to use the `json` package directly on the underlying value.