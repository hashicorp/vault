/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */
import './horizontal-bar-chart.scss';
import Component from '@glimmer/component';
import { HdsLinkStandalone } from '@hashicorp/design-system-components/components';
import type { SimpleDatum } from '../../types/index.ts';
import type { HdsApplicationStateSignature } from '@hashicorp/design-system-components/components/hds/application-state/index';
export interface SSUReportingHorizontalBarChartSignature {
    Args: {
        data: SimpleDatum[];
        title: string;
        description?: string;
        /** Custom text for the link (defaults to "View all") */
        linkText?: string;
        /** Icon to display with the link (defaults to "arrow-right") */
        linkIcon?: HdsLinkStandalone['icon'];
        /** URL for the link - if not provided, no link will be shown */
        linkUrl?: string;
        /** Target for the link - defaults to "_self" */
        linkTarget?: '_blank' | '_self';
    };
    Blocks: {
        default: [];
        /** We optionally yield application state to allow for overrides on empty state eg:
         * <SSUReportingHorizontalBarChart ...>
         *   <:empty as |A|>
         *     <A.Header @title="Custom Title" />
         *     <A.Body @text="Custom description" />
         *   </:empty>
         * </SSUReportingHorizontalBarChart>
         * */
        empty: HdsApplicationStateSignature['Blocks']['default'];
    };
    Element: HTMLElement;
}
export default class SSUReportingHorizontalBarChart extends Component<SSUReportingHorizontalBarChartSignature> {
    xRangeOffsetWidth: number;
    get hasData(): boolean;
    get data(): SimpleDatum[];
    get total(): number;
    get a11yLabel(): string;
    get yDomain(): string[];
    get xDomain(): number[];
    get rangeHeight(): number;
    get yRange(): number[];
    get emptyStateTitle(): string;
    get emptyStateDescription(): string;
    get emptyStateLinkText(): string;
    get description(): string | undefined;
    get linkUrl(): string | undefined;
    getXRange: (width: number) => number[];
    handleAxisOffset: (offsetWidth: number) => void;
}
//# sourceMappingURL=horizontal-bar-chart.d.ts.map