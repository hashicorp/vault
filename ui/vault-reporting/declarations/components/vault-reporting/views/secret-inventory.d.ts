/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */
import './secret-inventory.scss';
import Component from '@glimmer/component';
import { type Filter } from '@hashicorp/vault-reporting/utils/filters';
import { type InventoryTableColumn, type QuickFilter } from '../inventory-table';
import { type FilterFieldDefinition } from '../filter-bar';
import type ReportingApiService from '@hashicorp/vault-reporting/services/reporting-api';
export interface SecretInventorySignature {
    Args: {
        onFilterApplied: (value: Filter[]) => void;
        onSortApplied: (sort: string[]) => void;
        appliedFilters: Filter[];
        appliedSort: string[];
        onPageChange: (pagination: object) => void;
        pageSize: number;
        nextPageToken: string | null;
        previousPageToken: string | null;
    };
    Blocks: {
        default: [];
    };
    Element: HTMLElement;
}
export default class SecretInventory extends Component<SecretInventorySignature> {
    readonly reportingApi: ReportingApiService;
    page: Record<string, unknown>[];
    isLoading: boolean;
    nextPageToken: string | null;
    previousPageToken: string | null;
    pageSize: number;
    totalUnfilteredCount: number;
    constructor(owner: unknown, args: SecretInventorySignature['Args']);
    handlePageSizeChange: (newPageSize: number) => void;
    lastUpdatedTime: string;
    quickFilters: QuickFilter[];
    filterFieldDefinitions: FilterFieldDefinition[];
    columns: InventoryTableColumn[];
    fetchQuickFilterCount: (filter: QuickFilter) => Promise<QuickFilter>;
    fetchQuickFilterCounts: () => Promise<void>;
    handleDataUpdate: (filters: Filter[]) => Promise<void>;
}
//# sourceMappingURL=secret-inventory.d.ts.map