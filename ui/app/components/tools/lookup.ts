/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import apiErrorMessage from 'vault/utils/api-error-message';
import { addSeconds } from 'date-fns';

import type ApiService from 'vault/services/api';
import type FlashMessageService from 'vault/services/flash-messages';
import type { HTMLElementEvent } from 'vault/forms';
import type { ReadWrappingPropertiesResponse } from '@hashicorp/vault-client-typescript';

/**
 * @module ToolsLookup
 * ToolsLookup components are components that sys/wrapping/lookup functionality.
 *
 * @example
 * <Tools::Lookup />
 */
export default class ToolsLookup extends Component {
  @service declare readonly api: ApiService;
  @service declare readonly flashMessages: FlashMessageService;

  @tracked token = '';
  @tracked lookupData: ReadWrappingPropertiesResponse | null = null;
  @tracked errorMessage = '';

  @action
  reset() {
    this.token = '';
    this.lookupData = null;
    this.errorMessage = '';
  }

  get expirationDate() {
    if (this.lookupData) {
      const { creationTime, creationTtl } = this.lookupData;
      if (creationTime && creationTtl) {
        // returns new Date with seconds added.
        return addSeconds(creationTime, Number(creationTtl));
      }
    }
    return null;
  }

  @action
  async handleSubmit(evt: HTMLElementEvent<HTMLFormElement>) {
    evt.preventDefault();
    const payload = { token: this.token.trim() };
    try {
      const data = await this.api.sys.readWrappingProperties(payload);
      this.lookupData = data;
      this.flashMessages.success('Lookup was successful.');
    } catch (error) {
      this.errorMessage = await apiErrorMessage(error);
    }
  }
}
