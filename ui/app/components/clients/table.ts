/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import { paginate } from 'core/utils/paginate-list';

/**
 * @module ClientsTable
 * ClientsTable renders a paginated table for a passed dataset. HDS table components handle basic sorting
 * but because this table is paginated, all sorting happens manually.
 * See https://helios.hashicorp.design/components/table/table?tab=code#component-api
 * for the full component API and list of supported args.
 *
 * @example
 * <Clients::Table
 *    @data={{this.dataset}}
 *    @showPaginationSizeSelector={{false}}
 *    @pageSize={{100}}
 *    @initiallySortBy={{hash column="clients" direction="desc"}}
 *    @columns={{array
 *      (hash key="namespace" label="Namespace" isSortable=true)
 *      (hash key="label" label="Mount path" isSortable=true)
 *      (hash key="mount_type" label="Mount type" isSortable=true)
 *      (hash key="clients" label="Counts" isSortable=true)
 *    }}
 *  />
 */

interface TableColumn {
  key: string;
  label: string;
  isSortable: boolean;
}

type SortDirection = 'asc' | 'desc';

interface Args {
  columns: TableColumn[];
  data: Record<string, any>[];
  initiallySortBy?: { column?: string; direction?: SortDirection };
  setPageSize?: number;
  showPaginationSizeSelector?: boolean;
}

export default class ClientsTable extends Component<Args> {
  @tracked currentPage = 1;
  @tracked pageSize = 5; // Can be overridden by @setPageSize
  @tracked sortColumn = '';
  @tracked sortDirection: SortDirection;

  constructor(owner: unknown, args: Args) {
    super(owner, args);

    const { column = '', direction = 'asc' } = this.args.initiallySortBy || {};
    this.sortColumn = column;
    this.sortDirection = direction; // default is 'asc' for consistency with HDS defaults

    // Override default page size with a custom amount.
    // pageSize can be updated by the end user if @showPaginationSizeSelector is true
    if (this.args.setPageSize) {
      this.pageSize = this.args.setPageSize;
    }
  }

  get paginatedTableData(): Record<string, any>[] {
    const sorted = this.sortTableData(this.args.data);
    const paginated = paginate(sorted, {
      page: this.currentPage,
      pageSize: this.pageSize,
    });

    return paginated;
  }

  get columnKeys() {
    return this.args.columns.map((k: TableColumn) => k['key']);
  }

  // This table is paginated so we cannot use any out of the box filtering
  // from the HDS component and must manually sort data.
  sortTableData(data: Record<string, any>[]): Record<string, any>[] {
    if (this.sortColumn) {
      return [...data].sort((a, b) => {
        const valA = a[this.sortColumn];
        const valB = b[this.sortColumn];

        if (valA < valB) return this.sortDirection === 'asc' ? -1 : 1;
        if (valA > valB) return this.sortDirection === 'asc' ? 1 : -1;
        return 0;
      });
    }
    return data;
  }

  @action
  handlePaginationChange(action: 'currentPage' | 'pageSize', value: number) {
    this[action] = value;
  }

  @action
  updateSort(column: string, direction: SortDirection) {
    // Update tracked variables so paginatedTableData recomputes
    this.sortColumn = column;
    this.sortDirection = direction;
  }
}
