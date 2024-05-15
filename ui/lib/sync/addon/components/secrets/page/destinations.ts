/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import { getOwner } from '@ember/application';
import errorMessage from 'vault/utils/error-message';
import { findDestination, syncDestinations } from 'core/helpers/sync-destinations';
import { next } from '@ember/runloop';

import type SyncDestinationModel from 'vault/vault/models/sync/destination';
import type RouterService from '@ember/routing/router-service';
import type StoreService from 'vault/services/store';
import type FlashMessageService from 'vault/services/flash-messages';
import type { EngineOwner } from 'vault/vault/app-types';
import type { SyncDestinationName, SyncDestinationType } from 'vault/vault/helpers/sync-destinations';
import type Transition from '@ember/routing/transition';

interface Args {
  destinations: Array<SyncDestinationModel>;
  nameFilter: SyncDestinationName;
  typeFilter: SyncDestinationType;
}

export default class SyncSecretsDestinationsPageComponent extends Component<Args> {
  @service declare readonly router: RouterService;
  @service declare readonly store: StoreService;
  @service declare readonly flashMessages: FlashMessageService;

  @tracked destinationToDelete: SyncDestinationModel | null = null;
  // for some reason there isn't a full page refresh happening when transitioning on filter change
  // when the transition happens it causes the FilterInput component to lose focus since it can only focus on didInsert
  // to work around this, verify that a transition from this route was completed and then focus the input
  constructor(owner: unknown, args: Args) {
    super(owner, args);
    this.router.on('routeDidChange', this.focusNameFilter);
  }

  willDestroy(): void {
    this.router.off('routeDidChange', this.focusNameFilter);
    super.willDestroy();
  }

  focusNameFilter(transition?: Transition) {
    const route = 'vault.cluster.sync.secrets.destinations.index';
    if (transition?.from?.name === route && transition?.to?.name === route) {
      next(() => document.getElementById('name-filter')?.focus());
    }
  }

  // typeFilter arg comes in as destination type but we need to pass the destination display name into the SearchSelect
  get typeFilterName() {
    return findDestination(this.args.typeFilter)?.name;
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
  onFilterChange(key: string, value: { id: string; name: string }[] | string | undefined) {
    const queryValue = Array.isArray(value) ? value[0]?.name : value;
    this.router.transitionTo('vault.cluster.sync.secrets.destinations', {
      queryParams: { [key]: queryValue },
    });
  }

  @action
  async onDelete(destination: SyncDestinationModel) {
    try {
      const { name } = destination;
      const message = `Destination ${name} has been queued for deletion.`;
      await destination.destroyRecord();
      this.store.clearDataset('sync/destination');
      this.router.transitionTo('vault.cluster.sync.secrets.overview');
      this.flashMessages.success(message);
    } catch (error) {
      this.flashMessages.danger(`Error deleting destination \n ${errorMessage(error)}`);
    } finally {
      this.destinationToDelete = null;
    }
  }
}
