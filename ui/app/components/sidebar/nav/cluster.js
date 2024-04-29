/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';

export default class SidebarNavClusterComponent extends Component {
  @service currentCluster;
  @service flags;
  @service version;
  @service auth;
  @service namespace;

  get cluster() {
    return this.currentCluster.cluster;
  }

  get isRootNamespace() {
    // should only return true if we're in the true root namespace
    return this.namespace.inRootNamespace && !this.cluster?.hasChrootNamespace;
  }

  get badgeText() {
    const isManaged = this.flags.isManaged;
    const onLicense = this.version.hasSecretsSync;
    const isEnterprise = this.version.isEnterprise;

    if (isManaged) return 'Plus';
    if (isEnterprise && !onLicense) return 'Premium';
    if (!isEnterprise) return 'Enterprise';
    // no badge for Enterprise clusters with Secrets Sync on their license--the only remaining option.
    return '';
  }
}
