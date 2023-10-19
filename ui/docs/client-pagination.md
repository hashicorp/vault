# Client-side pagination

Our custom extended `store` service allows us to paginate LIST responses while maintaining good performance, particularly when the LIST response includes tens of thousands of keys in the data response. It does this by caching the entire response, and then filtering the full response into the datastore for the client.

## Using pagination

Rather than use `store.query`, use `store.lazyPaginatedQuery`. It generally uses the same inputs, but accepts additional keys in the query object `size`, `page`, `responsePath`, `pageFilter`

### Before

```js
export default class ExampleRoute extends Route {
  @service store;

  model(params) {
    const { secret } = params;
    const { backend } = this.paramsFor('vault.cluster.secrets.backend');
    return this.store.query('pki/role', { backend, id })
  }
```

### After

```js
export default class ExampleRoute extends Route {
  @service store;

  model(params) {
    const { page, pageFilter, secret } = params;
    const { backend } = this.paramsFor('vault.cluster.secrets.backend');
    return this.store.lazyPaginatedQuery('secret', {
      backend,
      id: secret,
      size,
      page,
      responsePath,
      pageFilter
    })
  }
```

The `size` param defaults to the default page size set in [the app config](../config/environment.js). `responsePath` and `page` are required, and typically `responsePath` is going to be `data.keys` since that is where the LIST responses typically return their array data.

### Serializing

In order to interrupt the regular serialization when using `lazyPaginatedData`, define `extractLazyPaginatedData` on the modelType's serializer. This will be called with the raw response before being cached on the store. `extractLazyPaginatedData` should return an array of objects.

## Gotchas

The data is cached from whenever the original API call is made, which means that if a user views a list and then creates or deletes an item, viewing the list page again will show outdated information unless the cache for the item is cleared first. For this reason, it is best practice to clear the dataset with `store.clearDataset(modelName)` after successfully deleting or creating an item.

## How it works

When using the `lazyPaginatedQuery` method, the full response is cached in a [tracked Map](https://github.com/tracked-tools/tracked-built-ins/tree/master) within the service. `store.lazyCaches` is actually a Map of [Maps](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Map), keyed first on the normalized modelType and then on a stringified version of the base query (all keys except ones related to pagination). So, at the top level `store.lazyCaches` looks like this:

```
lazyCaches = new Map({
  "secret": <Map>,
  "kmip": <Map>,
  "namespace": <Map>,
})
```

Within each top-level modelType, we need to separate cached responses based on the details of the query. Typically (but not always) this includes the backend name. In list items that can be nested (see KV V2 secrets or namespaces for example) `id` is also provided, so that the keys nested under the given ID is returned. The store.lazyCaches may look something like the following after a user navigates to a couple different KV v2 lists, and clicks into the `app/` item:

```
lazyCaches = new Map({
  "secret": {
    "{ backend: 'secret', id: '' }: <CachedData>,
    "{ backend: 'secret', id: 'app/' }: <CachedData>,
    "{ backend: 'kv2', id: '' }: <CachedData>,
  },
  ...
})
```

The cached data at the given key is an object with `response` and `dataset` keys. The response is the full response from the original API call, with the `responsePath` nulled out (it is repopulated before "sending" the data back to the store). `dataset` is the full, original value at `responsePath`, usually an array of strings. An example of what the data might look like:

```
lazyCaches = new Map({
  "secret": {
    "{ backend: 'secret', id: 'app/' }: {
      dataset: ['some', 'nested', 'secrets'],
      response: {
        request_id: 'foobar',
        data: {},
        ...
      }
    },
  },
  ...
})
```
