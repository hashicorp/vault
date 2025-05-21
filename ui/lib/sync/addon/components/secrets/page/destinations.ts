/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import { getOwner } from '@ember/owner';
import { findDestination, syncDestinations } from 'core/helpers/sync-destinations';
import { next } from '@ember/runloop';
import apiMethodResolver from 'sync/utils/api-method-resolver';

import type RouterService from '@ember/routing/router-service';
import type PaginationService from 'vault/services/pagination';
import type FlashMessageService from 'vault/services/flash-messages';
import type { CapabilitiesMap, EngineOwner } from 'vault/app-types';
import type { DestinationName, DestinationType, ListDestination } from 'vault/sync';
import type Transition from '@ember/routing/transition';
import type { PaginatedMetadata } from 'core/utils/paginate-list';
import type ApiService from 'vault/services/api';

interface Args {
  capabilities: CapabilitiesMap;
  destinations: ListDestination & PaginatedMetadata[];
  nameFilter: DestinationName;
  typeFilter: DestinationType;
}

export default class SyncSecretsDestinationsPageComponent extends Component<Args> {
  @service('app-router') declare readonly router: RouterService;
  @service declare readonly pagination: PaginationService;
  @service declare readonly flashMessages: FlashMessageService;
  @service declare readonly api: ApiService;

  @tracked destinationToDelete: ListDestination | null = null;
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
    const { typeFilter } = this.args;
    return typeFilter ? findDestination(typeFilter).name : undefined;
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
  async onDelete(destination: ListDestination) {
    try {
      const { name } = destination;
      const message = `Destination ${name} has been queued for deletion.`;
      const method = apiMethodResolver('delete', destination.type);
      await this.api.sys[method](destination.name, {});
      this.router.transitionTo('vault.cluster.sync.secrets.overview');
      this.flashMessages.success(message);
    } catch (error) {
      const { message } = await this.api.parseError(error);
      this.flashMessages.danger(`Error deleting destination \n ${message}`);
    } finally {
      this.destinationToDelete = null;
    }
  }
}
