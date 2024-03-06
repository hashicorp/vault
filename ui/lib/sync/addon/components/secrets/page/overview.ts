/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { service } from '@ember/service';
import { task } from 'ember-concurrency';
import Ember from 'ember';
import AdapterError from '@ember-data/adapter/error';

import type FlashMessageService from 'vault/services/flash-messages';
import type StoreService from 'vault/services/store';
import type VersionService from 'vault/services/version';
import type { SyncDestinationAssociationMetrics } from 'vault/vault/adapters/sync/association';
import type SyncDestinationModel from 'vault/vault/models/sync/destination';

interface Args {
  destinations: Array<SyncDestinationModel>;
  totalVaultSecrets: number;
  featureEnabled: boolean;
  adapterError: AdapterError | boolean;
}

export default class SyncSecretsDestinationsPageComponent extends Component<Args> {
  @service declare readonly flashMessages: FlashMessageService;
  @service declare readonly store: StoreService;
  @service declare readonly version: VersionService;

  @tracked destinationMetrics: SyncDestinationAssociationMetrics[] = [];
  @tracked page = 1;

  pageSize = Ember.testing ? 3 : 5; // lower in tests to test pagination without seeding more data

  constructor(owner: unknown, args: Args) {
    super(owner, args);
    if (this.args.destinations.length) {
      this.fetchAssociationsForDestinations.perform();
    }
  }

  fetchAssociationsForDestinations = task(this, {}, async (page = 1) => {
    try {
      const total = page * this.pageSize;
      const paginatedDestinations = this.args.destinations.slice(total - this.pageSize, total);
      this.destinationMetrics = await this.store
        .adapterFor('sync/association')
        .fetchByDestinations(paginatedDestinations);
      this.page = page;
    } catch (error) {
      this.destinationMetrics = [];
    }
  });
}
