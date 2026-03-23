/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Service, { inject as service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { keepLatestTask } from 'ember-concurrency';
import { macroCondition, isDevelopingApp } from '@embroider/macros';
import { ADMINISTRATIVE_NAMESPACE } from 'vault/services/namespace';

import type VersionService from 'vault/services/version';
import type ApiService from 'vault/services/api';

const FLAGS = {
  vaultCloudNamespace: 'VAULT_CLOUD_ADMIN_NAMESPACE',
};

export type ActivationFlags = {
  activated: string[];
  unactivated: string[];
};

/**
 * This service returns information about cluster flags. For now, the two available flags are from sys/internal/ui/feature-flags and sys/activation-flags.
 * The feature-flags endpoint returns VAULT_CLOUD_ADMIN_NAMESPACE which indicates that the Vault cluster is managed rather than self-managed.
 * The activation-flags endpoint returns which features are enabled.
 */

export default class FlagsService extends Service {
  @service declare readonly version: VersionService;
  @service declare readonly api: ApiService;

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
      // unable to use internalUiListEnabledFeatureFlags method since the response does not conform to expected format
      // example -> { feature_flags: string[] } instead of the standard { data: { feature_flags: string[] } }
      // since it is typed as JSONApiResponse and not VoidResponse the client attempts to parse the body at
      const response = await this.api.request.get('/sys/internal/ui/feature-flags');
      const { feature_flags } = await response.json();
      this.featureFlags = feature_flags || [];
    } catch (error) {
      const { response } = await this.api.parseError(error);
      if (macroCondition(isDevelopingApp())) {
        console.error(response);
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
      const { data } = await this.api.sys.readActivationFlags(
        this.api.buildHeaders({ token: '', namespace: '' })
      );
      this.activatedFlags = (data as ActivationFlags)?.activated;
      return;
    } catch (error) {
      const { response } = await this.api.parseError(error);
      if (macroCondition(isDevelopingApp())) {
        console.error(response);
      }
    }
  });

  fetchActivatedFlags() {
    return this.getActivatedFlags.perform();
  }
}
