/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

/**
 * @module ListViewConfig
 *
 * Formal TypeScript contract for Vault list views.
 * Consumed by the `vault-list-view` Bob skill (SKILL.md) and by generated
 * page components that delegate rendering to `Page::ListView`.
 *
 * Import path:  import type { ListViewConfig } from 'vault/list-view-config';
 */

// ── Sub-interfaces (exported individually so consumers can import only what they need) ──

export interface ListViewBreadcrumb {
  label: string;
  route?: string;
  icon?: string;
  model?: string;
  models?: string[];
  linkExternal?: boolean;
}

// ── Column value types ────────────────────────────────────────────────────────
//
// Declares how Page::ListView renders each cell automatically, without the
// consuming page component needing to yield a <:customTableItem> block.
//
// | valueType          | Renders                                              |
// |--------------------|------------------------------------------------------|
// | "text" (default)   | Plain string value — same as omitting valueType      |
// | "link"             | Hds::Link::Inline navigating to column.route using   |
// |                    | rowData[column.routeModelKey ?? "id"] as @model      |
// | "icon-text"        | Hds::Icon + plain text; icon name from               |
// |                    | rowData[column.iconKey] (falls back to column.icon)  |
// |                    | with an optional tooltip via column.tooltip or       |
// |                    | rowData[column.tooltipKey]                           |
// | "copy"             | Hds::Copy::Button (icon-only) + plain text           |
// | "header-tooltip"   | Plain text with a ⓘ tooltip on the column HEADER;   |
// |                    | tooltip text set via column.headerTooltip             |
// | "status"           | Hds::Badge; maps value to @color via column.statusMap|
// |                    | e.g. { Enabled: "success", Disabled: "warning" }     |
// | "date-time"        | {{date-format rowData[key] "MMM dd, yyyy h:mm a"}} wrapped in <time> |
// | "custom"           | Falls through to the <:customTableItem> yield block  |
// |                    | (same as setting customTableItem: true)              |

export type ListViewColumnValueType =
  | 'text'
  | 'link'
  | 'icon-text'
  | 'copy'
  | 'header-tooltip'
  | 'status'
  | 'date-time';

export interface ListViewColumn {
  key: string;
  label: string;
  width?: string;
  isExpandable?: boolean;
  isSortable?: boolean;

  /**
   * Declares how the cell is rendered by Page::ListView automatically.
   * Defaults to "text" when omitted.
   * Use "custom" (or the legacy `customTableItem: true`) when the consuming
   * page component must yield a <:customTableItem> block for this column.
   */
  valueType?: ListViewColumnValueType;

  /** @deprecated Use valueType: "custom" instead. Kept for backwards compat. */
  customTableItem?: boolean;

  // ── valueType: "link" ─────────────────────────────────────────────────────
  /**
   * Static Ember route name navigated to when the cell link is clicked.
   * Use when every row in the column links to the same route.
   * e.g. route: 'vault.cluster.access.identity.groups.group.details'
   */
  route?: string;
  /**
   * Row-data key whose value is the Ember route name for this cell.
   * Use when the destination route varies per row (e.g. different engine
   * types navigate to different routes).
   * Takes precedence over the static `route` field.
   * e.g. routeKey: 'backendLink' → reads rowData.backendLink at render time.
   *
   * Convention: always pair with a raw-data key on the row object.
   * If the value comes from a class getter, promote it to an own property
   * in the route's model() before paginating (same pattern as textKey).
   */
  routeKey?: string;
  /**
   * Row-data key whose value is used as the route when `routeKey` resolves
   * to a falsy value. Enables a conditional "primary route or fallback route"
   * pattern without requiring a <:customTableItem> block.
   * e.g. fallbackRouteKey: 'backendConfigurationLink'
   *   → supported backends use rowData.backendLink (routeKey)
   *   → unsupported backends fall through to rowData.backendConfigurationLink
   *
   * Resolution order: routeKey → fallbackRouteKey → route (static).
   */
  fallbackRouteKey?: string;
  /**
   * Key on the row data object whose value is passed as @model to the route.
   * Defaults to "id" when omitted.
   */
  routeModelKey?: string;
  /** Tooltip text displayed when the link has no route. */
  unavailableTooltip?: string;

