/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Service, { inject as service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { keepLatestTask } from 'ember-concurrency';
import { macroCondition, isDevelopingApp } from '@embroider/macros';
import { ADMINISTRATIVE_NAMESPACE } from 'vault/services/namespace';

import type Store from '@ember-data/store';
import type VersionService from 'vault/services/version';

const FLAGS = {
  vaultCloudNamespace: 'VAULT_CLOUD_ADMIN_NAMESPACE',
};

/**
 * This service returns information about cluster flags. For now, the two available flags are from sys/internal/ui/feature-flags and sys/activation-flags.
 * The feature-flags endpoint returns VAULT_CLOUD_ADMIN_NAMESPACE which indicates that the Vault cluster is managed rather than self-managed.
 * The activation-flags endpoint returns which features are enabled.
 */

export default class FlagsService extends Service {
  @service declare readonly version: VersionService;
  @service declare readonly store: Store;

  @tracked activatedFlags: string[] = [];
  @tracked featureFlags: string[] = [];

  get isHvdManaged(): boolean {
    return this.featureFlags?.includes(FLAGS.vaultCloudNamespace);
  }

  // for non-managed clusters the root namespace path is technically an empty string so we return null
  get hvdManagedNamespaceRoot(): string | null {
    return this.isHvdManaged ? ADMINISTRATIVE_NAMESPACE : null;
  }

  getFeatureFlags = keepLatestTask(async () => {
    try {
      const result = await fetch('/v1/sys/internal/ui/feature-flags', {
        method: 'GET',
      });

      if (result.status === 200) {
        const body = await result.json();
        this.featureFlags = body.feature_flags || [];
      }
    } catch (error) {
      if (macroCondition(isDevelopingApp())) {
        console.error(error);
      }
    }
  });

  fetchFeatureFlags() {
    return this.getFeatureFlags.perform();
  }

  get secretsSyncIsActivated(): boolean {
    return this.activatedFlags.includes('secrets-sync');
  }

  getActivatedFlags = keepLatestTask(async () => {
    // Response could change between user sessions.
    // Fire off endpoint without checking if activated features are already set.
    if (this.version.isCommunity) return;
    try {
      const response = await this.store
        .adapterFor('application')
        .ajax('/v1/sys/activation-flags', 'GET', { unauthenticated: true, namespace: null });
      this.activatedFlags = response.data?.activated;
      return;
    } catch (error) {
      if (macroCondition(isDevelopingApp())) {
        console.error(error);
      }
    }
  });

  fetchActivatedFlags() {
    return this.getActivatedFlags.perform();
  }
}
