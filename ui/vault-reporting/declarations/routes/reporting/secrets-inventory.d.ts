/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */
import Route from '@ember/routing/route';
import type ReportingInventoryController from '../../controllers/reporting/secrets-inventory';
export default class ReportingInventoryRoute extends Route {
    queryParams: {
        filters: {
            replace: boolean;
            refreshModel: boolean;
        };
        pagination: {
            replace: boolean;
            refreshModel: boolean;
        };
        sortingOrderBy: {
            replace: boolean;
            refreshModel: boolean;
        };
    };
    resetController(controller: ReportingInventoryController, isExiting: boolean): void;
}
//# sourceMappingURL=secrets-inventory.d.ts.map