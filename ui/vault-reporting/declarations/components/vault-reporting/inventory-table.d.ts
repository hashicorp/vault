/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */
import { type HdsTableColumn } from '@hashicorp/design-system-components/components';
import { type FilterFieldDefinition } from './filter-bar';
import './inventory-table.scss';
import Component from '@glimmer/component';
import { type Filter } from '@hashicorp/vault-reporting/utils/filters';
import type RouterService from '@ember/routing/router-service';
export type InventoryTableColumn = HdsTableColumn & {
    cellFormatter?: (data: Record<string, unknown>) => string;
};
export type QuickFilter = {
    label: string;
    tooltip?: string;
    getParams: () => {
        filters?: Filter[];
    };
    count?: number;
};
export interface InventoryTableSignature {
    Args: {
        onFilterApplied: (value: Filter[]) => void;
        onSortApplied: (sort: string[]) => void;
        onPageChange: (pagination: object) => void;
        onPageSizeChange?: (pageSize: number) => void;
        pageSize?: number;
        nextPageToken?: string | null;
        previousPageToken?: string | null;
        appliedFilters: Filter[];
        appliedSort: string[];
        filterFieldDefinitions: FilterFieldDefinition[];
        columns: InventoryTableColumn[];
        rows: Record<string, unknown>[];
        quickFilters: QuickFilter[];
        isLoading?: boolean;
        totalUnfilteredCount: number;
    };
    Blocks: {
        default: [];
    };
    Element: HTMLElement;
}
export default class InventoryTable extends Component<InventoryTableSignature> {
    readonly router: RouterService;
    handleApplyFilters: (filters: Filter[]) => void;
    handleSort: (sortBy: string, sortOrder: "asc" | "desc") => void;
    handlePageChange: (direction: string) => void;
    handlePageSizeChange: (pageSize: number) => void;
    get isDisabledNext(): boolean | undefined;
    get isDisabledPrev(): boolean | undefined;
    get routeName(): string;
    getData: (key?: string, row?: Record<string, unknown>) => string;
    getParams: (quickFilter: QuickFilter) => {
        filters?: Filter[];
    };
    get sortBy(): {
        key: string;
        order: "asc" | "desc";
    };
    get rows(): Record<string, unknown>[];
}
//# sourceMappingURL=inventory-table.d.ts.map