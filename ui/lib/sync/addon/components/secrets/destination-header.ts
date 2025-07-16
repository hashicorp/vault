/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { service } from '@ember/service';
import apiMethodResolver from 'sync/utils/api-method-resolver';

import type RouterService from '@ember/routing/router-service';
import type PaginationService from 'vault/services/pagination';
import type FlashMessageService from 'vault/services/flash-messages';
import type ApiService from 'vault/services/api';
import type CapabilitiesService from 'vault/services/capabilities';
import type { Destination } from 'vault/sync';
import type { CapabilitiesMap } from 'vault/app-types';

interface Args {
  destination: Destination;
  capabilities: CapabilitiesMap;
}

export default class DestinationsTabsToolbar extends Component<Args> {
  @service('app-router') declare readonly router: RouterService;
  @service declare readonly pagination: PaginationService;
  @service declare readonly flashMessages: FlashMessageService;
  @service declare readonly api: ApiService;
  @service declare readonly capabilities: CapabilitiesService;

  get showSyncBtn() {
    const { destination, capabilities } = this.args;
    const path = this.capabilities.pathFor('syncSetAssociation', destination);
    return capabilities[path]?.canUpdate && !destination.purgeInitiatedAt;
  }

  get showEditBtn() {
    const { destination, capabilities } = this.args;
    const path = this.capabilities.pathFor('syncDestination', destination);
    return capabilities[path]?.canUpdate && !destination.purgeInitiatedAt;
  }

  @action
  async deleteDestination() {
    try {
      const { destination } = this.args;
      const message = `Destination ${destination.name} has been queued for deletion.`;
      const method = apiMethodResolver('delete', destination.type);
      await this.api.sys[method](destination.name, {});
      this.router.transitionTo('vault.cluster.sync.secrets.overview');
      this.flashMessages.success(message);
    } catch (error) {
      const { message } = await this.api.parseError(error);
      this.flashMessages.danger(`Error deleting destination \n ${message}`);
    }
  }
}
