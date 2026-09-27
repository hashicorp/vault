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
   * filtering (via config.filter.filterKey) and pagination internally.
   */
  model: unknown[];
  /** Current page number. Changes on pagination. */
  page?: number;
  /** Current page size. Comes from the controller query param so it survives route transitions. */
  pageSize?: number;
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

  get pageSize() {
    return this.args.pageSize ?? 10;
  }
  @tracked activeModalConfig: ListViewModalConfig | null = null;
  @tracked activeModalItem: Record<string, unknown> | null = null;

  @action
  onFilterChange(value: string | null) {
    this.pageFilter = value ?? '';
    // Only reset the page QP to 1 when the user is on a page other than 1.
    // Skipping the transition when already on page 1 prevents a spurious
    // route re-render (and model re-fetch) every time the filter changes.
    if ((this.args.page ?? 1) === 1) return;
    // Wrapped in try/catch because transitionTo({ queryParams }) requires an
    // active route ("has last route info") — in integration tests and any
    // render context with no running route it would otherwise throw.
    try {
      this.router.transitionTo({ queryParams: { page: 1 } });
    } catch (_) {
      // no active route — filter state still updates, QP reset is a no-op
    }
  }

  /**
   * Called by ListTable when the user clicks a pagination control.
   * Transitions the `page` query param so the route model refreshes and the
   * URL stays in sync with the displayed page.
   */
  @action
  onPageChange(page: number) {
    try {
      this.router.transitionTo({ queryParams: { page } });
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
   * Applies client-side filtering (via config.filter.filterKey) and pagination
   * to the raw @model array. Called once per render cycle when model, pageFilter,
   * or page changes — cheap since paginate() is a pure array operation.
   */
  get filteredData() {
    return paginate(this.args.model as Parameters<typeof paginate>[0], {
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
    return !!(this.pageFilter && this.filteredData?.length === 0);
  }

  get isNoData() {
    return !this.pageFilter && !this.hasData;
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
