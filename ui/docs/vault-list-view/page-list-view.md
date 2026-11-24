# Page::ListView and ListTable

`Page::ListView` is the shared structural wrapper for all Vault list views. It owns the page header, breadcrumbs, toolbar, table, pagination, empty states, row actions popup menu, and confirm modal. Individual list-view routes pass a static config object and dynamic data — they do not manage layout or common controls themselves.

`ListTable` and `ListTable::Cell` are the lower-level components that `Page::ListView` delegates to for table rendering. They can be used independently when you need a paginated table without the full page wrapper, but the common case is `Page::ListView`.

## Architecture

```
Page::ListView
├── Page::Header        (title, breadcrumbs, description, primary action)
├── Toolbar             (filter input, toolbar actions)
├── ListTable           (columns, rows, pagination)
│   └── ListTable::Cell (per-cell value rendering by valueType)
├── Hds::ApplicationState  (no-data and filtered-empty states)
└── Hds::Modal          (built-in confirm dialog for destructive row actions)
```

The split between `@config` (static) and `@data` / `@pageFilter` / `@page` (dynamic) is intentional. Glimmer tracks them independently so the header and breadcrumbs do not re-render when only the data or filter changes.

## Basic usage

Define the config as a class property — not a getter — so the object reference is stable across renders:

```js
// controller.js
import Controller from '@ember/controller';
import { tracked } from '@glimmer/tracking';

export default class SecretsBackendsController extends Controller {
  // page triggers a model refresh (API call). pageSize is client-side only.
  queryParams = ['page', 'pageSize'];
  @tracked page = 1;
  @tracked pageSize = 10;
}
```

```hbs
{{! template.hbs }}
<Page::ListView
  @config={{this.model.listViewConfig}}
  @model={{this.model.backends}}
  @page={{this.page}}
  @pageSize={{this.pageSize}}
/>
```

## Args

| Arg | Type | Required | Description |
|-----|------|----------|-------------|
| `@config` | `ListViewDisplayConfig` | ✓ | Static display config — title, breadcrumbs, columns, empty-state copy, optional primary action and filter. Define as a class property for a stable reference. |
| `@model` | `unknown[]` | ✓ | Raw (unpaginated) data array from the route model. `Page::ListView` applies filtering and pagination internally. |
| `@page` | `number` | — | Current page number from the controller query param. Passed to `paginate()` to slice the correct page. |
| `@pageSize` | `number` | — | Current page size from the controller query param. Must be passed so `paginate()` uses the same size as the HDS pagination widget. Defaults to `10` when absent. |

## Config: ListViewDisplayConfig

Define a config constant in `ui/app/utils/constants/list-view-config/<resource>.ts` and import it into the controller.

```ts
import type { ListViewDisplayConfig } from 'core/types/list-view-config';

export const MY_RESOURCE_LIST_VIEW_CONFIG: ListViewDisplayConfig = {
  title: 'My Resources',
  breadcrumbs: [
    { label: 'Vault', icon: 'vault-color', route: 'vault.cluster' },
    { label: 'My Resources' },
  ],
  columns: [
    { key: 'name', label: 'Name', valueType: 'link', route: 'vault.cluster.my-resource.details', isExpandable: true },
    { key: 'status', label: 'Status', valueType: 'status', statusMap: { active: 'success', inactive: 'warning' } },
    { key: 'created_at', label: 'Created', valueType: 'date-time' },
    { key: 'popupMenu', label: '' },
  ],
  noDataTitle: 'No resources yet',
  noDataDescription: 'Create a resource to get started.',
  filteredEmptyTitle: 'No resources match',
  primaryAction: { label: 'Create', route: 'vault.cluster.my-resource.create', icon: 'plus' },
  filter: { type: 'text', placeholder: 'Filter resources', ariaLabel: 'Filter resources' },
  rowActions: [
    { label: 'Edit', dataTest: 'edit', kind: 'route', route: 'vault.cluster.my-resource.edit', modelKey: 'id' },
    {
      label: 'Delete', dataTest: 'delete', kind: 'modal', color: 'critical',
      modal: {
        title: 'Delete resource?',
        body: 'This will permanently delete ',
        itemDisplayKey: 'name',
        confirmActionName: 'deleteResource',
        color: 'critical',
      },
    },
  ],
};
```

### Column value types

Set `column.valueType` to control how a cell renders. No named block is needed for built-in types.

