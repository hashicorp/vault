/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { service } from '@ember/service';
import { action } from '@ember/object';
import { task } from 'ember-concurrency';
import { keyIsFolder } from 'core/utils/key-utils';

import type { Destination } from 'vault/sync';
import type RouterService from '@ember/routing/router-service';
import type ApiService from 'vault/services/api';
import type PaginationService from 'vault/services/pagination';
import type FlashMessageService from 'vault/services/flash-messages';
import type Store from '@ember-data/store';
import type { SearchSelectOption } from 'vault/app-types';

interface Args {
  destination: Destination;
}

export default class DestinationSyncPageComponent extends Component<Args> {
  @service('app-router') declare readonly router: RouterService;
  @service declare readonly api: ApiService;
  @service declare readonly flashMessages: FlashMessageService;
  @service declare readonly pagination: PaginationService;
  @service declare readonly store: Store;

  constructor(owner: unknown, args: Args) {
    super(owner, args);
    this.fetchMounts();
  }

  @tracked mounts: SearchSelectOption[] = [];
  @tracked mountPath = '';
  @tracked secretPath = '';
  @tracked error = '';
  @tracked syncedSecret = '';

  get isSecretDirectory() {
    return this.secretPath && keyIsFolder(this.secretPath);
  }

  get isSubmitDisabled() {
    return !this.mountPath || !this.secretPath || this.isSecretDirectory || this.setAssociation.isRunning;
  }

  // unable to use built-in fetch functionality of SearchSelect since we need to filter by kv type
  async fetchMounts() {
    const adapter = this.store.adapterFor('application');
    const mounts = [];
    try {
      const { data } = await adapter.ajax('/v1/sys/internal/ui/mounts', 'GET');
      const secret = data.secret;
      if (secret) {
        for (const path in secret) {
          const { type, options } = secret[path as keyof typeof secret];
          if (type === 'kv' && options?.['version'] === '2') {
            mounts.push({ name: path, id: path });
          }
        }
      }
      this.mounts = mounts;
    } catch (error) {
      // the user is still able to manually enter the mount path
      // InputSearch component will render in this case
    }
  }

  @action
  setMount(selected: Array<string>) {
    this.mountPath = selected[0] || '';
    if (this.mountPath === '') {
      // reset form path when mount is cleared
      this.secretPath = '';
      this.syncedSecret = '';
    }
  }

  setAssociation = task({}, async (event: Event) => {
    event.preventDefault();
    this.error = ''; // reset error
    try {
      this.syncedSecret = '';
      const { name, type } = this.args.destination;
      const mount = keyIsFolder(this.mountPath) ? this.mountPath.slice(0, -1) : this.mountPath; // strip trailing slash from mount path
      const payload = { mount, secretName: this.secretPath };
      await this.api.sys.systemWriteSyncDestinationsTypeNameAssociationsSet(name, type, payload);
      this.syncedSecret = this.secretPath;
      // reset the secret path to help make it clear that the sync was successful
      this.secretPath = '';
    } catch (error) {
      const { message } = await this.api.parseError(error);
      this.error = `Sync operation error: \n ${message}`;
    }
  });
}
