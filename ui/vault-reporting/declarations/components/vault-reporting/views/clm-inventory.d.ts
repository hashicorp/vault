/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */
import './clm-inventory.scss';
import Component from '@glimmer/component';
import { type Filter } from '@hashicorp/vault-reporting/utils/filters';
import { type FilterFieldDefinition } from '../filter-bar';
export interface CLMInventorySignature {
    Args: {
        onFilterApplied: (value: Filter[]) => void;
        appliedFilters: Filter[];
    };
    Blocks: {
        default: [];
    };
    Element: HTMLElement;
}
export default class CLMInventory extends Component<CLMInventorySignature> {
    quickFilters: {
        label: string;
        applyFilter: () => void;
    }[];
    filterFieldDefinitions: FilterFieldDefinition[];
}
//# sourceMappingURL=clm-inventory.d.ts.map