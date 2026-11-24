/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { action } from '@ember/object';
import { service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import Component from '@glimmer/component';
import { paginate } from 'core/utils/paginate-list';
import routerLookup from 'core/utils/router-lookup';
import type RouterService from '@ember/routing/router-service';
import type ApiService from 'vault/services/api';
import type FlashMessageService from 'vault/services/flash-messages';
import type { ListViewColumn, ListViewDisplayConfig, ListViewModalConfig } from 'core/types/list-view-config';

type SortDirection = 'asc' | 'desc';

/**
 * @module Page::ListView
 *
 * A shared structural wrapper for all Vault list views. Owns the Page::Header,
 * breadcrumbs, toolbar, ListTable, pagination, and empty states.
 *
 * Accepts two kinds of args:
 *  - @config  — static, per-route display settings (title, breadcrumbs, columns,
 *               empty-state copy, primaryAction button, filter config). Stable
 *               across renders; Glimmer will not re-render the header or
 *               breadcrumbs when only @model / @page change.
 *  - @model, @page — dynamic runtime values that update on every model refresh
 *               or query-param change.
 *
 * Config-driven rendering (no named block needed for common cases):
 *  - config.primaryAction  — renders Hds::Button in Page::Header actions area
 *  - config.filter         — renders FilterInput in toolbar; filter state is
 *                            owned internally as @tracked pageFilter. The
 *                            component resets to '' on route exit automatically
 *                            because Glimmer destroys the instance on navigation.
 *
 * Named blocks (for resource-specific content only):
 *  - <:headerActions>      — additional buttons beyond primaryAction (rare)
 *  - <:toolbarActions>     — right-side toolbar controls
 *  - <:popupMenu as |rowData|>  — ··· dropdown per parent row
 *  - <:confirmModal as |item config onClose|>
 *                          — escape hatch to render a custom confirm modal;
 *                            yields activeModalItem, activeModalConfig, and
 *                            closeModal. When absent, the built-in Hds::Modal
 *                            driven by config.rowActions is rendered instead.
 *
 * Built-in cell value types (set via column.valueType — no yield block needed):
 *  - "text"           — plain string (default)
 *  - "link"           — Hds::Link::Inline; uses column.route + rowData[column.routeModelKey ?? "id"]
 *  - "icon-text"      — Hds::Icon (from rowData[column.iconKey] or column.icon) + tooltip + text
 *  - "copy"           — Hds::Copy::Button (icon-only) + plain text
 *  - "header-tooltip" — plain text cell; ⓘ tooltip injected into the column header via
 *                       column.headerTooltip (copied to the HDS-native tooltip at normalization time)
 *  - "status"         — Hds::Badge; color from column.statusMap, icon from column.statusIconMap
 *  - "date-time"      — {{date-format}} wrapped in <time>; format: "MMM dd, yyyy h:mm a"
 *
 * @example
 * <Page::ListView
 *   @config={{this.model.listViewConfig}}
 *   @model={{this.model.groups}}
 *   @page={{this.page}}
 * >
 *   <:popupMenu as |rowData|>...</:popupMenu>
 *   <:confirmModal as |item config onClose|>...</:confirmModal>
 * </Page::ListView>
 */

interface Args {
  /**
   * Static per-route display settings: title, description, breadcrumbs, columns,
   * empty-state copy, optional primaryAction button, and optional filter config.
   * Defined as a class property (not a getter) on the controller for a stable
   * object reference — Glimmer uses stability to avoid re-rendering the header
   * and breadcrumbs when only model/filter/page changes.
   * See ListViewDisplayConfig in vault/list-view-config.
   */
  config: ListViewDisplayConfig;
  /**
   * Raw (unpaginated) data array from the route model. Page::ListView applies
   * filtering (via config.filter.filterKey) and pagination internally unless
   * the caller has already filtered the model.
   */
  model: unknown[];
  /**
   * Search text already applied by the caller, used to select filtered-empty
   * state.
   */
  filterValue?: string;
  /** Current page number. Changes on pagination. */
  page?: number;
  /** Current page size. Comes from the controller query param so it survives route transitions. */
  pageSize?: number;
  /**
   * Active sort column key. Comes from the controller sortBy query param so the
   * sort state is encoded in the URL and survives model refreshes triggered by
   * the page query param changing. When undefined no sort is applied.
   */
  sortBy?: string;
  /**
   * Active sort direction. Comes from the controller sortOrder query param.
   * Defaults to 'asc' when undefined.
   */
  sortOrder?: SortDirection;
  /**
   * Called when the user clicks a sortable column header. The parent route
   * template is expected to bubble this up to the controller so the sortBy and
   * sortOrder query params are updated in the URL.
   */
  onSortChange?: (sortBy: string, sortOrder: SortDirection) => void;
}

export default class PageListViewComponent extends Component<Args> {
  @service declare readonly api: ApiService;
  @service declare readonly flashMessages: FlashMessageService;

  // Use routerLookup so this component works in both the main app (service:router)
  // and in Ember engines that alias it as service:app-router.
  get router(): RouterService {
    return routerLookup(this) as RouterService;
  }

  @tracked pageFilter = '';

  // Returns the effective filter value, checking the internal pageFilter state and the external filterValue argument.
  get filterValue() {
    return this.pageFilter || this.args.filterValue || '';
  }

  get pageSize() {
    return this.args.pageSize ?? 10;
  }
  @tracked activeModalConfig: ListViewModalConfig | null = null;
  @tracked activeModalItem: Record<string, unknown> | null = null;

  @action
  onFilterChange(value: string | null) {
    this.pageFilter = value ?? '';
    // No router transition needed: paginate() already resets to page 1 when
    // the active page exceeds the filtered total (page > lastPage → currentPage = 1).
    // Calling transitionTo({ queryParams: { page: 1 } }) would trigger a full
    // route model refresh (page has refreshModel: true), which destroys the
    // component instance and wipes @tracked pageFilter — breaking search from
    // any page other than page 1.
  }

  /**
   * Called by ListTable when the user clicks a pagination control.
   * Transitions the `page` query param together with the current sortBy/sortOrder
   * so the URL always carries the full view state (page + sort) as a unit.
   * Without this, a transitionTo({queryParams:{page:2}}) would drop sortBy/sortOrder
   * from the URL and reset them to their controller defaults on the next render.
   */
  @action
  onPageChange(page: number) {
    try {
      this.router.transitionTo({
        queryParams: {
          page,
          ...(this.args.sortBy !== undefined ? { sortBy: this.args.sortBy } : {}),
          ...(this.args.sortOrder !== undefined ? { sortOrder: this.args.sortOrder } : {}),
        },
      });
    } catch (_) {
      // no active route (integration test context) — no-op
    }
  }

  /**
   * Called by ListTable when the user changes the page size selector.
   * Transitions the pageSize query param on the controller so the value
   * persists across page-number transitions (which reload the route model).
   * pageSize does not have refreshModel:true so no API call is made.
   */
  @action
  onPageSizeChange(size: number) {
    try {
      this.router.transitionTo({ queryParams: { pageSize: size } });
    } catch (_) {
      // no active route (integration test context) — no-op
    }
  }

  @action
  openModal(modalConfig: ListViewModalConfig, item: Record<string, unknown>) {
    this.activeModalConfig = modalConfig;
    this.activeModalItem = item;
  }

  @action
  closeModal() {
    this.activeModalConfig = null;
    this.activeModalItem = null;
  }

  @action
  updateSort(column: string, direction: SortDirection) {
    this.args.onSortChange?.(column, direction);
  }

  @action
  async confirmDelete() {
    const item = this.activeModalItem;
    const modal = this.activeModalConfig;
    if (!item || !modal) return;
    try {
      if (modal.apiCall) {
        const { service, method, argKey } = modal.apiCall;
        const svc = this.api[service as keyof ApiService] as Record<
          string,
          (arg: unknown) => Promise<unknown>
        >;
        await svc[method]?.(item[argKey]);
      } else {
        await modal.deleteAction?.(item);
      }
      this.flashMessages.success(`Successfully deleted ${item[modal.itemDisplayKey ?? 'id']}`);
      this.router.refresh(this.router.currentRouteName ?? undefined);
      this.closeModal();
    } catch (err) {
      this.flashMessages.danger((err as Error)?.message ?? 'An error occurred. Please try again.');
    }
  }

  // ── Computed data ────────────────────────────────────────────────────────────
  /**
   * Returns the full model array sorted by the active sort column/direction.
   * sortBy/sortOrder come from URL query params (via @args) so the sort survives
   * the model refresh triggered when the page query param changes.
   * Sorting happens before pagination so the entire dataset is ordered globally,
   * not just the current page's slice.
   */
  get sortedModel(): unknown[] {
    const data = this.args.model;
    const column = this.args.sortBy;
    if (!column || !Array.isArray(data)) return data;

    const direction: SortDirection = this.args.sortOrder ?? 'asc';

    return [...(data as Record<string, unknown>[])].sort((a, b) => {
      const valA = a[column];
      const valB = b[column];

      if (valA == null && valB == null) return 0;
      if (valA == null) return 1;
      if (valB == null) return -1;

      let result: number;
      if (typeof valA === 'string' && typeof valB === 'string') {
        // Use natural/numeric collation so values with numbers (e.g. "item-2", "item-10") sort in human order rather than lexicographically.
        result = valA.localeCompare(valB, undefined, { numeric: true, sensitivity: 'base' });
      } else {
        result = valA < valB ? -1 : valA > valB ? 1 : 0;
      }

      return direction === 'asc' ? result : -result;
    });
  }

  /**
   * Applies client-side filtering (via config.filter.filterKey) and pagination
   * to the sorted model array. Sorting is applied before slicing so page 1
   * always contains the globally first items under the active sort order.
   * Called once per render cycle when model, sortColumn, sortDirection,
   * pageFilter, or page changes — cheap since paginate() is a pure array operation.
   */
  get filteredData() {
    return paginate(this.sortedModel as Parameters<typeof paginate>[0], {
      page: this.args.page ?? 1,
      pageSize: this.pageSize,
      filter: this.pageFilter || undefined,
      filterKey: this.args.config.filter?.filterKey,
    });
  }

  get hasData() {
    return (this.filteredData?.meta?.total ?? this.filteredData?.length ?? 0) > 0;
  }

  get isFilteredEmpty() {
    return !!(this.filterValue && this.filteredData?.length === 0);
  }

  get isNoData() {
    return !this.filterValue && !this.hasData;
  }

  // ── Column normalization ────────────────────────────────────────────────────
  /**
   * Normalizes the column definitions so that `header-tooltip` columns have
   * their `headerTooltip` value copied to the HDS-native `tooltip` property —
   * Hds::AdvancedTable renders it as a ⓘ on the <th> automatically when
   * `tooltip` is present on the column object.
   * Called once per config change (not per row) so it is cheap.
   */
  get normalizedColumns(): ListViewColumn[] {
    return this.args.config.columns.map((col) => ({
      ...col,
      ...(col.headerTooltip ? { tooltip: col.headerTooltip } : {}),
    }));
  }
}
