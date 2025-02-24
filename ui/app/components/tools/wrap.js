/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import { stringify } from 'core/helpers/stringify';
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

  @tracked hasLintingErrors = false;
  @tracked token = '';
  @tracked wrapTTL = null;
  @tracked wrapData = null;
  @tracked errorMessage = '';
  @tracked showJson = true;

  get startingValue() {
    // must pass the third param called "space" in JSON.stringify to structure object with whitespace
    // otherwise the following codemirror modifier check will pass `this._editor.getValue() !== namedArgs.content` and _setValue will be called.
    // the method _setValue moves the cursor to the beginning of the text field.
    // the effect is that the cursor jumps after the first key input.
    return JSON.stringify({ '': '' }, null, 2);
  }

  get stringifiedWrapData() {
    return this?.wrapData ? stringify([this.wrapData], {}) : this.startingValue;
  }

  @action
  handleToggle() {
    this.showJson = !this.showJson;
    this.hasLintingErrors = false;
  }

  @action
  reset(clearData = true) {
    this.token = '';
    this.errorMessage = '';
    this.wrapTTL = null;
    this.hasLintingErrors = false;
    if (clearData) this.wrapData = null;
  }

  @action
  updateTtl(evt) {
    if (!evt) return;
    this.wrapTTL = evt.enabled ? `${evt.seconds}s` : '30m';
  }

  @action
  codemirrorUpdated(val, codemirror) {
    codemirror.performLint();
    this.hasLintingErrors = codemirror?.state.lint.marked?.length > 0;
    if (!this.hasLintingErrors) this.wrapData = JSON.parse(val);
  }

  @action
  async handleSubmit(evt) {
    evt.preventDefault();

    const data = this.wrapData;
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
