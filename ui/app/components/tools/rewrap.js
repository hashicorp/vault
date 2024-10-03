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
 * @module ToolsRewrap
 * ToolsRewrap components are components that sys/wrapping/rewrap functionality
 *
 * @example
 * <Tools::Rewrap />
 */

export default class ToolsRewrap extends Component {
  @service store;
  @service flashMessages;

  @tracked originalToken = '';
  @tracked rewrappedToken = '';
  @tracked errorMessage = '';

  @action
  reset() {
    this.originalToken = '';
    this.rewrappedToken = '';
    this.errorMessage = '';
  }

  @action
  async handleSubmit(evt) {
    evt.preventDefault();
    const data = { token: this.originalToken.trim() };

    try {
      const response = await this.store.adapterFor('tools').toolAction('rewrap', data);
      this.rewrappedToken = response.wrap_info.token;
      this.flashMessages.success('Rewrap was successful.');
    } catch (error) {
      this.errorMessage = errorMessage(error);
    }
  }
}
