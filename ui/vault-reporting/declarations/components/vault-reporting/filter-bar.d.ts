/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */
import Component from '@glimmer/component';
import { type Filter } from '../../utils/cel.ts';
export interface FilterFieldDefinition {
    name: string;
    label: string;
    type: 'text' | 'multiselect' | 'number' | 'daterange';
    options?: {
        name: string;
        value: string;
    }[];
}
export interface FilterBarSignature {
    Args: {
        onFiltersApplied: (filters: Filter[]) => void;
        appliedFilters: Filter[];
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
    handleMultiselectChange: (name: string, event: Event) => void;
    handleTextInputChange: (name: string, event: Event) => void;
    handleNumberChange: (event: Event) => void;
    handleDateRangeChange: (name: string, event: Event) => void;
    isEqual: (a: string, b: string) => boolean;
    isChecked: (name: string, value: string) => boolean;
    getValue: (name: string) => string;
    getOperator: (name: string) => "" | ">=" | "<=" | ">" | "<" | "=" | "!=" | "IN" | "NOT IN";
    friendlyAppliedString: (appliedFilter: Filter) => string;
    get appliedFilters(): Record<string, Filter>;
    get appliedFiltersCount(): number;
    get hasAppliedFilters(): boolean;
}
//# sourceMappingURL=filter-bar.d.ts.map