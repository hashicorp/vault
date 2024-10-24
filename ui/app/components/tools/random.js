/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import errorMessage from 'vault/utils/error-message';

/**
 * @module ToolsRandom
 * ToolsRandom components are components that perform sys/wrapping/random functionality.
 * @example
 * <Tools::Random />
 */

export default class ToolsRandom extends Component {
  @service store;
  @service flashMessages;

  @tracked bytes = 32;
  @tracked format = 'base64';
  @tracked randomBytes = null;
  @tracked errorMessage = '';

  @action
  reset() {
    this.bytes = 32;
    this.format = 'base64';
    this.randomBytes = null;
    this.errorMessage = '';
  }

  @action
  handleSelect(evt) {
    const { value } = evt.target;
    this.format = value;
  }

  @action
  async handleSubmit(evt) {
    evt.preventDefault();
    const data = { bytes: Number(this.bytes), format: this.format };
    try {
      const response = await this.store.adapterFor('tools').toolAction('random', data);
      this.randomBytes = response.data.random_bytes;
      this.flashMessages.success('Generated random bytes successfully.');
    } catch (error) {
      this.errorMessage = errorMessage(error);
    }
  }
}
