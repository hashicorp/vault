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
 * @module ToolRandom
 * ToolRandom components are components that perform sys/wrapping/random functionality.
 * @example
 * <ToolRandom />
 */

export default class ToolRandom extends Component {
  @service store;
  @service flashMessages;

  @tracked bytes = 32;
  @tracked format = 'base64';
  @tracked random_bytes = null;
  @tracked errorMessage = '';

  @action
  reset() {
    this.bytes = 32;
    this.format = 'base64';
    this.random_bytes = null;
    this.errorMessage = '';
  }

  @action
  handleEvent(evt) {
    const { name, value } = evt.target;
    this[name] = value;
  }

  @action
  async handleSubmit(evt) {
    evt.preventDefault();
    const data = { bytes: parseInt(this.bytes), format: this.format };

    try {
      const response = await this.store.adapterFor('tools').toolAction('random', data);
      this.random_bytes = response.data.random_bytes;
      this.flashMessages.success('Generated random bytes successfully.');
    } catch (error) {
      this.errorMessage = errorMessage(error);
    }
  }
}
