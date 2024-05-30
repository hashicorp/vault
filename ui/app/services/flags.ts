/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Service, { inject as service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { keepLatestTask } from 'ember-concurrency';
import { DEBUG } from '@glimmer/env';
import lazyCapabilities, { apiPath } from 'vault/macros/lazy-capabilities';
import type StoreService from 'vault/services/store';
import type VersionService from 'vault/services/version';
import type PermissionsService from 'vault/services/permissions';

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
  @service declare readonly permissions: PermissionsService;

  @tracked activatedFlags: string[] = [];
  @tracked featureFlags: string[] = [];

  get isHvdManaged(): boolean {
    return this.featureFlags?.includes(FLAGS.vaultCloudNamespace);
  }

  get hvdManagedNamespaceRoot(): string | null {
    return this.isHvdManaged ? 'admin' : null;
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
      if (DEBUG) console.error(error); // eslint-disable-line no-console
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
      if (DEBUG) console.error(error); // eslint-disable-line no-console
    }
  });

  fetchActivatedFlags() {
    return this.getActivatedFlags.perform();
  }

  @lazyCapabilities(apiPath`sys/activation-flags/secrets-sync/activate`) secretsSyncActivatePath;

  get canActivateSecretsSync() {
    return (
      this.secretsSyncActivatePath.get('canCreate') !== false ||
      this.secretsSyncActivatePath.get('canUpdate') !== false
    );
  }

  get showSecretsSync() {
    const isHvdManaged = this.isHvdManaged;
    const onLicense = this.version.hasSecretsSync;
    const isEnterprise = this.version.isEnterprise;
    const isActivated = this.secretsSyncIsActivated;

    if (!isEnterprise) return false;
    if (isHvdManaged) return true;
    if (isEnterprise && !onLicense) return false;
    if (isActivated) {
      // if the feature is activated but the user does not have permissions on the `sys/sync` endpoint, hide navigation link.
      return this.permissions.hasNavPermission('sync');
    }
    // only remaining option is Enterprise with Secrets Sync on the license but the feature is not activated. In this case, we want to show the upsell page and message about either activating or having an admin activate.
    return true;
  }
}
