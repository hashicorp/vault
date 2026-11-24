/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
import { guidFor } from '@ember/object/internals';
import { run } from '@ember/runloop';
import { keyIsFolder, parentKeyForKey } from 'core/utils/key-utils';
import { SECRET_TYPE_FETCHERS } from 'core/utils/secret-type-fetchers';

import type ApiService from 'vault/services/api';
import type { SecretType } from 'vault/sync';

/**
 * @module SuggestionInput
 * @description
 * Generalized input component that fetches suggestions based on secret type
 * and displays them in a dropdown. Supports multiple secret types with
 * extensible configuration.
 *
 * @example
 * <SuggestionInput
 *   @type="kv"
 *   @label="Select a secret to sync"
 *   @subText="Enter the full path to the secret."
 *   @placeholder="Path to secret"
 *   @noMatchesMessage="No suggestions for this path"
 *   @value={{this.secretPath}}
 *   @mountPath="my-kv/"
 *   @onChange={{fn (mut this.secretPath)}}
 * />
 *
 * <SuggestionInput
 *   @type="database"
 *   @label="Select a static role"
 *   @placeholder="Static role name"
 *   @noMatchesMessage="No matching static roles found"
 *   @value={{this.roleName}}
 *   @mountPath="my-database/"
 *   @onChange={{fn (mut this.roleName)}}
 * />
 */

interface Args {
  type: SecretType;
  label: string;
  subText?: string;
  placeholder?: string;
  noMatchesMessage?: string;
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

export default class SuggestionInputComponent extends Component<Args> {
  @service declare readonly api: ApiService;

  @tracked suggestions: string[] = [];
  powerSelectAPI: PowerSelectAPI | undefined;
  _cachedSuggestions: string[] = [];
  _lastMountPath = '';
  inputId = `suggestion-input-${guidFor(this)}`;
  pathToSecret = ''; // only used for KV type

  constructor(owner: unknown, args: Args) {
    super(owner, args);
    if (this.args.mountPath) {
      this.updateSuggestions();
    }
  }

  get fetcher() {
    return SECRET_TYPE_FETCHERS[this.args.type];
  }

  get isDirectory() {
    return this.args.type === 'kv' && keyIsFolder(this.args.value);
  }

  get placeholder() {
    return this.args.placeholder ?? 'Search...';
  }

  get noMatchesMessage() {
    return this.args.noMatchesMessage ?? 'No matches found';
  }

  async fetchSuggestions() {
    const items = await this.fetcher.fetch(this.api, this.args.mountPath, this.args.value);

    if (this.args.type === 'kv') {
      const parentDirectory = parentKeyForKey(this.args.value);
      this.pathToSecret = this.isDirectory ? this.args.value : parentDirectory;
    }

    return items;
  }

  filterSuggestions(items: string[] | undefined = []) {
    return this.fetcher.filter(items, this.args.value, this.isDirectory);
  }

  @action
  async updateSuggestions() {
    // Reset cache when mount changes so the dropdown doesn't auto-open (consistent with initial selection)
    if (this.args.mountPath !== this._lastMountPath) {
      this._cachedSuggestions = [];
      this._lastMountPath = this.args.mountPath;
    }
    const isFirstUpdate = !this._cachedSuggestions.length;
    if (!this.args.mountPath) {
      this.suggestions = [];
    } else if (this.args.value && !this.isDirectory && this.suggestions) {
      this.suggestions = this.filterSuggestions(this._cachedSuggestions);
    } else {
      const items = await this.fetchSuggestions();
      this._cachedSuggestions = items;
      this.suggestions = this.filterSuggestions(items);
    }

    if (!isFirstUpdate) {
      const action = this.suggestions.length ? 'open' : 'close';
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
    if (this.suggestions.length) {
      this.powerSelectAPI?.actions.open();
    }
  }

  @action
  onSuggestionSelect(item: string) {
    const selectedValue = this.fetcher.onSelect(item, this.pathToSecret);
    this.args.onChange(selectedValue);
    this.updateSuggestions();
    run(() => document.getElementById(this.inputId)?.focus());
  }
}
