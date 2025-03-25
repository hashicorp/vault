/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */
import Component from '@glimmer/component';
import './dashboard.scss';
import type { IUsageDashboardService, UsageDashboardData, SimpleDatum } from '../../../types';
import type { IconName } from '@hashicorp/flight-icons/svg';
interface CounterBlock {
    title: string;
    data: number;
    icon?: IconName;
    suffix?: string;
    link?: string;
}
export interface SSUViewDashboardSignature {
    Args: {
        service: IUsageDashboardService;
    };
    Blocks: {
        default: [];
    };
    Element: HTMLElement;
}
export default class SSUViewDashboard extends Component<SSUViewDashboardSignature> {
    data?: UsageDashboardData;
    lastUpdatedTime: string;
    constructor(owner: unknown, args: SSUViewDashboardSignature['Args']);
    fetchAllData: () => Promise<void>;
    getBarChartData: (map: Record<string, number>) => SimpleDatum[];
    get counters(): CounterBlock[];
}
export {};
//# sourceMappingURL=dashboard.d.ts.map