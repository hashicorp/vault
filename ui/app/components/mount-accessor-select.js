/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';
import { task } from 'ember-concurrency';
import { action } from '@ember/object';

/**
 * @module MountAccessorSelect
 * The MountAccessorSelect component is used to selectDrop down mount options.
 *
 * @example
 * ```js
 * <MountAccessorSelect @value={this.aliasMountAccessor} @onChange={this.onChange} />
 * ```
 * @param {string} value - the selected value.
 * @param {function} onChange - the parent function that handles when a new value is selected.
 * @param {boolean} [showAccessor] - whether or not you should show the value or the more detailed accessor off the class.
 * @param {boolean} [noDefault] - whether or not there is a default value.
 * @param {boolean} [filterToken] - whether or not you should filter out type "token".
 * @param {string} [name] - name on the label.
 * @param {string} [label] - label above the select input.
 * @param {string} [helpText] - text shown in tooltip.
 */

export default class MountAccessorSelect extends Component {
  @service store;

  get filterToken() {
    return this.args.filterToken || false;
  }

  get noDefault() {
    return this.args.noDefault || false;
  }

  constructor() {
    super(...arguments);
    this.authMethods.perform();
  }

  @task *authMethods() {
    const methods = yield this.store.findAll('auth-method');
    if (!this.args.value && !this.args.noDefault) {
      const getValue = methods[0].accessor;
      this.args.onChange(getValue);
    }
    return methods;
  }

  @action change(event) {
    this.args.onChange(event.target.value);
  }
}
