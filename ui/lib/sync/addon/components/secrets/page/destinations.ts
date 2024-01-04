/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';
import { action } from '@ember/object';
import { getOwner } from '@ember/application';
import errorMessage from 'vault/utils/error-message';
import { findDestination, syncDestinations } from 'core/helpers/sync-destinations';

import type SyncDestinationModel from 'vault/vault/models/sync/destination';
import type RouterService from '@ember/routing/router-service';
import type StoreService from 'vault/services/store';
import type FlashMessageService from 'vault/services/flash-messages';
import type { EngineOwner } from 'vault/vault/app-types';
import type { SyncDestinationName, SyncDestinationType } from 'vault/vault/helpers/sync-destinations';

interface Args {
  destinations: Array<SyncDestinationModel>;
  nameFilter: SyncDestinationName;
  typeFilter: SyncDestinationType;
}

export default class SyncSecretsDestinationsPageComponent extends Component<Args> {
  @service declare readonly router: RouterService;
  @service declare readonly store: StoreService;
  @service declare readonly flashMessages: FlashMessageService;

  // typeFilter arg comes in as destination type but we need to pass the destination display name into the SearchSelect
  get typeFilterName() {
    return findDestination(this.args.typeFilter)?.name;
  }

  get destinationNames() {
    return this.args.destinations.map((destination) => ({ id: destination.name, name: destination.name }));
  }

  get destinationTypes() {
    return syncDestinations().map((d) => ({ id: d.name, name: d.type }));
  }

  get mountPoint(): string {
    const owner = getOwner(this) as EngineOwner;
    return owner.mountPoint;
  }

  get paginationQueryParams() {
    return (page: number) => ({ page });
  }

  get noResultsMessage() {
    const { nameFilter, typeFilter } = this.args;
    if (nameFilter && typeFilter) {
      return `There are no ${this.typeFilterName || typeFilter} destinations matching "${nameFilter}".`;
    }
    if (nameFilter) {
      return `There are no destinations matching "${nameFilter}".`;
    }
    if (typeFilter) {
      return `There are no ${this.typeFilterName || typeFilter} destinations.`;
    }
    return '';
  }

  @action
  onFilterChange(key: string, selectObject: Array<{ id: string; name: string } | undefined>) {
    this.router.transitionTo('vault.cluster.sync.secrets.destinations', {
      queryParams: { [key]: selectObject[0]?.name },
    });
  }

  @action
  async onDelete(destination: SyncDestinationModel) {
    try {
      const { name } = destination;
      const message = `Successfully deleted destination ${name}.`;
      await destination.destroyRecord();
      this.store.clearDataset('sync/destination');
      this.router.transitionTo('vault.cluster.sync.secrets.destinations');
      this.flashMessages.success(message);
    } catch (error) {
      this.flashMessages.danger(`Error deleting destination \n ${errorMessage(error)}`);
    }
  }
}
