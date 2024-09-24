/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { service } from '@ember/service';
import errorMessage from 'vault/utils/error-message';

import type SyncDestinationModel from 'vault/models/sync/destination';
import type RouterService from '@ember/routing/router-service';
import type StoreService from 'vault/services/store';
import type FlashMessageService from 'vault/services/flash-messages';

interface Args {
  destination: SyncDestinationModel;
}

export default class DestinationsTabsToolbar extends Component<Args> {
  @service declare readonly router: RouterService;
  @service declare readonly store: StoreService;
  @service declare readonly flashMessages: FlashMessageService;

  @action
  async deleteDestination() {
    try {
      const { destination } = this.args;
      const message = `Destination ${destination.name} has been queued for deletion.`;
      await destination.destroyRecord();
      this.store.clearDataset('sync/destination');
      this.router.transitionTo('vault.cluster.sync.secrets.overview');
      this.flashMessages.success(message);
    } catch (error) {
      this.flashMessages.danger(`Error deleting destination \n ${errorMessage(error)}`);
    }
  }
}
