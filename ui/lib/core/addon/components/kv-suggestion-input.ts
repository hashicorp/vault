/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
import { guidFor } from '@ember/object/internals';
import { run } from '@ember/runloop';
import { keyIsFolder, parentKeyForKey, keyWithoutParentKey } from 'core/utils/key-utils';
import { SecretsApiKvV2ListListEnum } from '@hashicorp/vault-client-typescript';

import type ApiService from 'vault/services/api';

/**
 * @module KvSuggestionInput
 * @description
 * Input component that fetches secrets at a provided mount path and displays them as suggestions in a dropdown
 * As the user types the result set will be filtered providing suggestions for the user to select
 * After the input debounce wait time (500ms), if the value ends in a slash, secrets will be fetched at that path
 * The new result set will then be displayed in the dropdown as suggestions for the newly inputted path
 * Selecting a suggestion will append it to the input value
 * This allows the user to build a full path to a secret for the provided mount
 * This is useful for helping the user find deeply nested secrets given the path based policy system
 * If the user does not have list permission they are still able to enter a path to a secret but will not see suggestions
 * Input is disabled when mount path is not provided
 *
 * @example
 * <KvSuggestionInput @label="Select a secret to sync" @subText="Enter the full path to the secret. Suggestions will display below if permitted by policy." @value={{this.secretPath}} @mountPath="my-kv/" @onChange={{fn (mut this.secretPath)}} />
 *
 * <KvSuggestionInput @label="Select a secret to sync" @subText="Disabled because no mount path provided" @value={{this.secretPath}} @mountPath={{false}} @onChange={{fn (mut this.secretPath)}} />
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
  @service declare readonly api: ApiService;

  @tracked secrets: string[] = [];
  powerSelectAPI: PowerSelectAPI | undefined;
  _cachedSecrets: string[] = []; // cache the response for filtering purposes
  inputId = `suggestion-input-${guidFor(this)}`; // add unique segment to id in case multiple instances of component are used on the same page
  pathToSecret = ''; // keeps track of the full path to the secret as user builds it out

  constructor(owner: unknown, args: Args) {
    super(owner, args);
    if (this.args.mountPath) {
      this.updateSuggestions();
    }
  }

  get isDirectory() {
    return keyIsFolder(this.args.value);
  }

  async fetchSecrets() {
    try {
      const { mountPath } = this.args;
      const backend = keyIsFolder(mountPath) ? mountPath.slice(0, -1) : mountPath;
      const parentDirectory = parentKeyForKey(this.args.value);
      this.pathToSecret = this.isDirectory ? this.args.value : parentDirectory;
      // kvV2List => GET /:secret-mount-path/metadata/:secret_path/?list=true
      // This request can either list secrets at the mount root or for a specified :secret_path.
      // Since :secret_path already contains a trailing slash, e.g. /metadata/my-secret//
      // the request URL is sanitized by the api service to remove duplicate slashes.
      const { keys } = await this.api.secrets.kvV2List(
        this.pathToSecret,
        backend,
        SecretsApiKvV2ListListEnum.TRUE
      );
      // this will be used to filter the existing result set when the search term changes within the same path
      this._cachedSecrets = keys || [];
      return this._cachedSecrets;
    } catch (error) {
      console.log(error); // eslint-disable-line
      return [];
    }
  }

  filterSecrets(secrets: string[] | undefined = []) {
    const { value } = this.args;
    const secretName = keyWithoutParentKey(value) || '';
    return secrets.filter((path) => {
      if (!value || this.isDirectory) {
        return true;
      }
      if (secretName === path) {
        // don't show suggestion if it's currently selected
        return false;
      }
      return path.toLowerCase().includes(secretName.toLowerCase());
    });
  }

  @action
  async updateSuggestions() {
    const isFirstUpdate = !this._cachedSecrets.length;
    if (!this.args.mountPath) {
      this.secrets = [];
    } else if (this.args.value && !this.isDirectory && this.secrets) {
      // if we don't need to fetch from a new path, filter the previous result set with the updated search term
      this.secrets = this.filterSecrets(this._cachedSecrets);
    } else {
      const secrets = await this.fetchSecrets();
      this.secrets = this.filterSecrets(secrets);
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
  onSuggestionSelect(secret: string) {
    // user may partially type a value to filter result set and then select a suggestion
    // in this case the partially typed value must be replaced with suggestion value
    // pathToSecret contains the previous selections or typed path segments
    this.args.onChange(`${this.pathToSecret}${secret}`);
    this.updateSuggestions();
    // refocus the input after selection
    run(() => document.getElementById(this.inputId)?.focus());
  }
}
