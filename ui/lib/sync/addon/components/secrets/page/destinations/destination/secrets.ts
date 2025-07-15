/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import { getOwner } from '@ember/owner';

import type RouterService from '@ember/routing/router-service';
import type ApiService from 'vault/services/api';
import type FlashMessageService from 'vault/services/flash-messages';
import type { EngineOwner, CapabilitiesMap } from 'vault/app-types';
import type { Destination, AssociatedSecret } from 'vault/sync';

interface Args {
  destination: Destination;
  associations: AssociatedSecret[];
  capabilities: CapabilitiesMap;
}

export default class SyncSecretsDestinationsPageComponent extends Component<Args> {
  @service('app-router') declare readonly router: RouterService;
  @service declare readonly api: ApiService;
  @service declare readonly flashMessages: FlashMessageService;

  @tracked secretToUnsync: AssociatedSecret | null = null;

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
    this.router.transitionTo(
      'vault.cluster.sync.secrets.destinations.destination.secrets',
      this.args.destination.type,
      this.args.destination.name
    );
  }

  @action
  async update(association: AssociatedSecret, operation: string) {
    try {
      const { name, type } = this.args.destination;
      const { mount, secretName } = association;
      const body = { mount, secretName };

      if (operation === 'set') {
        await this.api.sys.systemWriteSyncDestinationsTypeNameAssociationsSet(name, type, body);
      } else {
        await this.api.sys.systemWriteSyncDestinationsTypeNameAssociationsRemove(name, type, body);
      }
      const action: string = operation === 'set' ? 'Sync' : 'Unsync';
      this.flashMessages.success(`${action} operation initiated.`);
    } catch (error) {
      const { message } = await this.api.parseError(error);
      this.flashMessages.danger(`Sync operation error: \n ${message}`);
    } finally {
      this.secretToUnsync = null;
      this.refreshRoute();
    }
  }
}
