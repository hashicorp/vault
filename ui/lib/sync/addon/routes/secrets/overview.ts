/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { hash } from 'rsvp';

import type RouterService from '@ember/routing/router-service';
import type FeatureFlagService from 'vault/services/feature-flag';
import type VersionService from 'vault/services/version';
import type StoreService from 'vault/services/store';

enum secretsSyncPersona {
  SHOW_ENTERPRISE_CTA = 'SHOW_ENTERPRISE_CTA',
  SHOW_PREMIUM_CTA = 'SHOW_PREMIUM_CTA',
  SHOW_ACTIVATION_CTA = 'SHOW_ACTIVATION_CTA',
  SHOW_SECRETS_SYNC = 'SHOW_SECRETS_SYNC',
}

export default class SyncSecretsOverviewRoute extends Route {
  @service declare readonly router: RouterService;
  @service declare readonly store: StoreService;
  @service declare readonly featureFlag: FeatureFlagService;
  @service declare readonly version: VersionService;

  beforeModel(): void | Promise<unknown> {
    if (this.featureFlag.managedNamespaceRoot !== null) {
      this.router.transitionTo('vault.cluster.dashboard');
    }
  }

  async model() {
    const { secretsSyncPersona } = this.modelFor('secrets') as {
      secretsSyncPersona: secretsSyncPersona;
    };

    return hash({
      secretsSyncPersona,
      destinations: this.version.secretsSyncIsActivated
        ? this.store.query('sync/destination', {}).catch(() => [])
        : [],
      associations: this.version.secretsSyncIsActivated
        ? this.store
            .adapterFor('sync/association')
            .queryAll()
            .catch(() => [])
        : [],
    });
  }
}
