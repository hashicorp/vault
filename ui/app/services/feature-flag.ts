/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Service from '@ember/service';
import { tracked } from '@glimmer/tracking';

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
export default class FeatureFlagService extends Service {
  @tracked featureFlags: string[] = [];

  setFeatureFlags(flags: string[]) {
    this.featureFlags = flags;
  }

  get managedNamespaceRoot() {
    if (this.featureFlags && this.featureFlags.includes(FLAGS.vaultCloudNamespace)) {
      return 'admin';
    }
    return null;
  }
}
