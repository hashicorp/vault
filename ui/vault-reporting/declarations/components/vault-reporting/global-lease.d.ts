/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */
import Component from '@glimmer/component';
import './global-lease.scss';
import type { HdsApplicationStateSignature } from '@hashicorp/design-system-components/components/hds/application-state/index';
export interface GlobalLeaseSignature {
    Args: {
        count?: number;
        quota?: number;
    };
    Blocks: {
        default: [];
        /** We optionally yield application state to allow for overrides on empty state eg:
         * <SSUReportingGlobalLease ...>
         *   <:empty as |A|>
         *     <A.Header @title="Custom Title" />
         *     <A.Body @text="Custom description" />
         *   </:empty>
         * </SSUReportingGlobalLease>
         * */
        empty: HdsApplicationStateSignature['Blocks']['default'];
    };
    Element: HTMLElement;
}
export default class GlobalLease extends Component<GlobalLeaseSignature> {
    get percentage(): number;
    get progressFillClass(): "" | "ssu-global-lease__progress-fill--exceeded";
    get formattedCount(): string;
    get percentageString(): string;
    get hasData(): boolean | 0 | undefined;
    get description(): import("@ember/template").SafeString | undefined;
    get linkUrl(): "https://developer.hashicorp.com/vault/tutorials/operations/resource-quotas#global-default-lease-count-quota" | undefined;
    get alert(): {
        color: 'warning' | 'neutral';
        description: string;
    } | undefined;
}
//# sourceMappingURL=global-lease.d.ts.map