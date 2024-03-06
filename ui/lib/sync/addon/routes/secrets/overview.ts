/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { hash } from 'rsvp';
import AdapterError from 'ember-data/adapter';

import type StoreService from 'vault/services/store';

export default class SyncSecretsOverviewRoute extends Route {
  @service declare readonly store: StoreService;

  async model() {
    const { featureEnabled } = this.modelFor('secrets') as { featureEnabled: boolean };
    const { adapterError } = this.modelFor('secrets') as { adapterError: AdapterError | boolean };
    return hash({
      destinations: this.store.query('sync/destination', {}).catch(() => []),
      associations: this.store
        .adapterFor('sync/association')
        .queryAll()
        .catch(() => []),
      featureEnabled,
      adapterError,
    });
  }
}
