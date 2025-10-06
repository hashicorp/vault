/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { cached, tracked } from '@glimmer/tracking';
import { paginate } from 'core/utils/paginate-list';
import { next } from '@ember/runloop';
import { service } from '@ember/service';

import type VersionService from 'vault/services/version';
import type { ClientFilterTypes } from 'vault/vault/client-counts/activity-api';

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
  @service declare readonly version: VersionService;

  @tracked currentPage = 1;
  @tracked pageSize = 5; // Can be overridden by @setPageSize
  @tracked sortColumn = '';
  @tracked sortDirection: SortDirection = 'asc'; // default is 'asc' for consistency with HDS defaults

  //  WORKAROUND to manually re-render Hds::Pagination::Numbered to force update @currentPage
  @tracked renderPagination = true;

  constructor(owner: unknown, args: Args) {
    super(owner, args);

    const { column, direction } = this.args.initiallySortBy || {};
    if (column) {
      this.sortColumn = column;
    }
    if (direction) {
      this.sortDirection = direction;
    }

    // Override default page size with a custom amount.
    // pageSize can be updated by the user if @showPaginationSizeSelector is true
    if (this.args.setPageSize) {
      this.pageSize = this.args.setPageSize;
    }
  }

  @cached
  get columnKeys() {
    return this.args.columns.map((k: TableColumn) => k['key']);
  }

  @cached
  get paginatedTableData(): Record<string, any>[] {
    const sorted = this.sortTableData(this.args.data);
    const paginated = paginate(sorted, {
      page: this.currentPage,
      pageSize: this.pageSize,
    });

    return paginated;
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
  async resetPagination() {
    // setPageSize is intentionally NOT reset here so user changes to page size
    // are preserved regardless of whether or not the table data updates.
    this.renderPagination = false;
    this.currentPage = 1;
    //  WORKAROUND to manually re-render Hds::Pagination::Numbered to force update @currentPage
    next(() => {
      this.renderPagination = true;
    });
  }

  @action
  updateSort(column: string, direction: SortDirection) {
    // Update tracked variables so paginatedTableData recomputes
    this.sortColumn = column;
    this.sortDirection = direction;
  }

  // TEMPLATE HELPERS
  isObject = (value: any) => typeof value === 'object';

  generateQueryParams = (datum: Record<ClientFilterTypes, any>) => {
    const { namespace_path = '', mount_path = '', mount_type = '' } = datum;
    return { namespace_path, mount_path, mount_type };
  };
}
