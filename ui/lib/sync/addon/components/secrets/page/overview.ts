/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { service } from '@ember/service';
import { task } from 'ember-concurrency';
import { waitFor } from '@ember/test-waiters';
import { action } from '@ember/object';
import errorMessage from 'vault/utils/error-message';
import Ember from 'ember';

import type FlashMessageService from 'vault/services/flash-messages';
import type StoreService from 'vault/services/store';
import type RouterService from '@ember/routing/router-service';
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
  @service declare readonly router: RouterService;
  @service declare readonly version: VersionService;
  @service declare readonly flags: FlagsService;

  @tracked destinationMetrics: SyncDestinationAssociationMetrics[] = [];
  @tracked page = 1;
  @tracked showActivateSecretsSyncModal = false;
  @tracked hasConfirmedDocs = false;
  @tracked error = null;
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
  resetOptInModal() {
    this.showActivateSecretsSyncModal = false;
    this.hasConfirmedDocs = false;
  }

  @task
  @waitFor
  *onFeatureConfirm() {
    // must return null instead of root for non managed cluster.
    // child namespaces are not sent.
    const namespace = this.args.isHvdManaged ? 'admin' : null;
    try {
      yield this.store
        .adapterFor('application')
        .ajax('/v1/sys/activation-flags/secrets-sync/activate', 'POST', { namespace });
      // must refresh and not transition because transition does not refresh the model from within a namespace
      yield this.router.refresh();
    } catch (error) {
      this.error = errorMessage(error);
      this.flashMessages.danger(`Error enabling feature \n ${errorMessage(error)}`);
    } finally {
      this.resetOptInModal();
    }
  }
}
