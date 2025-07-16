/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */
import Component from '@glimmer/component';
import './dashboard.scss';
import type { UsageDashboardData, SimpleDatum, getUsageDataFunction } from '../../../types';
import type { IconName } from '@hashicorp/flight-icons/svg';
import type ReportingAnalyticsService from '../../../services/reporting-analytics';
interface CounterBlock {
    title: string;
    tooltipMessage: string;
    data: number;
    icon?: IconName;
    suffix?: string;
    link?: string;
    emptyText?: string;
    emptyLink?: string;
}
export interface SSUViewDashboardSignature {
    Args: {
        onFetchUsageData: getUsageDataFunction;
        isVaultDedicated: boolean;
    };
    Blocks: {
        default: [];
    };
    Element: HTMLElement;
}
export default class SSUViewDashboard extends Component<SSUViewDashboardSignature> {
    readonly reportingAnalytics: ReportingAnalyticsService;
    data?: UsageDashboardData;
    lastUpdatedTime: string;
    error?: unknown;
    constructor(owner: unknown, args: SSUViewDashboardSignature['Args']);
    fetchAllData: () => void;
    handleTrackAnalyticsEvent: (eventName: string, properties?: object, options?: object) => void;
    handleTrackSurveyLink: () => void;
    handleRefresh: () => void;
    getBarChartData: (map: Record<string, number>, exclude?: string[]) => SimpleDatum[];
    get isVaultDedicated(): boolean;
    get kvSecretsTooltipMessage(): string;
    get counters(): CounterBlock[];
    get namespace(): string;
}
export {};
//# sourceMappingURL=dashboard.d.ts.map