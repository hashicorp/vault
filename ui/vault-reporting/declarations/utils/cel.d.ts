/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */
export interface Filter {
    field: string;
    operator: '>=' | '<=' | '>' | '<' | '=' | '!=' | 'IN' | 'NOT IN';
    value: {
        type: 'timestamp' | 'list' | 'string' | 'number';
        value: unknown;
    };
}
export declare const expressionToFilters: (expression: string) => Filter[];
export declare const filtersToExpression: (filters: Filter[]) => string;
//# sourceMappingURL=cel.d.ts.map