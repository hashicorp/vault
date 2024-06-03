/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';

/**
 * @module ToolRandom
 * ToolRandom components are components that sys/wrapping/random functionality.  Most of the functionality is passed through as actions from the tool-actions-form and then called back with properties.
 *
 * @example
 * <Tools::Random
 *  @onClear={{action "onClear"}}
 *  @format={{format}}
 *  @bytes={{bytes}}
 *  @random_bytes={{random_bytes}}
 *  @errors={{errors}}/>
 *
 * @param {Function} onClear - parent action that is passed through. Must be passed as {{action "onClear"}}
 * @param {string} format - property returned from parent.
 * @param {string} bytes - property returned from parent.
 * @param {string} random_bytes - property returned from parent.
 * @param {object} errors=null - errors passed from parent as default then from child back to parent.
 */

export default class ToolsRandom extends Component {
  @action
  handleEvent(evt) {
    const { name, value } = evt.target;
    this.args.onChange(name, value);
  }
}
