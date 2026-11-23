# Client-side pagination

LIST endpoints return every key in one response, so lists are paginated in the browser. A route fetches the full list from the API service and passes it to `paginate`, which filters and slices it and attaches the `meta` the pagination components read.

## Using pagination

```js
import { paginate } from 'core/utils/paginate-list';
import { SecretsApiDatabaseListConnectionsListEnum } from '@hashicorp/vault-client-typescript';

export default class ExampleRoute extends Route {
  @service api;

  async model(params) {
    const { page, pageFilter } = params;
    const { backend } = this.paramsFor('vault.cluster.secrets.backend');
    const { keys = [] } = await this.api.secrets.databaseListConnections(
      backend,
      SecretsApiDatabaseListConnectionsListEnum.TRUE
    );
    // filterKey names the property to match for lists of objects; omit it for lists of strings
    return paginate(
      keys.map((name) => ({ id: name, name, backend })),
      { page, filter: pageFilter, filterKey: 'name' }
    );
  }
}
```

`paginate(data, { page, pageSize, filter, filterKey })`:

- `page` defaults to 1. A page beyond the last page falls back to the first.
- `pageSize` defaults to the default page size set in [the app config](../config/environment.js).
- `filter` keeps items that contain the text, ignoring case. `filterKey` names the property to match, or an array of properties where any may match.

The result is the page of items with a read-only `meta` property: `currentPage`, `lastPage`, `nextPage`, `prevPage`, `total`, `filteredTotal`, `totalItems` and `pageSize`. Pass `meta.totalItems` to `Hds::Pagination::Numbered @totalItems` since it reflects the active filter. See [paginate-list.ts](../lib/core/addon/utils/paginate-list.ts).

## Capabilities for a page

Fetch capabilities only for the items on the current page, in one request. See [Capabilities](./capabilities.md).

## Gotchas

The full list is fetched each time the route's model hook runs, so a list is current after an item is created or deleted as long as the route refreshes. `paginate` assigns `meta` to the returned array, so assign properties onto its items in place (for example capabilities) rather than mapping to a new array.
