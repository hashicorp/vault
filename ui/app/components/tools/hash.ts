/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { service } from '@ember/service';
import { tracked } from '@glimmer/tracking';

import type ApiService from 'vault/services/api';
import type FlashMessageService from 'vault/services/flash-messages';
import type { HTMLElementEvent } from 'vault/forms';

/**
 * @module ToolsHash
 * ToolsHash components are components that sys/wrapping/hash functionality.
 *
 * @example
 * <Tools::Hash />
 */
export default class ToolsHash extends Component {
  @service declare readonly api: ApiService;
  @service declare readonly flashMessages: FlashMessageService;

  @tracked algorithm = 'sha2-256';
  @tracked format = 'base64';
  @tracked hashData = '';
  @tracked sum = '';
  @tracked errorMessage = '';

  get breadcrumbs() {
    return [{ label: 'Vault', route: 'vault.cluster.dashboard', icon: 'vault' }, { label: 'Hash data' }];
  }

  @action
  reset() {
    this.algorithm = 'sha2-256';
    this.format = 'base64';
    this.hashData = '';
    this.sum = '';
    this.errorMessage = '';
  }

  @action
  handleEvent(evt: HTMLElementEvent<HTMLInputElement>) {
    const { name, value } = evt.target;
    const key = name as 'algorithm' | 'format' | 'hashData';
    this[key] = value;
  }

  @action
  async handleSubmit(evt: HTMLElementEvent<HTMLFormElement>) {
    evt.preventDefault();
    const data = {
      input: this.hashData,
      format: this.format,
      algorithm: this.algorithm,
    };

    try {
      const { sum } = await this.api.sys.generateHash(data);
      this.sum = sum || '';
      this.flashMessages.success('Hash was successful.');
    } catch (error) {
      const { message } = await this.api.parseError(error);
      this.errorMessage = message;
    }
  }
}
