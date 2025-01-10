/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';

import type flagsService from 'vault/services/flags';
import NamespaceService from 'vault/services/namespace';

export type Args = {
  isRootNamespace: boolean;
  replication: unknown;
  secretsEngines: unknown;
  vaultConfiguration: unknown;
  version: { isEnterprise: boolean };
};

export default class OverviewComponent extends Component<Args> {
  @service declare readonly flags: flagsService;
  @service declare readonly namespace: NamespaceService;

  /**
   * the client count card should show in the following conditions
   * Self Managed clusters that are running enterprise and showing the `root` namespace
   * Managed clusters that are running enterprise and show the `admin` namespace
   */
  // for self managed clusters, this is the `root` namespace
  // for HVD clusters, this is the `admin` namespace
  get shouldShowClientCount() {
    const { version, isRootNamespace } = this.args;
    const { flags, namespace } = this;

    // don't show client count if this isn't an enterprise cluster
    if (!version.isEnterprise) return false;

    // HVD clusters
    if (flags.isHvdManaged && namespace.currentNamespace === 'admin') return true;

    // SM clusters
    if (isRootNamespace) return true;

    return false;
  }
}
