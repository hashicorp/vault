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
 * `ListTable` component is used for rendering a list of items in a table.
 *
 * @example
 * <ListTable
 *   @columns={{this.tableColumns}}
 *   @data={{this.data}}
 *   @isSelectable={{true}}
 *   @onSelectionChange={{this.updateSelectedItems}}
 * >
 *
 * @param {array} columns - An array of type TableColumn, for populating table headers and any optional functionality (ie. isSortable...)
 * @param {array} data - An array of data to display corresponding to columns. (ie. the key from a column corresponds to the parameter value of the data object your passing in)
 * @param {string} [selectionKey] - string of desired param value to be used as a unique identifier for a selected row - setting this arg automatically sets 'isSelectable' on the table to make rows selectable
 * @param {function} [onSelectionChange] - Provided function for handling when rows are selected
 *
 * For custom column items that are not 1 to 1 with their dataset (ie. these items have conditional icons, colors, generated text etc)
 * within the column type, the 'customTableItem' flag will allow the parent component to {{yield}} any custom implementation from the parent for those items
 * but the yield block in the parent must be 'customTableItem'.
 *
 * similarly, if 'isSelectable' is true and 'onSelectionChange' is being handled
 * For displaying new content based on what's selected, the parent component can also pass in a yield block 'selectedItems'
 *
 * If there's an 'Action' column (ie. possibly for manipulating data rows, or navigating to a page per that row data, etc)
 * The parent component must specify the key as 'popupMenu' for that column and pass in a yield block 'popupMenu' for it to render per each item under the 'action' column.
 *
 */

interface TableColumn {
  key: string;
  label: string;
  selectionKey?: string;
  customTableItem?: boolean;
  onSelectionChange?: CallableFunction;
}

interface Args {
  data: Array<object>;
  columns: TableColumn[];
}

export default class ListTable extends Component<Args> {
  @tracked currentPage = 1;
  @tracked pageSize = 10;

  //  WORKAROUND to manually re-render Hds::Pagination::Numbered to force update @currentPage
  @tracked renderPagination = true;

  get paginatedTableData() {
    const paginated = paginate(this.args.data, {
      page: this.currentPage,
      pageSize: this.pageSize,
    });
    return paginated;
  }

  get columnKeys() {
    return this.args.columns.map((k: TableColumn) => k['key'] ?? k['label']);
  }

  @action
  handlePaginationChange(action: 'currentPage' | 'pageSize', value: number) {
    this[action] = value;
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
}