  // ── valueType: "icon-text" ────────────────────────────────────────────────
  /**
   * Row-data key whose value is the icon name (e.g. "icon", "glyph").
   * Takes precedence over the static `icon` field.
   */
  iconKey?: string;
  /**
   * Static fallback icon name used when `iconKey` is absent or its value is
   * falsy (e.g. "lock", "user").
   */
  icon?: string;
  /**
   * Row-data key whose value is the tooltip text shown on the icon.
   * Takes precedence over the static `tooltip` field.
   */
  tooltipKey?: string;
  /** Static tooltip text used when `tooltipKey` is absent. */
  tooltip?: string;
  /**
   * Row-data key whose value is used as the display text for "icon-text" cells.
   * When omitted the cell falls back to rowData[column.key].
   * Use this when the display value comes from a computed getter rather than
   * the raw API field — e.g. key: "type", textKey: "engineType".
   */
  textKey?: string;

  // ── valueType: "header-tooltip" ───────────────────────────────────────────
  /** Tooltip text rendered as a ⓘ button next to the column header label. */
  headerTooltip?: string;

  // ── valueType: "status" ───────────────────────────────────────────────────
  /**
   * Maps cell values to Hds::Badge @color tokens.
   * e.g. { Enabled: "success", Disabled: "warning" }
   * Unrecognized values fall back to the default (neutral) badge color.
   */
  statusMap?: Record<string, string>;
  /**
   * Maps cell values to Hds::Badge @icon names.
   * e.g. { Enabled: "running-static", Disabled: "alert-triangle" }
   */
  statusIconMap?: Record<string, string>;
}

export interface ListViewModalConfig {
  /** Title / header for the modal (e.g. "Delete engine path?") */
  title: string;
  /** Key on the item object to append to the modal title */
  titleItemDisplayKey?: string;
  /** Optional icon for header, defaults to "alert-circle" */
  icon?: string;
  /** Modal body description or prefix text */
  body: string;
  /** Key on the item object to display in bold inside the body (e.g. "path" or "id") */
  itemDisplayKey?: string;
  /**
   * Optional bullet-list items rendered below the body text.
   * Each entry is a plain string; use `itemDisplayKey` for the dynamic item name in the first bullet.
   * e.g. ["secrets engine", "Engine configuration and third-party integrations"]
   * When present, the first entry is prefixed by the bold itemDisplayKey value (if set).
   */
  bodyItems?: string[];
  /** Optional secondary helper text under the body */
  bodySubtext?: string;
  /** Color/type: "critical" (red, default) or "warning" (yellow) */
  color?: 'critical' | 'warning';
  /** Confirm button text, defaults to "Confirm" (or "Delete" if color is critical) */
  confirmButtonText?: string;
  /** Cancel button text, defaults to "Cancel" */
  cancelButtonText?: string;
  /**
   * When set, renders a type-to-confirm text input inside the modal.
   * The user must type this exact string before the confirm button is enabled.
   * Passed as @confirmText to ConfirmModal.
   * e.g. "delete-engine"
   */
  confirmText?: string;
  /**
   * Label shown above the type-to-confirm input.
   * Passed as @confirmLabel to ConfirmModal.
   * Defaults to "Confirm" when omitted.
   * e.g. "Confirm deletion"
   */
  confirmLabel?: string;
  /**
   * Declarative API call executed when the user confirms the action.
   * Page::ListView resolves `api[service][method](item[argKey])` using its
   * injected ApiService — no function reference needed in the config.
   * e.g. { service: 'sys', method: 'mountsDisableSecretsEngine', argKey: 'id' }
   * Provide either `apiCall` or `deleteAction`, not both. `apiCall` takes precedence.
   */
  apiCall?: {
    /** Top-level key on ApiService, e.g. "sys", "identity", "secrets" */
    service: string;
    /** Method name on that service, e.g. "mountsDisableSecretsEngine" */
    method: string;
    /** Key on the row item whose value is passed as the sole argument, e.g. "id" */
    argKey: string;
  };
  /**
   * Escape hatch for cases where the API call cannot be expressed declaratively.
   * Prefer `apiCall` for standard single-argument API methods.
   * Page::ListView owns the flash message and router.refresh() wrapper.
   */
  deleteAction?: (item: Record<string, unknown>) => Promise<void>;
}

