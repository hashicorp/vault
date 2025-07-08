/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */
import Component from '@glimmer/component';
export interface ConnectClusterSignature {
    Args: {
        onClusterChange: (cluster: ClusterConnection) => void;
    };
    Blocks: {
        default: [];
    };
    Element: HTMLElement;
}
export interface ClusterConnection {
    clusterName: string;
    clusterUrl: string;
    token: string;
}
export default class ConnectCluster extends Component<ConnectClusterSignature> {
    clusters: ClusterConnection[];
    activeCluster: ClusterConnection | null;
    showAddModal: boolean;
    constructor(owner: unknown, args: ConnectClusterSignature['Args']);
    handleConnect: (e: Event) => void;
    getClusterList: () => void;
    isActiveCluster: (cluster: ClusterConnection) => boolean;
    handleClusterChange: (cluster: ClusterConnection) => void;
    handleCloseModal: () => void;
    handleShowAddModal: () => void;
    get selectClusterToggleText(): string;
}
//# sourceMappingURL=connect-cluster.d.ts.map