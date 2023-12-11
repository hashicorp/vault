/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
import { guidFor } from '@ember/object/internals';
import { run } from '@ember/runloop';
import { keyIsFolder, parentKeyForKey, keyWithoutParentKey } from 'core/utils/key-utils';

import type StoreService from 'vault/services/store';
import type KvSecretMetadataModel from 'vault/models/kv/metadata';

/**
 * @module KvSuggestionInput
 * Input component that fetches secrets at a provided mount path and displays them as suggestions in a dropdown
 * As the user types the result set will be filtered providing suggestions for the user to select
 * After the input debounce wait time (500ms), if the value ends in a slash, secrets will be fetched at that path
 * The new result set will then be displayed in the dropdown as suggestions for the newly inputted path
 * Selecting a suggestion will append it to the input value
 * This allows the user to build a full path to a secret for the provided mount
 * This is useful for helping the user find deeply nested secrets given the path based policy system
 * If the user does not have list permission they are still able to enter a path to a secret but will not see suggestions
 * 
 * @example
 * <KvSuggestionInput
    @label="Select a secret to sync"
    @subText="Enter the full path to the secret. Suggestions will display below if permitted by policy."
    @value={{this.secretPath}}
    @mountPath={{this.mountPath}} // input disabled when mount path is not provided
    @onChange={{fn (mut this.secretPath)}}
  /> 
 */

interface Args {
  label: string;
  subText?: string;
  mountPath: string;
  value: string;
  onChange: CallableFunction;
}

interface PowerSelectAPI {
  actions: {
    open(): void;
    close(): void;
  };
}

export default class KvSuggestionInputComponent extends Component<Args> {
  @service declare readonly store: StoreService;

  @tracked secrets: KvSecretMetadataModel[] = [];
  powerSelectAPI: PowerSelectAPI | undefined;
  _cachedSecrets: KvSecretMetadataModel[] = []; // cache the response for filtering purposes
  inputId = `suggestion-input-${guidFor(this)}`; // add unique segment to id in case multiple instances of component are used on the same page

  constructor(owner: unknown, args: Args) {
    super(owner, args);
    if (this.args.mountPath) {
      this.updateSuggestions();
    }
  }

  async fetchSecrets(isDirectory: boolean) {
    const { mountPath } = this.args;
    try {
      const backend = keyIsFolder(mountPath) ? mountPath.slice(0, -1) : mountPath;
      const parentDirectory = parentKeyForKey(this.args.value);
      const pathToSecret = isDirectory ? this.args.value : parentDirectory;
      const kvModels = (await this.store.query('kv/metadata', {
        backend,
        pathToSecret,
      })) as unknown;
      // this will be used to filter the existing result set when the search term changes within the same path
      this._cachedSecrets = kvModels as KvSecretMetadataModel[];
      return this._cachedSecrets;
    } catch (error) {
      console.log(error); // eslint-disable-line
      return [];
    }
  }

  filterSecrets(kvModels: KvSecretMetadataModel[] | undefined = [], isDirectory: boolean) {
    const { value } = this.args;
    const secretName = keyWithoutParentKey(value) || '';
    return kvModels.filter((model) => {
      if (!value || isDirectory) {
        return true;
      }
      if (value === model.fullSecretPath) {
        // don't show suggestion if it's currently selected
        return false;
      }
      return model.path.toLowerCase().includes(secretName.toLowerCase());
    });
  }

  @action
  async updateSuggestions() {
    const isFirstUpdate = !this._cachedSecrets.length;
    const isDirectory = keyIsFolder(this.args.value);
    if (!this.args.mountPath) {
      this.secrets = [];
    } else if (this.args.value && !isDirectory && this.secrets) {
      // if we don't need to fetch from a new path, filter the previous result set with the updated search term
      this.secrets = this.filterSecrets(this._cachedSecrets, isDirectory);
    } else {
      const kvModels = await this.fetchSecrets(isDirectory);
      this.secrets = this.filterSecrets(kvModels, isDirectory);
    }
    // don't do anything on first update -- allow dropdown to open on input click
    if (!isFirstUpdate) {
      const action = this.secrets.length ? 'open' : 'close';
      this.powerSelectAPI?.actions[action]();
    }
  }

  @action
  onInput(value: string) {
    this.args.onChange(value);
    this.updateSuggestions();
  }

  @action
  onInputClick() {
    if (this.secrets.length) {
      this.powerSelectAPI?.actions.open();
    }
  }

  @action
  onSuggestionSelect(secret: KvSecretMetadataModel) {
    // user may partially type a value to filter result set and then select a suggestion
    // in this case the partially typed value must be replaced with suggestion value
    // the fullSecretPath contains the previous selections or typed path segments
    this.args.onChange(secret.fullSecretPath);
    this.updateSuggestions();
    // refocus the input after selection
    run(() => document.getElementById(this.inputId)?.focus());
  }
}
