/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { action } from '@ember/object';
import { next } from '@ember/runloop';
import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { paginate } from 'core/utils/paginate-list';

type SortDirection = 'asc' | 'desc';

/**
 * @module ListTable
 * `ListTable` renders paginated table rows with optional row selection.
 *
 * @example
 * <ListTable
 *   @columns={{this.tableColumns}}
 *   @data={{this.data}}
 *   @selectionKeyField="path"
 *   @onSelectionChange={{this.updateSelectedItems}}
 * >
 *
 * @param {TableColumn[]} columns - Used to populate table headers and specify column display options or functionality (e.g. `isSortable`). See the `columns` in the component API for available HDS parameters @see https://helios.hashicorp.design/components/table/advanced-table?tab=code#advancedtable
 * @param {object[]} data - An array of data to display corresponding to columns. (ie. the key from a column corresponds to the parameter value of the data object your passing in)
 * @param {string} [selectionKeyField] - string of desired param to be use as a unique identifier for a selected row, if provided 'isSelectable' is set to "true" and table rows are selectable
 * @param {OnSelectionChange} [onSelectionChange] - Provided function for handling when rows are selected
 *
 * If there's an 'Action' column (ie. possibly for manipulating data rows, or navigating to a page per that row data, etc)
 * The parent component must specify the key as 'popupMenu' for that column and pass in a yield block 'popupMenu' for it to render per each item under the 'action' column.
 *
 */

interface TableColumn {
  key: string;
  label: string;
  isSortable?: boolean;
  customTableItem?: boolean; // when true, the parent yields a custom display for that column
  sortingFunction?: (a: Record<string, unknown>, b: Record<string, unknown>) => number | boolean;
}

interface SelectableRowState {
  selectionKey: string; // value of selected item
  isSelected: boolean;
}

interface OnSelectionArgs {
  selectionKey: string;
  selectionCheckboxElement: HTMLInputElement;
  selectedRowsKeys: string[];
  selectableRowsStates: SelectableRowState[];
}

type OnSelectionChange = (callbackArgs: OnSelectionArgs) => void;

interface Args {
  data: Array<Record<string, unknown>>;
  columns: TableColumn[];
  selectionKeyField?: string;
  childrenKey?: string;
  hasResizableColumns?: boolean; // when false, column resizing is disabled; defaults to true unless expandable rows are present
  page?: number; // optional page number to set current page, needed to keep pagination sync with url query param
  pageSize?: number; // optional page size, needed to keep pagination sync with url query param & keep page size
  hidePagination?: boolean; // when true, pagination controls are not rendered
  /**
   * When provided, overrides @data.length as the total item count passed to
   * Hds::Pagination::Numbered. Use this when the parent (e.g. Page::ListView)
   * has already sliced @data to a single page — the pagination control still
   * needs to know the true unsliced total to render page buttons correctly.
   * When absent, @data.length is used as before (self-contained pagination).
   */
  totalItems?: number;
  onSelectionChange?: OnSelectionChange;
  onPageChange?: CallableFunction;
  onPageSizeChange?: CallableFunction;
  /**
   * When provided by a parent that owns sort state (e.g. Page::ListView), these
   * three args keep the Hds::AdvancedTable header indicator in sync and tell
   * ListTable to skip its own internal sort (the parent has already sorted
   * @data before slicing for the current page).
   */
  sortBy?: string;
  sortOrder?: SortDirection;
  onSort?: (column: string, direction: SortDirection) => void;
}

export default class ListTable extends Component<Args> {
  @tracked currentPage;
  @tracked pageSize;
  @tracked sortColumn?: string;
  @tracked sortDirection: SortDirection = 'asc';
  //  WORKAROUND to manually re-render Hds::Pagination::Numbered to force update @currentPage
  @tracked renderPagination = true;

  constructor(owner: unknown, args: Args) {
    super(owner, args);

    this.currentPage = args.page || 1;
    this.pageSize = args.pageSize || 10;
  }

