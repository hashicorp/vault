/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */
import Component from '@glimmer/component';
import { REPLICATION_ENABLED_STATE } from '../../types/index.ts';
import './cluster-replication.scss';
export interface ClusterReplicationSignature {
    Args: {
        isDisasterRecoveryPrimary: boolean;
        disasterRecoveryState: REPLICATION_ENABLED_STATE | 'disabled';
        isPerformancePrimary: boolean;
        performanceState: REPLICATION_ENABLED_STATE | 'disabled';
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