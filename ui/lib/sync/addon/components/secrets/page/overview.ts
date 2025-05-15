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
import { findDestination } from 'core/helpers/sync-destinations';

import type FlashMessageService from 'vault/services/flash-messages';
import type VersionService from 'vault/services/version';
import type FlagsService from 'vault/services/flags';
import type ApiService from 'vault/services/api';
import type { SystemReadSyncDestinationsTypeNameAssociationsResponse } from '@hashicorp/vault-client-typescript';
import type { ListDestination, DestinationMetrics, AssociatedSecret, DestinationType } from 'vault/sync';

interface Args {
  destinations: ListDestination[];
  totalVaultSecrets: number;
  canActivateSecretsSync: boolean;
}

export default class SyncSecretsDestinationsPageComponent extends Component<Args> {
  @service declare readonly flashMessages: FlashMessageService;
  @service declare readonly api: ApiService;
  @service declare readonly version: VersionService;
  @service declare readonly flags: FlagsService;

  @tracked destinationMetrics: DestinationMetrics[] = [];
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
      const requests = paginatedDestinations.map(({ name, type }) => {
        return this.api.sys.systemReadSyncDestinationsTypeNameAssociations(name, type);
      });
      const responses = await Promise.all(requests);
      this.destinationMetrics = this.normalizeFetchByDestinations(responses);
      this.page = page;
    } catch (error) {
      this.destinationMetrics = [];
    }
  });

  normalizeFetchByDestinations(
    responses: SystemReadSyncDestinationsTypeNameAssociationsResponse[]
  ): DestinationMetrics[] {
    return responses.map((response) => {
      const { storeName, storeType, associatedSecrets } = response;
      const type = storeType as DestinationType;
      const secrets = associatedSecrets as Record<string, AssociatedSecret>;
      const unsynced = [];
      let lastUpdated;

      for (const key in secrets) {
        const association = secrets[key];
        // for display purposes, any status other than SYNCED is considered unsynced
        if (association) {
          if (association.syncStatus !== 'SYNCED') {
            unsynced.push(association.syncStatus);
          }
          // use the most recent updated_at value as the last synced date
          const updated = new Date(association.updatedAt);
          if (!lastUpdated || updated > lastUpdated) {
            lastUpdated = updated;
          }
        }
      }

      const associationCount = Object.entries(secrets).length;
      return {
        icon: findDestination(type).icon,
        name: storeName,
        type,
        associationCount,
        status: associationCount ? (unsynced.length ? `${unsynced.length} Unsynced` : 'All synced') : null,
        lastUpdated,
      };
    });
  }

  @action
  clearActivationErrors() {
    this.activationErrors = null;
  }

  @action
  onModalError(errorMsg: string) {
    if (DEBUG) console.error(errorMsg); // eslint-disable-line no-console

    const errors = [errorMsg];

    if (this.flags.isHvdManaged) {
      errors.push(
        'Secrets Sync is available for Plus tier clusters only. Please check the tier of your cluster to enable Secrets Sync.'
      );
    }
    this.activationErrors = errors;
  }
}