  /**
   * Decorates column definitions with a custom `sortingFunction` that uses natural / numeric
   * collation so that <Hds::AdvancedTable> does not re-sort the model with standard lexicographical order.
   */
  get tableColumns(): TableColumn[] {
    return this.args.columns.map((column) => {
      if (!column.isSortable) return column;

      return {
        ...column,
        sortingFunction: (a: Record<string, unknown>, b: Record<string, unknown>) => {
          const sortOrder = this.activeSortDirection;
          const valA = a[column.key];
          const valB = b[column.key];

          if (valA == null && valB == null) return 0;
          if (valA == null) return sortOrder === 'asc' ? 1 : -1;
          if (valB == null) return sortOrder === 'asc' ? -1 : 1;

          let result: number;
          if (typeof valA === 'string' && typeof valB === 'string') {
            result = valA.localeCompare(valB, undefined, { numeric: true, sensitivity: 'base' });
          } else {
            result = valA < valB ? -1 : valA > valB ? 1 : 0;
          }

          return sortOrder === 'asc' ? result : -result;
        },
      };
    });
  }

  get hasResizableColumns() {
    // explicit arg takes precedence
    if (this.args.hasResizableColumns !== undefined) {
      return this.args.hasResizableColumns;
    }
    // if there are nested rows, resizable columns must be disabled
    // check if there are children defined in the data structure or the childrenKey arg is present
    if (this.args.childrenKey) {
      return false;
    }
    return !this.args.data.some((item) => item['children']);
  }

  /**
   * The active sort column — driven by the parent when @sortBy is supplied
   * (e.g. Page::ListView owns sort state), otherwise tracked internally.
   */
  get activeSortColumn() {
    return this.args.sortBy ?? this.sortColumn;
  }

  /** The active sort direction — driven by parent @sortOrder when provided. */
  get activeSortDirection(): SortDirection {
    return this.args.sortOrder ?? this.sortDirection;
  }

  get sortedTableData() {
    if (this.activeSortColumn) {
      const column = this.activeSortColumn;
      const direction = this.activeSortDirection;

      return [...this.args.data].sort((a, b) => {
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
    return this.args.data;
  }

  /**
   * When the parent supplies @totalItems the data is already pre-sliced to one
   * page AND pre-sorted (Page::ListView sorts before paginating). Pass it
   * straight through — do not re-sort or re-paginate here.
   * For standalone ListTable usage (no @totalItems), apply the internal sort
   * and paginate as before.
   */
  get paginatedTableData() {
    if (this.args.totalItems !== undefined) {
      return this.args.data;
    }
    return paginate(this.sortedTableData, {
      page: this.currentPage,
      pageSize: this.pageSize,
    });
  }

  /** True total for Hds::Pagination::Numbered — use external value when supplied. */
  get paginationTotalItems() {
    return this.args.totalItems ?? this.args.data.length;
  }

  @action
  updateSort(column: string, direction: SortDirection) {
    if (this.args.onSort) {
      // Parent owns sort state — delegate so it can sort the full dataset.
      this.args.onSort(column, direction);
    } else {
      this.sortColumn = column;
      this.sortDirection = direction;
    }
  }

  @action
  async handlePaginationChange(action: 'currentPage' | 'pageSize', value: number) {
    if (action === 'pageSize') {
      await this.resetPagination();
      // external callback to handle page size changes and bubble up to parent component
      this.args.onPageSizeChange?.(value);
    }
    this[action] = value;

    // external callback to handle current page changes and bubble up to parent component
    if (action === 'currentPage') {
      this.args.onPageChange?.(value);
    }
  }

  @action
  async resetPagination() {
    // When the parent drives pagination externally (@totalItems is set), it
    // owns the page state via the URL query param — do not reset currentPage
    // here or it fights the URL. Still toggle renderPagination so
    // Hds::Pagination::Numbered picks up the new @currentPageSize.
    if (this.args.totalItems === undefined) {
      this.currentPage = 1;
    }
    this.renderPagination = false;
    //  WORKAROUND to manually re-render Hds::Pagination::Numbered to force update @currentPage
    next(() => {
      this.renderPagination = true;
    });
  }

  // TEMPLATE HELPERS
  isObject = (value: unknown) => typeof value === 'object' && value !== null;

  identifier = (cellData: Record<string, unknown>) => {
    const firstColumn = this.args.columns[0]?.key;
    // Use selectionKeyField if provided, otherwise default to value of the first column
    const identifier = this.args.selectionKeyField || firstColumn;
    return identifier ? cellData[identifier] : null;
  };

  // Returns true if the row's data has the children key, regardless of whether the array is empty.
  // Child rows never have this key; only top-level (parent) rows do.
  isParentRow = (cellData: Record<string, unknown>) => {
    const key = this.args.childrenKey ?? 'children';
    return key in cellData;
  };
}
