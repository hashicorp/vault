/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { hash } from 'rsvp';

import type StoreService from 'vault/services/store';
import type VersionService from 'vault/services/version';
import type AdapterError from '@ember-data/adapter';

export default class SyncSecretsOverviewRoute extends Route {
  @service declare readonly store: StoreService;
  @service declare readonly version: VersionService;

  async model() {
    const { activatedFeatures, adapterError } = this.modelFor('secrets') as {
      activatedFeatures: Array<string>;
      adapterError: AdapterError;
      licenseFeatures: Array<string>;
    };
    return hash({
      destinations: this.store.query('sync/destination', {}).catch(() => []),
      associations: this.store
        .adapterFor('sync/association')
        .queryAll()
        .catch(() => []),
      activatedFeatures,
      adapterError,
      licenseFeatures: this.version.features,
    });
  }
}
