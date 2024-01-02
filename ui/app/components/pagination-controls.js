/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';

/**
 * @module PaginationControls
 * PaginationControls components are used to paginate through item lists
 *
 * @example
 * ```js
 * <PaginationControls @startPage={{1}} @total={{100}} @size={{15}} @onChange={{this.onPageChange}} />
 * ```
 * @param {number} total - total number of items
 * @param {number} [startPage=1] - initial page number to select
 * @param {number} [size=15] - number of items to display per page
 * @param {function} onChange - callback fired on page change
 */

export default class PaginationControls extends Component {
  @tracked page;

  constructor() {
    super(...arguments);
    this.page = this.args.startPage || 1;
    this.size = this.args.size || 15; // size selector may be added in future version
  }

  get totalPages() {
    return Math.ceil(this.args.total / this.size);
  }
  get displayInfo() {
    const { total } = this.args;
    const end = this.page * this.size;
    return `${end - this.size + 1}-${end > total ? total : end} of ${total}`;
  }
  get pages() {
    // show 5 pages with 2 on either side of the current page
    let start = this.page - 2 >= 1 ? this.page - 2 : 1;
    const incrementer = start + 4;
    const end = incrementer <= this.totalPages ? incrementer : this.totalPages;
    const pageNumbers = [];
    while (start <= end) {
      pageNumbers.push(start);
      start++;
    }
    return pageNumbers;
  }
  get hasMorePages() {
    return this.pages.lastObject !== this.totalPages;
  }

  @action
  changePage(page) {
    this.page = page;
    this.args.onChange(page);
  }
}
