/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */
import Controller from '@ember/controller';
import type { Filter } from '../../utils/filters';
export default class ReportingInventoryController extends Controller {
    queryParams: string[];
    filters: Filter[];
    pagination: {
        page_size: number;
        next_page_token: string | null;
        previous_page_token: string | null;
    };
    sortingOrderBy: string[];
    handleApplyFilter: (filters: Filter[]) => void;
    resetPagination: () => void;
    handlePageChange: (pagination: {
        next_page_token: string | null;
        previous_page_token: string | null;
        page_size: number;
    }) => void;
    handleApplySort: (orderBy: string[]) => void;
}
//# sourceMappingURL=secrets-inventory.d.ts.map