/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */
import './horizontal-bar-chart.scss';
import Component from '@glimmer/component';
import type { SimpleDatum } from '../../../types/reporting/index.ts';
export interface SSUReportingHorizontalBarChartSignature {
    Args: {
        data: SimpleDatum[];
        title: string;
        description?: string;
        linkUrl?: string;
    };
    Blocks: {
        default: [];
    };
    Element: HTMLElement;
}
export default class SSUReportingHorizontalBarChart extends Component<SSUReportingHorizontalBarChartSignature> {
    xRangeOffsetWidth: number;
    get data(): SimpleDatum[];
    get total(): number;
    get a11yLabel(): string;
    get yDomain(): string[];
    get xDomain(): number[];
    get rangeHeight(): number;
    get yRange(): number[];
    getXRange: (width: number) => number[];
    handleAxisOffset: (offsetWidth: number) => void;
}
//# sourceMappingURL=horizontal-bar-chart.d.ts.map