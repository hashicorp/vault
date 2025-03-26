/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { hash } from 'rsvp';
import lazyCapabilities, { apiPath } from 'vault/macros/lazy-capabilities';

import type FlagsService from 'vault/services/flags';
import type RouterService from '@ember/routing/router-service';
import type Store from '@ember-data/store';
import type VersionService from 'vault/services/version';
import type CapabilitiesModel from 'vault/models/capabilities';

export default class SyncSecretsOverviewRoute extends Route {
  @service('app-router') declare readonly router: RouterService;
  @service declare readonly store: Store;
  @service declare readonly flags: FlagsService;
  @service declare readonly version: VersionService;

  @lazyCapabilities(apiPath`sys/activation-flags/secrets-sync/activate`)
  declare syncPath: CapabilitiesModel;

  async model() {
    const isActivated = this.flags.secretsSyncIsActivated;

    return hash({
      canActivateSecretsSync:
        this.syncPath.get('canCreate') !== false || this.syncPath.get('canUpdate') !== false,
      destinations: isActivated ? this.store.query('sync/destination', {}).catch(() => []) : [],
      associations: isActivated
        ? this.store
            .adapterFor('sync/association')
            .queryAll()
            .catch(() => [])
        : [],
    });
  }

  redirect() {
    if (!this.version.hasSecretsSync) {
      this.router.replaceWith('vault.cluster.dashboard');
    }
  }
}
