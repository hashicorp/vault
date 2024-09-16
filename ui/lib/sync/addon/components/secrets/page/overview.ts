/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { service } from '@ember/service';
import { task } from 'ember-concurrency';
import { action } from '@ember/object';
import Ember from 'ember';
import { DEBUG } from '@glimmer/env';

import type FlashMessageService from 'vault/services/flash-messages';
import type StoreService from 'vault/services/store';
import type VersionService from 'vault/services/version';
import type FlagsService from 'vault/services/flags';
import type { SyncDestinationAssociationMetrics } from 'vault/vault/adapters/sync/association';
import type SyncDestinationModel from 'vault/vault/models/sync/destination';

interface Args {
  destinations: Array<SyncDestinationModel>;
  totalVaultSecrets: number;
  isActivated: boolean;
  licenseHasSecretsSync: boolean;
  isHvdManaged: boolean;
}

export default class SyncSecretsDestinationsPageComponent extends Component<Args> {
  @service declare readonly flashMessages: FlashMessageService;
  @service declare readonly store: StoreService;
  @service declare readonly version: VersionService;
  @service declare readonly flags: FlagsService;

  @tracked destinationMetrics: SyncDestinationAssociationMetrics[] = [];
  @tracked page = 1;
  @tracked showActivateSecretsSyncModal = false;
  @tracked activationErrors: null | string[] = null;
  // eventually remove when we deal with permissions on activation-features
  @tracked hideOptIn = false;
  @tracked hideError = false;

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

  @action
  clearActivationErrors() {
    this.activationErrors = null;
  }

  @action
  onModalError(errorMsg: string) {
    if (DEBUG) console.error(errorMsg); // eslint-disable-line no-console

    const errors = [errorMsg];

    if (this.args.isHvdManaged) {
      errors.push(
        'Secrets Sync is available for Plus tier clusters only. Please check the tier of your cluster to enable Secrets Sync.'
      );
    }
    this.activationErrors = errors;
  }
}