export interface ListViewRowAction {
  label: string;
  /** kebab-case of label — used as the `data-test-popup-menu` value */
  dataTest: string;
  kind: 'route' | 'action' | 'modal';
  /**
   * Static Ember route name. Use when every row links to the same route.
   * Required when kind is "route" and routeKey is absent.
   */
  route?: string;
  /**
   * Row-data key whose value is the Ember route name for this action.
   * Use when the destination route varies per row.
   * Takes precedence over the static `route` field.
   * Resolution order: routeKey → route.
   * [RULE] Same getter-promotion convention as column routeKey — promote
   * the getter to an own property in the route's model() spread.
   */
  routeKey?: string;
  /**
   * Row-data key whose value is passed as @model to the route.
   * Defaults to "id" when omitted.
   */
  modelKey?: string;
  /** Required when kind is "action" */
  actionName?: string;
  /** Modal configuration rendered when kind is "modal" */
  modal?: ListViewModalConfig;
  color?: 'critical';
  /** Optional capability key on the capabilities object (e.g. "canDelete") */
  capability?: string;
}

export interface ListViewBulkAction {
  label: string;
  actionName: string;
  color?: 'critical';
}

export interface ListViewFilter {
  type: 'text' | 'search-select';
  placeholder: string;
  ariaLabel: string;
  /**
   * Raw API field name (snake_case) to filter on. Must exactly match the field
   * on the data object — the same value that would be passed as `filterKey` to
   * paginate(). Required when filter is defined.
   * e.g. "path", "name", "id"
   */
  filterKey: string;
}

export interface ListViewRoutePaths {
  /**
   * Target route file path — always `.ts`.
   * e.g. "ui/app/routes/vault/cluster/access/identity/index.ts"
   */
  routeFile: string;
  /**
   * True when the current file on disk is a `.js` file.
   * When true, the skill writes a new `.ts` file and the old `.js` can be
   * deleted after the replacement is confirmed — giving JS→TS modernization
   * for free in the PR diff.
   */
  existingRouteIsJs: boolean;
  /** Always `.js` — e.g. "ui/app/controllers/vault/cluster/access/identity/index.js" */
  controllerFile: string;
  /** e.g. "ui/app/templates/vault/cluster/access/identity/index.hbs" */
  templateFile: string;
  /** Always `.js` — e.g. "ui/tests/acceptance/access/identity/groups/index-test.js" */
  testFile: string;
  /** "ui/app/router.js" for app-level routes, "ui/lib/<engine>/addon/routes.js" for engines */
  routerFile: string;
  /** Name of the parent this.route() block, e.g. "identity" */
  parentBlock: string;
  /** New route segment to register, e.g. "groups" */
  routeSegment: string;
  /** Only set if the URL path differs from the segment; defaults to "/" for index children */
  urlPath?: string;
}

// ── Static display config — passed as @config to Page::ListView ──────────────

/**
 * Static, per-route display configuration for Page::ListView.
 * Contains only fields that are fixed for a given route and do not change
 * during a render cycle: title, description, breadcrumbs, column definitions,
 * empty-state copy, and the primary action button.
 *
 * Dynamic runtime fields (data, pageFilter, page) are passed as separate
 * @args on Page::ListView so Glimmer can track them independently and avoid
 * re-rendering the static header/breadcrumbs when only the data changes.
 *
 * Filtering is standardized: if `filter` is defined in config, Page::ListView
 * renders FilterInput automatically and handles the query-param transition
 * internally — no callback or named block needed.
 *
 * Usage in a route template:
 *   <Page::ListView
 *     @config={{this.listViewConfig}}
 *     @data={{this.model.groups}}
 *     @pageFilter={{this.pageFilter}}
 *     @page={{this.model.page}}
 *   >
 *     <:popupMenu as |rowData|>…</:popupMenu>
 *   </Page::ListView>
 */
