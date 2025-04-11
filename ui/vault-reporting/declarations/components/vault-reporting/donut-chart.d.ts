/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */
import './donut-chart.scss';
import Component from '@glimmer/component';
export interface SSUReportingDonutChartSignature {
    Args: {
        data: {
            value: number;
            label: string;
        }[];
        title: string;
    };
    Blocks: {
        default: [];
    };
    Element: HTMLElement;
}
export default class SSUReportingDonutChart extends Component<SSUReportingDonutChartSignature> {
    get data(): {
        scaleIndex: number;
        value: number;
        label: string;
    }[];
    get total(): number;
    get a11yLabel(): string;
    getOffset(width: number, height: number): string;
    getInnerRadius(width: number, height: number): number;
    getOuterRadius(width: number, height: number): number;
}
//# sourceMappingURL=donut-chart.d.ts.map