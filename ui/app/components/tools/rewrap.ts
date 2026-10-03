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
 * @module ToolsRewrap
 * ToolsRewrap components are components that sys/wrapping/rewrap functionality
 *
 * @example
 * <Tools::Rewrap />
 */

export default class ToolsRewrap extends Component {
  @service declare readonly api: ApiService;
  @service declare readonly flashMessages: FlashMessageService;

  @tracked originalToken = '';
  @tracked rewrappedToken = '';
  @tracked errorMessage = '';

  get breadcrumbs() {
    return [{ label: 'Vault', route: 'vault.cluster.dashboard', icon: 'vault' }, { label: 'Rewrap token' }];
  }

  @action
  reset() {
    this.originalToken = '';
    this.rewrappedToken = '';
    this.errorMessage = '';
  }

  @action
  async handleSubmit(evt: HTMLElementEvent<HTMLFormElement>) {
    evt.preventDefault();
    const data = { token: this.originalToken.trim() };

    try {
      const { wrap_info } = await this.api.sys.rewrap(data);
      this.rewrappedToken = wrap_info?.token || '';
      this.flashMessages.success('Rewrap was successful.');
    } catch (error) {
      const { message } = await this.api.parseError(error);
      this.errorMessage = message;
    }
  }
}
