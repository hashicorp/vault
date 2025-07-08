/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */
import Component from '@glimmer/component';
import { type Filter } from '@hashicorp/vault-reporting/utils/filters';
import './filter-bar.scss';
export interface FilterFieldDefinition {
    name: string;
    label: string;
    type: 'text' | 'single-select' | 'multi-select' | 'lookback' | 'list' | 'search';
    options?: {
        name: string;
        value: string;
    }[];
}
export interface FilterBarSignature {
    Args: {
        onFiltersApplied: (filters: Filter[]) => void;
        appliedFilters?: Filter[];
        filterFieldDefinitions: FilterFieldDefinition[];
    };
    Blocks: {
        default: [];
    };
    Element: HTMLElement;
}
export default class FilterBar extends Component<FilterBarSignature> {
    updateFilters: (filters: Record<string, Filter>) => void;
    handleClearFilters: () => void;
    handleDismissFilter: (key: string) => void;
    handleRadioChange: (name: string, value: string) => void;
    handleCheckboxChange: (name: string, value: string, event: Event) => void;
    handleTextInputChange: (name: string, event: Event) => void;
    handleLookbackChange: (name: string, value: string) => void;
    isEqual: (a: string, b: string) => boolean;
    isCheckboxChecked: (name: string, value: string) => boolean;
    isLookbackChecked: (name: string, value: string) => boolean;
    getValue: (name: string) => string;
    getOperator: (name: string) => "" | ">" | "<" | "=";
    friendlyAppliedString: (appliedFilter: Filter) => string;
    get appliedFilters(): Record<string, Filter>;
    get appliedFilterTags(): {
        key: string;
        text: string;
    }[];
    get appliedFiltersCount(): number;
    get hasAppliedFilters(): boolean;
}
//# sourceMappingURL=filter-bar.d.ts.map