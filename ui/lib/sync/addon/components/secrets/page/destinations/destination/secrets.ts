/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import { getOwner } from '@ember/owner';
import errorMessage from 'vault/utils/error-message';

import SyncDestinationModel from 'vault/vault/models/sync/destination';
import type SyncAssociationModel from 'vault/vault/models/sync/association';
import type RouterService from '@ember/routing/router-service';
import type StoreService from 'vault/services/store';
import type FlashMessageService from 'vault/services/flash-messages';
import type { EngineOwner } from 'vault/vault/app-types';

interface Args {
  destination: SyncDestinationModel;
  associations: Array<SyncAssociationModel>;
}

export default class SyncSecretsDestinationsPageComponent extends Component<Args> {
  @service('app-router') declare readonly router: RouterService;
  @service declare readonly store: StoreService;
  @service declare readonly flashMessages: FlashMessageService;

  @tracked secretToUnsync: SyncAssociationModel | null = null;

  get mountPoint(): string {
    const owner = getOwner(this) as EngineOwner;
    return owner.mountPoint;
  }

  get paginationQueryParams() {
    return (page: number) => ({ page });
  }

  @action
  refreshRoute() {
    // refresh route to update displayed secrets
    this.store.clearDataset('sync/association');
    this.router.transitionTo(
      'vault.cluster.sync.secrets.destinations.destination.secrets',
      this.args.destination.type,
      this.args.destination.name
    );
  }

  @action
  async update(association: SyncAssociationModel, operation: string) {
    try {
      await association.save({ adapterOptions: { action: operation } });
      const action: string = operation === 'set' ? 'Sync' : 'Unsync';
      this.flashMessages.success(`${action} operation initiated.`);
    } catch (error) {
      this.flashMessages.danger(`Sync operation error: \n ${errorMessage(error)}`);
    } finally {
      this.secretToUnsync = null;
      this.refreshRoute();
    }
  }
}