export interface ListViewDisplayConfig {
  /** Page title rendered in Page::Header */
  title: string;
  /** Optional badge rendered next to title */
  badge?: string;
  /** Optional badge icon rendered next to badge text */
  badgeIcon?: string;
  /** Optional description rendered below the title */
  description?: string;
  /** Breadcrumb items passed to Page::Breadcrumbs */
  breadcrumbs: ListViewBreadcrumb[];
  /** Column definitions passed to ListTable */
  columns: ListViewColumn[];
  /** Title for the no-data empty state */
  noDataTitle: string;
  /** Optional body text for the no-data empty state */
  noDataDescription?: string;
  /** Title prefix for the filtered-but-empty state (filter value is appended automatically) */
  filteredEmptyTitle: string;
  /** Key on each data row whose value is an array of child rows */
  childrenKey?: string;
  /**
   * Primary action button rendered in the Page::Header actions area.
   * When present, Page::ListView renders an Hds::Button automatically —
   * no <:headerActions> named block is needed for the common single-button case.
   */
  primaryAction?: {
    label: string;
    route: string;
    icon?: string;
  };
  /**
   * Filter input configuration. When present, Page::ListView renders
   * FilterInput automatically in the toolbar and transitions to
   * { queryParams: { pageFilter: value, page: 1 } } when the user types or clears.
   */
  filter?: ListViewFilter;
  /**
   * Row action definitions. When present, Page::ListView renders the ···
   * popup menu automatically for every row — no <:popupMenu> named block needed.
   * Each action is rendered as an Hds::Dropdown dd.Interactive item.
   * The <:popupMenu> named block remains available as an escape hatch for
   * non-standard cases (e.g. confirm modal triggers, capability-gated actions).
   * When both are present, the named block takes precedence.
   */
  rowActions?: ListViewRowAction[];
}

// ── Main config interface ──────────────────────────────────────────────────────────────

export interface ListViewConfig {
  // ── 0. Intent ───────────────────────────────────────────────────────────────
  /** "new" = route does not exist yet; "existing" = surgically modernise an existing route */
  intent: 'new' | 'existing';

  // ── 1. Routing ──────────────────────────────────────────────────────────────
  /** e.g. "vault.cluster.access.identity.index" */
  routeModule: string;
  /** All file paths derived from routeModule — see SKILL.md derivation rules */
  routePaths: ListViewRoutePaths;

  // ── 2. Header ───────────────────────────────────────────────────────────────
  header: {
    title: string;
    description?: string;
    secondaryAction?: {
      label: string;
      route: string;
      icon?: string;
    };
  };

  // ── 3. Breadcrumbs ──────────────────────────────────────────────────────────
  /** Must have at least 2 entries; the last entry must have no `route` (current page) */
  breadcrumbs: ListViewBreadcrumb[];

  // ── 4. Toolbar ──────────────────────────────────────────────────────────────
  toolbar: {
    hasSortToggle?: boolean;
    filtersDropdown?: boolean;
    filters?: ListViewFilter[];
    /** Required when a create/add button is visible in the header actions */
    primaryAction?: {
      label: string;
      route: string;
      icon?: string;
    };
    /** Empty array when no bulk actions are present */
    bulkActions?: ListViewBulkAction[];
  };

  // ── 5. Table ────────────────────────────────────────────────────────────────
  /** Value of the field used as the unique row identifier; set when rows are selectable */
  selectionKeyField?: string;
  /**
   * Key on each data row whose value is an array of child rows.
   * Passed as @childrenKey to Page::ListView → ListTable → Hds::AdvancedTable,
   * which renders expand/collapse natively.
   * Every parent object MUST carry this key as an array.
   */
  childrenKey?: string;
  /** Must contain exactly one entry with key "popupMenu" for the actions column */
  columns: ListViewColumn[];

  // ── 6. Row actions ──────────────────────────────────────────────────────────
  rowActions: ListViewRowAction[];

  // ── 7. Confirm modal ────────────────────────────────────────────────────────
  /** Required when any rowAction has color: "critical" */
  confirmModal?: {
    /** actionName of the critical rowAction that triggers the modal */
    triggerAction: string;
    /** Name of the @action that performs the confirmed delete */
    deleteAction: string;
    title: string;
    message: string;
  };

  // ── 8. Empty states ─────────────────────────────────────────────────────────
  emptyStates: {
    noData: {
      title: string;
      description?: string;
    };
    filtered: {
      /** Template string — use "{{filter}}" as the placeholder for the filter value */
      titleTemplate: string;
    };
  };

  // ── 9. Route model hints ────────────────────────────────────────────────────
  routeModel: {
    /** Always includes "page" and "pageFilter" by project convention */
    queryParams: Array<{
      name: string;
      refreshModel: boolean;
    }>;
    /**
     * Live `this.api` call expression, or null for a TODO stub.
     * e.g. "this.api.identity.entityListById(ListEnum.TRUE)"
     */
    apiCall: string | null;
    fetchCapabilities: boolean;
    /**
     * The raw API field name that paginate() filters on (snake_case).
     * Must exactly match the field on the data object — never camelCase.
     * e.g. "path", "id", "name", "creation_time"
     */
    filterKey: string;
  };
}