| `valueType` | Additional config fields | Renders |
|-------------|--------------------------|---------|
| _(omitted)_ | — | Plain text. Objects are JSON-stringified. |
| `"link"` | `route`, `routeKey`, `fallbackRouteKey`, `routeModelKey` (default `"id"`), `iconKey`, `icon` | `Hds::Link::Standalone` (when icon present) or `Hds::Link::Inline`. Route resolves: `routeKey` → `fallbackRouteKey` → `route`. |
| `"icon-text"` | `iconKey`, `icon` (fallback `"lock"`), `tooltipKey`, `tooltip`, `textKey` | `Hds::Icon` + optional `Hds::TooltipButton` + text span. |
| `"copy"` | — | `Hds::Copy::Button` (icon-only) + plain text. Hidden when value is falsy. |
| `"header-tooltip"` | `headerTooltip` | Plain text cell; ⓘ tooltip on the column header. Set `headerTooltip` on the column — `Page::ListView` copies it to the HDS-native `tooltip` field at render time. |
| `"status"` | `statusMap`, `statusIconMap` | `Hds::Badge` — maps value to color and icon via the maps. Unknown values render a neutral badge. |
| `"date-time"` | — | `<time>` wrapping `{{date-format}}` output, format `"MMM dd, yyyy h:mm a"`. Hidden when value is falsy. |
| `"custom"` | `customTableItem: true` (legacy alias) | Yields `<:customTableItem as \|row col val\|>` to the calling template. |

### Row actions

`config.rowActions` drives the `···` popup menu automatically. Each action has a `kind`:

| `kind` | Required fields | What happens on click |
|--------|-----------------|-----------------------|
| `"route"` | `route` or `routeKey` | Navigates to the resolved route with `@model` from `row[modelKey ?? "id"]`. |
| `"action"` | `actionName` | Calls `@actions[action.actionName](rowData)`. `@actions` must be passed. |
| `"modal"` | `modal` (a `ListViewModalConfig`) | Opens the built-in confirm modal. On confirm, calls `@actions[modal.confirmActionName](item, closeModal)`. |

The `<:popupMenu as |rowData|>` named block overrides the config-driven menu when present.

## Named blocks

All named blocks are optional escape hatches. Use config-driven rendering for common cases.

| Block | Yields | Purpose |
|-------|--------|---------|
| `<:headerActions>` | — | Additional buttons in the page header beyond `config.primaryAction`. |
| `<:toolbarActions>` | — | Right-side toolbar controls (e.g. sort toggle, bulk action button). |
| `<:popupMenu as \|rowData\|>` | `rowData` | Fully custom `···` menu per row. Overrides `config.rowActions` entirely. |
| `<:customTableItem as \|row col val\|>` | `row`, `col`, `val` | Custom cell rendering for columns with `valueType: "custom"`. |
| `<:confirmModal as \|item config onClose\|>` | `activeModalItem`, `activeModalConfig`, `closeModal` | Custom confirm modal. Overrides the built-in `Hds::Modal`. |

## Empty states

`Page::ListView` renders one of three states based on `@data` and `@pageFilter`:

| State | Condition | Rendered from |
|-------|-----------|---------------|
| No data | `pageFilter` is falsy and `data` is empty | `config.noDataTitle` + `config.noDataDescription` |
| Filtered empty | `pageFilter` is set and `data` is empty | `config.filteredEmptyTitle` + the current filter value |
| Table | Otherwise | `ListTable` with configured columns |

## ListTable

`ListTable` wraps `Hds::AdvancedTable` and adds `Hds::Pagination::Numbered`. Use it directly when you need a standalone paginated table without the full page wrapper.

```hbs
<ListTable @columns={{this.columns}} @data={{@data}} @childrenKey="children">
  <:popupMenu as |rowData|>
    {{! custom menu }}
  </:popupMenu>
</ListTable>
```

| Arg | Description |
|-----|-------------|
| `@columns` | Column definitions array (same shape as `config.columns`). |
| `@data` | Full unsliced data array — `ListTable` handles pagination internally. |
| `@childrenKey` | Key on each row whose value is an array of child rows (for expandable rows). |
| `@selectionKeyField` | When set, enables row selection and passes a unique key per row to `Hds::AdvancedTable`. |
| `@hidePagination` | Hides the pagination bar (useful for small fixed-size tables). |

## TypeScript types

All types are exported from `core/types/list-view-config`:

```ts
import type {
  ListViewDisplayConfig,   // @config shape for Page::ListView
  ListViewColumn,          // single column definition
  ListViewRowAction,       // single row action definition
  ListViewModalConfig,     // confirm modal configuration
  ListViewFilter,          // filter input configuration
  ListViewBreadcrumb,      // breadcrumb item
  ListViewConfig,          // full config used by the vault-list-view Bob skill
} from 'core/types/list-view-config';
```

## Related

- [Vault list-view skill](./skill.md) — the Bob skill that generates list view scaffolding from a screenshot or Figma URL
- [Client pagination](../client-pagination.md) — how `paginate()` works with `@data`
- [`ListViewDisplayConfig` type](../../lib/core/addon/types/list-view-config.ts)
- [`Page::ListView` component](../../lib/core/addon/components/page/list-view.ts)
- [`ListTable::Cell` component](../../lib/core/addon/components/list-table/cell.ts)
- [Secrets backends config (reference implementation)](../../app/utils/constants/list-view-config/secrets-backends.ts)
