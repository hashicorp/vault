/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import { hash } from 'rsvp';

import type StoreService from 'vault/services/store';
import type AdapterError from '@ember-data/adapter';

export default class SyncSecretsOverviewRoute extends Route {
  @service declare readonly store: StoreService;

  async model() {
    const { activatedFeatures } = this.modelFor('secrets') as {
      activatedFeatures: Array<string> | AdapterError;
    };
    return hash({
      destinations: this.store.query('sync/destination', {}).catch(() => []),
      associations: this.store
        .adapterFor('sync/association')
        .queryAll()
        .catch(() => []),
      activatedFeatures,
    });
  }
}
