/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */
export interface Filter {
    field: string;
    operator: '>' | '<' | '=';
    value: {
        type: 'string' | 'date' | 'duration' | 'list';
        value: string | number | string[];
    };
}
type FilterMap = {
    [key: string]: string | number | Date;
};
export declare const flattenFilters: <T = FilterMap>(filters?: Filter[]) => T;
export {};
//# sourceMappingURL=filters.d.ts.map