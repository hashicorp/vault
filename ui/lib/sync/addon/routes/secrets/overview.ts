/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { hash } from 'rsvp';

import type FlagsService from 'vault/services/flags';
import type StoreService from 'vault/services/store';

export default class SyncSecretsOverviewRoute extends Route {
  @service declare readonly store: StoreService;
  @service declare readonly flags: FlagsService;

  async model() {
    const isActivated = this.flags.secretsSyncIsActivated;
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
