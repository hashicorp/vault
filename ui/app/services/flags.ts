/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Service, { inject as service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { keepLatestTask } from 'ember-concurrency';
import { DEBUG } from '@glimmer/env';
import type StoreService from 'vault/services/store';
import type VersionService from 'vault/services/version';

const FLAGS = {
  vaultCloudNamespace: 'VAULT_CLOUD_ADMIN_NAMESPACE',
};

/**
 * This service is for managing feature flags relevant to the Vault Cluster.
 * For now, the only available feature flag is VAULT_CLOUD_ADMIN_NAMESPACE
 * indicates that the Vault cluster is managed rather than self-managed.
 * Flags are fetched in the application route once from sys/internal/ui/feature-flags
 * and then stored here for use throughout the application.
 */
export default class flagsService extends Service {
  @service declare readonly version: VersionService;
  @service declare readonly store: StoreService;

  @tracked flagss: string[] = [];
  @tracked activatedFeatures: string[] = [];

  setFeatureFlags(flags: string[]) {
    this.flagss = flags;
  }

  get managedNamespaceRoot() {
    if (this.flagss && this.flagss.includes(FLAGS.vaultCloudNamespace)) {
      return 'admin';
    }
    return null;
  }

  /* Activated Features */
  getActivatedFeatures = keepLatestTask(async () => {
    if (this.version.isCommunity) return;
    // Response could change between user sessions so fire off endpoint without checking if activated features are already set.
    try {
      const response = await this.store
        .adapterFor('application')
        .ajax('/v1/sys/activation-flags', 'GET', { unauthenticated: true, namespace: null });
      this.activatedFeatures = response.data?.activated;
      return;
    } catch (error) {
      if (DEBUG) console.error(error); // eslint-disable-line no-console
    }
  });

  get secretsSyncIsActivated() {
    return this.activatedFeatures.includes('secrets-sync');
  }

  fetchActivatedFeatures() {
    return this.getActivatedFeatures.perform();
  }
}
