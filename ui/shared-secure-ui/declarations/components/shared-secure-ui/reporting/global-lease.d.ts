/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */
import Component from '@glimmer/component';
import './global-lease.scss';
export interface GlobalLeaseSignature {
    Args: {
        count: number;
        quota: number;
    };
    Blocks: {
        default: [];
    };
    Element: HTMLElement;
}
export default class GlobalLease extends Component<GlobalLeaseSignature> {
    animatedPercentage: number;
    displayPercentage: number;
    initialState: boolean;
    constructor(owner: unknown, args: GlobalLeaseSignature['Args']);
    get actualPercentage(): number;
    get progressFillClass(): "ssu-global-lease__progress-fill--initial" | "ssu-global-lease__progress-fill--low" | "ssu-global-lease__progress-fill--medium" | "ssu-global-lease__progress-fill--high";
    get formattedCount(): string;
    startAnimation(): void;
    animatePercentageText(): void;
}
//# sourceMappingURL=global-lease.d.ts.map