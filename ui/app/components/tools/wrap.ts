/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import { stringify } from 'core/helpers/stringify';

import type ApiService from 'vault/services/api';
import type FlashMessageService from 'vault/services/flash-messages';
import type { TtlEvent } from 'vault/app-types';
import type { HTMLElementEvent } from 'vault/forms';

/**
 * @module ToolsWrap
 * ToolsWrap components are components that sys/wrapping/wrap functionality.
 *
 * @example
 * <Tools::Wrap />
 */

export default class ToolsWrap extends Component {
  @service declare readonly api: ApiService;
  @service declare readonly flashMessages: FlashMessageService;

  @tracked hasLintingErrors = false;
  @tracked token = '';
  @tracked wrapTTL = '';
  @tracked wrapData = null;
  @tracked errorMessage = '';
  @tracked showJson = true;

  get breadcrumbs() {
    return [{ label: 'Vault', route: 'vault.cluster.dashboard', icon: 'vault' }, { label: 'Wrap data' }];
  }

  get startingValue() {
    // must pass the third param called "space" in JSON.stringify to structure object with whitespace
    // otherwise the following codemirror modifier check will pass `this._editor.getValue() !== namedArgs.content` and _setValue will be called.
    // the method _setValue moves the cursor to the beginning of the text field.
    // the effect is that the cursor jumps after the first key input.
    return stringify([{ '': '' }], { skipFormat: false });
  }

  get stringifiedWrapData() {
    return this.wrapData ? stringify([this.wrapData], { skipFormat: false }) : this.startingValue;
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
    this.wrapTTL = '';
    this.hasLintingErrors = false;
    if (clearData) this.wrapData = null;
  }

  @action
  updateTtl(evt: TtlEvent) {
    if (!evt) return;
    this.wrapTTL = evt.enabled ? `${evt.seconds}s` : '30m';
  }

  @action
  editorUpdated(val: string) {
    this.hasLintingErrors = false;

    try {
      this.wrapData = JSON.parse(val);
    } catch {
      this.hasLintingErrors = true;
    }
  }

  @action
  async handleSubmit(evt: HTMLElementEvent<HTMLFormElement>) {
    evt.preventDefault();

    const data = this.wrapData || {};
    const wrap = this.wrapTTL || '';

    try {
      const { wrap_info } = await this.api.sys.wrap(data, this.api.buildHeaders({ wrap }));
      this.token = wrap_info?.token || '';
      this.flashMessages.success('Wrap was successful.');
    } catch (error) {
      const { message } = await this.api.parseError(error);
      this.errorMessage = message;
    }
  }
}
