/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { hash } from 'rsvp';

import type RouterService from '@ember/routing/router-service';
import type FeatureFlagService from 'vault/services/feature-flag';
import type StoreService from 'vault/services/store';

export default class SyncSecretsOverviewRoute extends Route {
  @service declare readonly router: RouterService;
  @service declare readonly store: StoreService;
  @service declare readonly featureFlag: FeatureFlagService;

  beforeModel(): void | Promise<unknown> {
    if (this.featureFlag.managedNamespaceRoot !== null) {
      this.router.transitionTo('vault.cluster.dashboard');
    }
  }

  async model() {
    const { activatedFeatures } = this.modelFor('secrets') as {
      activatedFeatures: Array<string>;
    };
    const isActivated = activatedFeatures.includes('secrets-sync');
    return hash({
      isActivated,
      destinations: isActivated ? this.store.query('sync/destination', {}).catch(() => []) : [],
      associations: isActivated
        ? this.store
            .adapterFor('sync/association')
            .queryAll()
            .catch(() => [])
        : [],
    });
  }
}
