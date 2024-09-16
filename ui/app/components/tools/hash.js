/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import errorMessage from 'vault/utils/error-message';

/**
 * @module ToolsHash
 * ToolsHash components are components that sys/wrapping/hash functionality.
 *
 * @example
 * <Tools::Hash />
 */
export default class ToolsHash extends Component {
  @service store;
  @service flashMessages;

  @tracked algorithm = 'sha2-256';
  @tracked format = 'base64';
  @tracked hashData = '';
  @tracked sum = null;
  @tracked errorMessage = '';

  @action
  reset() {
    this.algorithm = 'sha2-256';
    this.format = 'base64';
    this.hashData = '';
    this.sum = null;
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
    const data = {
      input: this.hashData,
      format: this.format,
      algorithm: this.algorithm,
    };

    try {
      const response = await this.store.adapterFor('tools').toolAction('hash', data);
      this.sum = response.data.sum;
      this.flashMessages.success('Hash was successful.');
    } catch (error) {
      this.errorMessage = errorMessage(error);
    }
  }
}
