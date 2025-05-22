/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */
import { type FilterFieldDefinition } from '../filter-bar';
import './secret-inventory.scss';
import Component from '@glimmer/component';
import { type Filter } from '@hashicorp/vault-reporting/utils/cel';
export interface SecretInventorySignature {
    Args: {
        onFilterApplied: (value: string) => void;
        filterString: string;
    };
    Blocks: {
        default: [];
    };
    Element: HTMLElement;
}
export default class SecretInventory extends Component<SecretInventorySignature> {
    quickFilters: {
        label: string;
        applyFilter: () => void;
    }[];
    filterFieldDefinitions: FilterFieldDefinition[];
    handleApplyFilters: (filters: Filter[]) => void;
    get appliedFilters(): Filter[];
}
//# sourceMappingURL=secret-inventory.d.ts.map