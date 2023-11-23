/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { service } from '@ember/service';
import { action } from '@ember/object';
import { task } from 'ember-concurrency';
import { keyIsFolder, parentKeyForKey, keyWithoutParentKey } from 'core/utils/key-utils';
import errorMessage from 'vault/utils/error-message';

import type SyncDestinationModel from 'vault/models/sync/destination';
import type KvSecretMetadataModel from 'vault/models/kv/metadata';
import type RouterService from '@ember/routing/router-service';
import type StoreService from 'vault/services/store';
import type FlashMessageService from 'vault/services/flash-messages';
import type { SearchSelectOption } from 'vault/vault/app-types';

interface Args {
  destination: SyncDestinationModel;
}
interface PowerSelectAPI {
  actions: {
    open(): void;
  };
}

export default class DestinationSyncPageComponent extends Component<Args> {
  @service declare readonly router: RouterService;
  @service declare readonly store: StoreService;
  @service declare readonly flashMessages: FlashMessageService;

  constructor(owner: unknown, args: Args) {
    super(owner, args);
    this.fetchMounts();
  }

  @tracked mounts: SearchSelectOption[] = [];
  @tracked secrets: KvSecretMetadataModel[] = [];
  @tracked mountPath = '';
  @tracked secretPath = '';
  @tracked error = '';
  powerSelectAPI: PowerSelectAPI | undefined;
  isShowingSuggestions = false;
  _lastSecretFetch: KvSecretMetadataModel[] | undefined; // cache the response for filtering purposes

  // strip trailing slash from mount path
  get mountName() {
    return this.mountPath ? this.mountPath.slice(0, -1) : null;
  }

  // unable to use built-in fetch functionality of SearchSelect since we need to filter by kv type
  async fetchMounts() {
    try {
      const secretEngines = await this.store.query('secret-engine', {});
      this.mounts = secretEngines.reduce((filtered, model) => {
        if (model.type === 'kv') {
          filtered.push({ name: model.path, id: model.path });
        }
        return filtered;
      }, []);
    } catch (error) {
      // the user is still able to manually enter the mount path
      // InputSearch component will render in this case
    }
  }

  async fetchSecrets(isDirectory: boolean) {
    try {
      const parentDirectory = parentKeyForKey(this.secretPath);
      const pathToSecret = isDirectory ? this.secretPath : parentDirectory;
      const kvModels = (await this.store.query('kv/metadata', {
        backend: this.mountName,
        pathToSecret,
      })) as unknown;
      // this will be used to filter the existing result set when the search term changes within the same path
      this._lastSecretFetch = kvModels as KvSecretMetadataModel[];
      return this._lastSecretFetch;
    } catch (error) {
      return [];
    }
  }

  filterSecrets(kvModels: KvSecretMetadataModel[] | undefined = [], isDirectory: boolean) {
    const secretName = keyWithoutParentKey(this.secretPath) || '';
    return kvModels.filter((model) => {
      if (!this.secretPath || isDirectory) {
        return true;
      }
      if (this.secretPath === model.fullSecretPath) {
        // don't show suggestion if it's currently selected
        return false;
      }
      return secretName.toLowerCase().includes(model.path.toLowerCase());
    });
  }

  async updateSecretSuggestions() {
    const isDirectory = keyIsFolder(this.secretPath);
    if (!this.mountPath) {
      this.secrets = [];
    } else if (this.secretPath && !isDirectory && this.secrets) {
      // if we don't need to fetch from a new path, filter the previous result set with the updated search term
      this.secrets = this.filterSecrets(this._lastSecretFetch, isDirectory);
    } else {
      const kvModels = await this.fetchSecrets(isDirectory);
      this.secrets = this.filterSecrets(kvModels, isDirectory);
    }
  }

  @action
  cancel() {
    this.router.transitionTo('vault.cluster.sync.secrets.destinations.destination.secrets');
  }

  @action
  setMount(selected: Array<string>) {
    this.mountPath = selected[0] || '';
    this.updateSecretSuggestions();
  }

  @action
  onSecretInput(value: string) {
    this.secretPath = value;
    this.updateSecretSuggestions();
  }

  @action
  onSecretInputClick() {
    this.powerSelectAPI?.actions?.open();
  }

  @action
  onSuggestionSelect(secret: KvSecretMetadataModel) {
    this.secretPath = this.secretPath + secret.path;
    this.updateSecretSuggestions();
  }

  @task
  *setAssociation(event: Event) {
    event.preventDefault();
    try {
      const { name, type } = this.args.destination;
      const association = this.store.createRecord('sync/association', {
        destinationName: name,
        destinationType: type,
        mount: this.mountName,
        secretName: this.secretPath,
      });
      yield association.save({ adapterOptions: { action: 'set' } });
      // this message can be expanded after testing -- deliberately generic for now
      this.flashMessages.success(
        'Sync operation successfully initiated. Status will be updated on secret when complete.'
      );
    } catch (error) {
      this.error = `Sync operation error: \n ${errorMessage(error)}`;
    }
  }
}
