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
 * This service returns information about cluster flags. For now, the two available flags are from sys/internal/ui/feature-flags and sys/activation-flags.
 * The feature-flags endpoint returns VAULT_CLOUD_ADMIN_NAMESPACE which indicates that the Vault cluster is managed rather than self-managed.
 * The activation-flags endpoint returns which features are enabled.
 */

export default class flagsService extends Service {
  @service declare readonly version: VersionService;
  @service declare readonly store: StoreService;

  @tracked flags: string[] = [];
  @tracked activatedFlags: string[] = [];

  setFeatureFlags(flags: string[]) {
    this.flags = flags;
  }

  get managedNamespaceRoot() {
    if (this.flags && this.flags.includes(FLAGS.vaultCloudNamespace)) {
      return 'admin';
    }
    return null;
  }

  // TODO getter will be used in the upcoming persona service
  get secretsSyncIsActivated() {
    return this.activatedFlags.includes('secrets-sync');
  }

  getActivatedFlags = keepLatestTask(async () => {
    if (this.version.isCommunity) return;
    // Response could change between user sessions.
    // Fire off endpoint without checking if activated features are already set.
    try {
      const response = await this.store
        .adapterFor('application')
        .ajax('/v1/sys/activation-flags', 'GET', { unauthenticated: true, namespace: null });
      this.activatedFlags = response.data?.activated;
      return;
    } catch (error) {
      if (DEBUG) console.error(error); // eslint-disable-line no-console
    }
  });

  fetchActivatedFlags() {
    return this.getActivatedFlags.perform();
  }
}
