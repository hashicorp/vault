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
 * @module ToolsUnwrap
 * ToolsUnwrap components are components that sys/wrapping/rewrap functionality
 *
 * @example
 * <Tools::Unwrap />
 */

export default class ToolsUnwrap extends Component {
  @service declare readonly api: ApiService;
  @service declare readonly flashMessages: FlashMessageService;

  @tracked token = '';
  @tracked unwrapData: unknown = '';
  @tracked unwrapDetails = {};
  @tracked errorMessage = '';

  get breadcrumbs() {
    return [{ label: 'Vault', route: 'vault.cluster.dashboard', icon: 'vault' }, { label: 'Unwrap data' }];
  }

  @action
  reset() {
    this.token = '';
    this.unwrapData = '';
    this.unwrapDetails = {};
    this.errorMessage = '';
  }

  @action
  async handleSubmit(evt: HTMLElementEvent<HTMLFormElement>) {
    evt.preventDefault();
    const data = { token: this.token.trim() };

    try {
      const resp = await this.api.sys.unwrap(data);
      this.unwrapData = (resp && resp.data) || resp.auth;
      this.unwrapDetails = {
        'Request ID': resp.request_id,
        'Lease ID': resp.lease_id || 'None',
        Renewable: resp.renewable,
        'Lease Duration': resp.lease_duration || 'None',
      };
      this.flashMessages.success('Unwrap was successful.');
    } catch (error) {
      const { message } = await this.api.parseError(error);
      this.errorMessage = message;
    }
  }
}
