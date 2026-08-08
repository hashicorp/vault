/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import { getOwner } from '@ember/owner';
import { SECRET_TYPE_CONFIGS, getSecretTypeFromAccessor } from 'sync/utils/secret-type-config';

import type RouterService from '@ember/routing/router-service';
import type ApiService from 'vault/services/api';
import type FlashMessageService from 'vault/services/flash-messages';
import type { EngineOwner, CapabilitiesMap } from 'vault/app-types';
import type { Destination, AssociatedSecret, SecretType } from 'vault/sync';

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
  getSecretType(association: AssociatedSecret): SecretType | null {
    return association.accessor ? getSecretTypeFromAccessor(association.accessor) : null;
  }

  getConfig(association: AssociatedSecret) {
    const secretType = this.getSecretType(association);
    return secretType ? SECRET_TYPE_CONFIGS[secretType] : null;
  }

  @action
  getRoute(association: AssociatedSecret): string {
    return this.getConfig(association)?.route || 'kvSecretOverview';
  }

  @action
  getModels(association: AssociatedSecret): string[] {
    const config = this.getConfig(association);
    if (config) {
      return config.getModels(association.mount, association.secret_name);
    }
    return [association.mount, association.secret_name];
  }

  @action
  getQuery(association: AssociatedSecret): Record<string, string> | undefined {
    return this.getConfig(association)?.getQuery?.();
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
      const { mount, secret_name } = association;
      const body = { mount, secret_name };

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
