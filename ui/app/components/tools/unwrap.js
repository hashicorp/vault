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
 * @module ToolsUnwrap
 * ToolsUnwrap components are components that sys/wrapping/rewrap functionality
 *
 * @example
 * <Tools::Unwrap />
 */

export default class ToolsUnwrap extends Component {
  @service store;
  @service flashMessages;

  @tracked token = '';
  @tracked unwrapData = '';
  @tracked unwrapDetails = {};
  @tracked errorMessage = '';

  @action
  reset() {
    this.token = '';
    this.unwrapData = '';
    this.unwrapDetails = {};
    this.errorMessage = '';
  }

  @action
  async handleSubmit(evt) {
    evt.preventDefault();
    const data = { token: this.token.trim() };

    try {
      const resp = await this.store.adapterFor('tools').toolAction('unwrap', data);
      this.unwrapData = (resp && resp.data) || resp.auth;
      this.unwrapDetails = {
        'Request ID': resp.request_id,
        'Lease ID': resp.lease_id || 'None',
        Renewable: resp.renewable,
        'Lease Duration': resp.lease_duration || 'None',
      };
      this.flashMessages.success('Unwrap was successful.');
    } catch (error) {
      this.errorMessage = errorMessage(error);
    }
  }
}
