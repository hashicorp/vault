/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';
import { action } from '@ember/object';
import { getOwner } from '@ember/application';
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
  @service declare readonly router: RouterService;
  @service declare readonly store: StoreService;
  @service declare readonly flashMessages: FlashMessageService;

  get mountPoint(): string {
    const owner = getOwner(this) as EngineOwner;
    return owner.mountPoint;
  }

  get paginationQueryParams() {
    return (page: number) => ({ page });
  }

  @action
  async update(association: SyncAssociationModel, operation: string) {
    try {
      await association.save({ adapterOptions: { action: operation } });
      // this message can be expanded after testing -- deliberately generic for now
      this.flashMessages.success(
        'Sync operation successfully initiated. Status will be updated on secret when complete.'
      );
    } catch (error) {
      this.flashMessages.danger(`Sync operation error: \n ${errorMessage(error)}`);
    }
  }
}
