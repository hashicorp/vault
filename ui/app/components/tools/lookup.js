/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import errorMessage from 'vault/utils/error-message';
import { addSeconds, parseISO } from 'date-fns';

/**
 * @module ToolsLookup
 * ToolsLookup components are components that sys/wrapping/lookup functionality.
 *
 * @example
 * <Tools::Lookup />
 */
export default class ToolsLookup extends Component {
  @service store;
  @service flashMessages;

  @tracked token = '';
  @tracked lookupData = null;
  @tracked errorMessage = '';

  @action
  reset() {
    this.token = '';
    this.lookupData = null;
    this.errorMessage = '';
  }

  get expirationDate() {
    const { creation_time, creation_ttl } = this.lookupData;
    if (creation_time && creation_ttl) {
      // returns new Date with seconds added.
      return addSeconds(parseISO(creation_time), creation_ttl);
    }
    return null;
  }

  @action
  async handleSubmit(evt) {
    evt.preventDefault();
    const payload = { token: this.token.trim() };
    try {
      const resp = await this.store.adapterFor('tools').toolAction('lookup', payload);
      this.lookupData = resp.data;
      this.flashMessages.success('Lookup was successful.');
    } catch (error) {
      this.errorMessage = errorMessage(error);
    }
  }
}
