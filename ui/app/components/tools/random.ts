/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';

import type ApiService from 'vault/services/api';
import type FlashMessageService from 'vault/services/flash-messages';
import type { HTMLElementEvent } from 'vault/forms';

/**
 * @module ToolsRandom
 * ToolsRandom components are components that perform sys/wrapping/random functionality.
 * @example
 * <Tools::Random />
 */

export default class ToolsRandom extends Component {
  @service declare readonly api: ApiService;
  @service declare readonly flashMessages: FlashMessageService;

  @tracked bytes = 32;
  @tracked format = 'base64';
  @tracked randomBytes = '';
  @tracked errorMessage = '';

  get breadcrumbs() {
    return [{ label: 'Vault', route: 'vault.cluster.dashboard', icon: 'vault' }, { label: 'Random bytes' }];
  }

  @action
  reset() {
    this.bytes = 32;
    this.format = 'base64';
    this.randomBytes = '';
    this.errorMessage = '';
  }

  @action
  handleSelect(evt: HTMLElementEvent<HTMLSelectElement>) {
    const { value } = evt.target;
    this.format = value;
  }

  @action
  async handleSubmit(evt: HTMLElementEvent<HTMLFormElement>) {
    evt.preventDefault();
    const data = { bytes: Number(this.bytes), format: this.format };
    try {
      const { random_bytes } = await this.api.sys.generateRandom(data);
      this.randomBytes = random_bytes || '';
      this.flashMessages.success('Generated random bytes successfully.');
    } catch (error) {
      const { message } = await this.api.parseError(error);
      this.errorMessage = message;
    }
  }
}
