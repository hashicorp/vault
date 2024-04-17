/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Service, { service } from '@ember/service';

import type FlagService from 'vault/services/flags';
import type VersionService from 'vault/services/version';

/**
 * This service returns a persona which can be used to hide or show various states. Currently being used for Secrets Sync, but designed so that other personas can be added.
 */

export default class PersonaService extends Service {
  @service declare readonly flags: FlagService;
  @service declare readonly version: VersionService;

  /**
 * Secret Sync persona options are: ['SHOW_ENTERPRISE_CTA, SHOW_PREMIUM_CTA, SHOW_ACTIVATION_CTA and SHOW_SECRETS_SYNC]. The persona return is based on the following criteria:
  * Community/OSS cluster:
    1. - Secrets Sync is not available, return SHOW_ENTERPRISE_CTA. This will never show for managed clusters because they are always Enterprise.
  * Managed cluster:
    2. - Secrets Sync is activated, return SHOW_SECRETS_SYNC.
    3. - Secrets Sync is not activated, return SHOW_ACTIVATION_CTA.
  * On-prem Enterprise cluster:
    4. - Secrets Sync is on the license and not activated, return SHOW_ACTIVATION_CTA.
    5. - Secrets Sync is on the license and activated, return SHOW_SECRETS_SYNC.
    6. - Secrets Sync is not on the license, return SHOW_PREMIUM_CTA.
*/

  get secretsSyncPersona() {
    if (!this.version.isEnterprise) return 'SHOW_ENTERPRISE_CTA';
    // TODO - until HVD backend changes are made, we can't show secrets sync for managed clusters
    // if (this.flags.managedNamespaceRoot) {
    //   return this.flags.secretsSyncIsActivated ? `SHOW_SECRETS_SYNC` : `SHOW_ACTIVATION_CTA`;
    // }
    if (this.version.hasSecretsSync) {
      return this.flags.secretsSyncIsActivated ? `SHOW_SECRETS_SYNC` : `SHOW_ACTIVATION_CTA`;
    }
    // only option left is enterprise without it on their license so show premium CTA
    return 'SHOW_PREMIUM_CTA';
  }
}
