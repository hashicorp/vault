/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Service from '@ember/service';
import { service } from '@ember/service';

import type FeatureFlagService from 'vault/services/feature-flag';
import type VersionService from 'vault/services/version';
import type StoreService from 'vault/services/store';

/**
 * This service returns a persona which can be used to hide or show various states. Currently being used for Secrets Sync, but designed so that other persona's can be added.
 */

export default class PersonaService extends Service {
  @service declare readonly featureFlag: FeatureFlagService;
  @service declare readonly store: StoreService;
  @service declare readonly version: VersionService;

  get isManagedNamespaceRoot() {
    return this.featureFlag.managedNamespaceRoot ? true : false;
  }

  /**
 * Secret Sync persona options are: ['SHOW_ENTERPRISE_CTA, SHOW_PREMIUM_CTA, SHOW_ACTIVATION_CTA and SHOW_SECRETS_SYNC]. The persona return is based on the following criteria:
  * OSS versions of on-prem Vault cluster:
    1. - Secrets Sync is not available, return SHOW_ENTERPRISE_CTA.   
  * Managed cluster:
    2. - Secrets Sync is activated, return SHOW_SECRETS_SYNC.
    3. - Secrets Sync is not activated, return SHOW_ACTIVATION_CTA.
  * Enterprise versions of on-prem Vault cluster:
    2. - Secrets Sync is not on the license, return SHOW_PREMIUM_CTA.
    3. - Secrets Sync is on the license and not activated, return SHOW_ACTIVATION_CTA.
    4. - Secrets Sync is on the license and activated, return SHOW_SECRETS_SYNC.
*/

  get secretsSyncPersona() {
    if (!this.version.isEnterprise) return 'SHOW_ENTERPRISE_CTA';
    // Everything else is enterprise
    return this.isManagedNamespaceRoot
      ? this.version.secretsSyncIsActivated
        ? 'SHOW_SECRETS_SYNC'
        : 'SHOW_ACTIVATION_CTA'
      : this.version.hasSecretsSync
      ? 'SHOW_PREMIUM_CTA'
      : this.version.secretsSyncIsActivated
      ? 'SHOW_ENTERPRISE_CTA'
      : 'SHOW_ACTIVATION_CTA';
  }
}
