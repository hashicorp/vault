/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */
import Component from '@glimmer/component';
import type { UsageDashboardData } from '../../../types';
import type ReportingAnalyticsService from '../../../services/reporting-analytics';
export interface DashboardExportSignature {
    Args: {
        data?: UsageDashboardData;
    };
    Blocks: {
        default: [];
    };
    Element: HTMLElement;
}
export default class DashboardExport extends Component<DashboardExportSignature> {
    #private;
    readonly reportingAnalytics: ReportingAnalyticsService;
    handleTrackExportToggle: () => void;
    handleTrackExportOption: (option: string) => void;
    get dataAsDownloadableJSONString(): string;
    get dataAsDownloadableCSVString(): string;
}
//# sourceMappingURL=export.d.ts.map