/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Service, { inject as service } from '@ember/service';
import { keepLatestTask, task } from 'ember-concurrency';
import { tracked } from '@glimmer/tracking';
import { getFeatureVersion } from 'vault/utils/feature-versions';
import { compareVersions } from 'vault/utils/version-utils';

/**
 * This service returns information about a cluster's license/features, version and type (community vs enterprise).
 */

export default class VersionService extends Service {
  @service store;
  @service flags;
  @tracked features = [];
  @tracked version = null;
  @tracked clusterName = null;
  @tracked type = null;

  get isEnterprise() {
    return this.type === 'enterprise';
  }

  get isCommunity() {
    return !this.isEnterprise;
  }

  /* Features */
  get hasPKIOnly() {
    return this.features.includes('PKI-only Secrets');
  }

  get hasPerfReplication() {
    return this.features.includes('Performance Replication');
  }

  get hasDRReplication() {
    return this.features.includes('DR Replication');
  }

  get hasSentinel() {
    return this.features.includes('Sentinel');
  }

  get hasNamespaces() {
    return this.features.includes('Namespaces');
  }

  get hasControlGroups() {
    return this.features.includes('Control Groups');
  }

  // Consumption Billing will only be present on the platform-standard module introduced in 2.0.0.
  get hasConsumptionBilling() {
    return this.features.includes('Consumption Billing');
  }

  // Agent Registry requires either the platform-standard or agentic-iam license module.
  get hasAgenticIam() {
    return this.features.includes('Agentic IAM');
  }

  get hasSecretsSync() {
    const isEnterprise = this.isEnterprise;
    const isHvdManaged = this.flags.isHvdManaged;
    const onLicense = this.features.includes('Secrets Sync');

    if (!isEnterprise) return false;
    if (isHvdManaged) return true;
    if (isEnterprise && onLicense) return true;
    return false;
  }

  get versionDisplay() {
    if (!this.version) {
      return '';
    }
    return this.isEnterprise ? `v${this.version.slice(0, this.version.indexOf('+'))}` : `v${this.version}`;
  }

  @task({ drop: true })
  *getVersion() {
    if (this.version) return;
    // Fetch seal status with token to get version
    const response = yield this.store.adapterFor('cluster').sealStatus(false);
    this.version = response?.version;
    this.clusterName = response?.cluster_name;
  }

  @task
  *getType() {
    if (this.type !== null) return;
    const response = yield this.store.adapterFor('cluster').health();
    if (response.has_chroot_namespace) {
      // chroot_namespace feature is only available in enterprise
      this.type = 'enterprise';
      return;
    }
    this.type = response.enterprise ? 'enterprise' : 'community';
  }

  @keepLatestTask
  *getFeatures() {
    if (this.features?.length || this.isCommunity) {
      return;
    }
    try {
      const response = yield this.store.adapterFor('cluster').features();
      this.features = response.features;
      return;
    } catch (err) {
      // if we fail here, we're likely in DR Secondary mode and don't need to worry about it
    }
  }

  /**
   * Returns true if the current Vault binary version satisfies the minimum
   * version required for a version-gated feature.
   *
   * @param {string} key - The stable feature key defined in FEATURE_VERSIONS.
   * @returns {boolean} true when the running Vault version meets or exceeds
   *   the feature's minimum version; false if the key is unknown, the version
   *   has not yet loaded, or the version is below the minimum.
   *
   * Usage:
   *   if (this.version.hasFeature('agents')) {
   *     const result = await this.api.versioned().sys.newMethod();
   *   }
   *
   * @see vault/utils/feature-versions — add new entries to FEATURE_VERSIONS there
   * @see ui/docs/versioned-api-calls.md — full guide to the version-gating pattern
   */
  hasFeature(key) {
    const feature = getFeatureVersion(key);
    if (!feature || !this.version) {
      return false;
    }
    return compareVersions(this.version, feature.version) >= 0;
  }

  fetchVersion() {
    return this.getVersion.perform();
  }

  fetchType() {
    return this.getType.perform();
  }

  fetchFeatures() {
    return this.getFeatures.perform();
  }
}
