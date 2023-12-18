/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

/**
 * @module Toggle
 * Toggle components are used to indicate boolean values which can be toggled on or off.
 * They are a stylistic alternative to checkboxes, but still use the input[type="checkbox"] under the hood.
 *
 * @example
 * ```js
 * <Toggle @requiredParam={requiredParam} @optionalParam={optionalParam} @param1={{param1}}/>
 * ```
 * @param {function} onChange - onChange is triggered on checkbox change (select, deselect). Must manually mutate checked value
 * @param {string} name - name is passed along to the form field, as well as to generate the ID of the input & "for" value of the label
 * @param {boolean} [checked=false] - checked status of the input, and must be passed in and mutated from the parent
 * @param {boolean} [disabled=false] - disabled makes the switch unclickable
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';

export default class ToggleComponent extends Component {
  // tracked because the Input mutates the property and therefor cannot be a getter
  @tracked
  checked = this.args.checked || false;

  get disabled() {
    return this.args.disabled || false;
  }

  get name() {
    return this.args.name || '';
  }

  get safeId() {
    return `toggle-${this.name.replace(/\W/g, '')}`;
  }

  @action
  handleChange(e) {
    this.args.onChange(e.target.checked);
  }
}
