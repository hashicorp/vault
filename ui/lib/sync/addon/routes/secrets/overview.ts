/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { hash } from 'rsvp';

import type RouterService from '@ember/routing/router-service';
import type VersionService from 'vault/services/version';
import type StoreService from 'vault/services/store';

export default class SyncSecretsOverviewRoute extends Route {
  @service declare readonly router: RouterService;
  @service declare readonly store: StoreService;
  @service declare readonly version: VersionService;

  // ARG TODO return to
  // beforeModel(): void | Promise<unknown> {
  //   if (this.featureFlag.managedNamespaceRoot !== null) {
  //     this.router.transitionTo('vault.cluster.dashboard');
  //   }
  // }

  async model() {
    const { persona } = this.modelFor('secrets') as {
      persona: string;
    };

    return hash({
      persona,
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
