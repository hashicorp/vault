/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */
import Component from '@glimmer/component';
import './cluster-replication.scss';
export declare const ENABLED_STATE = "enabled";
export declare const DISABLED_STATE = "disabled";
export interface ClusterReplicationSignature {
    Args: {
        isDisasterRecoveryPrimary: boolean;
        disasterRecoveryState: string;
        isPerformancePrimary: boolean;
        performanceState: string;
    };
    Blocks: {
        default: [];
    };
    Element: HTMLElement;
}
export default class ClusterReplication extends Component<ClusterReplicationSignature> {
    get disasterRecoveryBadge(): {
        icon: 'check' | 'x';
        text: string;
        color: 'success' | 'neutral';
    };
    get performanceBadge(): {
        icon: 'check' | 'x';
        text: string;
        color: 'success' | 'neutral';
    };
    get disasterRecoveryRole(): "Primary" | "Secondary";
    get performanceRole(): "Primary" | "Secondary";
}
//# sourceMappingURL=cluster-replication.d.ts.map