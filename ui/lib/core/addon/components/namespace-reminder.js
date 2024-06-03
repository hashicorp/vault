/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';

/**
 * @module NamespaceReminder
 * Renders a namespace reminder, typically used when creating a new item.
 * _The namespace reminder only renders within a namespace, we cannot stub the namespace service here
 * so manually wrote the component in the **example** below so it renders in docfy_
 *
 * @example
 * <NamespaceReminder @mode="save" @noun="Auth Method" />
 *
 * <p class="namespace-reminder" id="namespace-reminder">
 *  This Auth Method will be saved in the <span class="tag">admin/</span>namespace.
 * </p>
 *
 * @param {string} noun - item being created by form
 * @param {string} [mode=edit] - action happening in form
 */
export default class NamespaceReminder extends Component {
  @service namespace;

  get showMessage() {
    return !this.namespace.inRootNamespace;
  }

  get mode() {
    return this.args.mode || 'edit';
  }

  get modeVerb() {
    if (!this.mode) {
      return '';
    }
    return this.mode.endsWith('e') ? `${this.mode}d` : `${this.mode}ed`;
  }
}
