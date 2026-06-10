/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { action } from '@ember/object';
import { next } from '@ember/runloop';
import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { paginate } from 'core/utils/paginate-list';

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
  customTableItem?: boolean; // when true, the parent yields a custom display for that column
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
  data: Array<object>;
  columns: TableColumn[];
  selectionKeyField?: string;
  page?: number; // optional page number to set current page, needed to keep pagination sync with url query param
  pageSize?: number; // optional page size, needed to keep pagination sync with url query param & keep page size
  onSelectionChange?: OnSelectionChange;
  onPageChange?: CallableFunction;
  onPageSizeChange?: CallableFunction;
}

export default class ListTable extends Component<Args> {
  @tracked currentPage;
  @tracked pageSize;
  //  WORKAROUND to manually re-render Hds::Pagination::Numbered to force update @currentPage
  @tracked renderPagination = true;

  constructor(owner: unknown, args: Args) {
    super(owner, args);

    this.currentPage = args.page || 1;
    this.pageSize = args.pageSize || 10;
  }

  get paginatedTableData() {
    const paginated = paginate(this.args.data, {
      page: this.currentPage,
      pageSize: this.pageSize,
    });
    return paginated;
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
    this.renderPagination = false;
    this.currentPage = 1;
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
}
