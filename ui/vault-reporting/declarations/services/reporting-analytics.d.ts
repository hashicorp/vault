/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */
import Service from '@ember/service';
export default class ReportingAnalytics extends Service {
    get analytics(): {
        trackEvent: (event: string, properties?: object, options?: object) => void;
    } | undefined;
    trackEvent(event: string, properties?: object, options?: object): void;
}
declare module '@ember/service' {
    interface Registry {
        reportingAnalytics: ReportingAnalytics;
    }
}
//# sourceMappingURL=reporting-analytics.d.ts.map