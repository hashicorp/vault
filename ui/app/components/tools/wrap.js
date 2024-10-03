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
 * @module ToolsWrap
 * ToolsWrap components are components that sys/wrapping/wrap functionality.
 *
 * @example
 * <Tools::Wrap />
 */

export default class ToolsWrap extends Component {
  @service store;
  @service flashMessages;

  @tracked buttonDisabled = false;
  @tracked token = '';
  @tracked wrapTTL = null;
  @tracked wrapData = '{\n}';
  @tracked errorMessage = '';

  @action
  reset(clearData = true) {
    this.token = '';
    this.errorMessage = '';
    this.wrapTTL = null;
    if (clearData) this.wrapData = '{\n}';
  }

  @action
  updateTtl(evt) {
    if (!evt) return;
    this.wrapTTL = evt.enabled ? `${evt.seconds}s` : '30m';
  }

  @action
  codemirrorUpdated(val, codemirror) {
    codemirror.performLint();
    const hasErrors = codemirror?.state.lint.marked?.length > 0;
    this.buttonDisabled = hasErrors;
    if (!hasErrors) this.wrapData = val;
  }

  @action
  async handleSubmit(evt) {
    evt.preventDefault();
    const data = JSON.parse(this.wrapData);
    const wrapTTL = this.wrapTTL || null;

    try {
      const response = await this.store.adapterFor('tools').toolAction('wrap', data, { wrapTTL });
      this.token = response.wrap_info.token;
      this.flashMessages.success('Wrap was successful.');
    } catch (error) {
      this.errorMessage = errorMessage(error);
    }
  }
}
