/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';

import type CurrentClusterService from 'vault/services/current-cluster';
import type VersionService from 'vault/services/version';
import type NamespaceService from 'vault/services/namespace';
import type ClusterModel from 'vault/models/cluster';

interface Args {
  isEngine?: boolean;
}

export default class SidebarNavToolsComponent extends Component<Args> {
  @service declare readonly currentCluster: CurrentClusterService;
  @service declare readonly version: VersionService;
  @service declare readonly namespace: NamespaceService;

  get cluster() {
    return this.currentCluster.cluster as ClusterModel | null;
  }

  get hasChrootNamespace() {
    return this.cluster?.hasChrootNamespace;
  }

  get isRootNamespace() {
    // should only return true if we're in the true root namespace
    return this.namespace.inRootNamespace && !this.hasChrootNamespace;
  }
}
