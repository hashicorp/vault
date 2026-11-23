# Capabilities

The `capabilities` service answers "can the current token do X on this path?". It replaces the `lazyCapabilities` model macros, which sent one request per path and per record.

- [Defining paths](#defining-paths)
- [Fetching capabilities](#fetching-capabilities)
- [Where to fetch](#where-to-fetch)
- [Errors](#errors)

## Defining paths

Paths live in `PATH_MAP` in [capabilities.ts](../app/utils/constants/capabilities.ts), keyed by a friendly name and built with the `apiPath` tagged template. Dynamic segments are named, and values are filled in by `pathFor`.

```js
// app/utils/constants/capabilities.ts
authMethodDelete: apiPath`sys/auth/${'path'}`,

// anywhere the service is injected
this.capabilities.pathFor('authMethodDelete', { path: 'userpass/' }); // 'sys/auth/userpass/'
```

## Fetching capabilities

- `fetch(paths)` makes one request for any number of paths and returns a map of `path -> { canCreate, canDelete, canList, canPatch, canRead, canSudo, canUpdate }`. Prefer this for lists so every row shares a single request.
- `for(key, params)` fetches the capabilities of one path from its `PATH_MAP` key.
- `canRead(path)`, `canUpdate(path)` and `canPatch(path)` resolve to a single boolean.

```js
// route model hook for a list view
const paths = items.map((item) => this.capabilities.pathFor('authMethodDelete', { path: item.id }));
const capabilities = await this.capabilities.fetch(paths);
return items.map((item, index) => ({ ...item, canDelete: capabilities[paths[index]].canDelete }));
```

In templates, the `has-capability` helper reads a capabilities map for a `PATH_MAP` key:

```hbs
(has-capability @model.capabilities "delete" pathKey="authMethodDelete" params=method.id)
```

## Where to fetch

Fetch capabilities in the route model hook, or on a user action such as opening a dropdown, and pass the result down as plain data. Do not fetch them from models or serializers. Capabilities are not cached between requests, so a view that needs the same paths again fetches them again.

Pass `routeForCache` in the fetch options when the route should prefill the policy flyout (`<CodeGeneratorPolicyFlyout />`) with the paths it checked.

## Errors

If the capabilities request fails, `fetch` resolves every path to `true` and lets the API enforce access. Pass `{ strict: true }` to throw instead.
