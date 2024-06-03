/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';

/**
 * @module ToolHash
 * ToolHash components are components that sys/wrapping/hash functionality.  Most of the functionality is passed through as actions from the tool-actions-form and then called back with properties.
 *
 * @example
 * <Tools::Hash @onClear={{action "onClear"}} @onChange={{action "onChange"}} @sum={{sum}} @algorithm={{algorithm}} @format={{format}} @errors={{errors}} />
 *
 * @param {Function} onClear - parent action that is passed through. Must be passed as {{action "onClear"}}
 * @param {String} sum=null - property passed from parent to child and then passed back up to parent.
 * @param {String} algorithm - property returned from parent.
 * @param {String} format - property returned from parent.
 * @param {Object} errors=null - errors passed from parent as default then from child back to parent.
 */
export default class ToolsHash extends Component {
  @action
  handleEvent(evt) {
    const { name, value } = evt.target;
    this.args.onChange(name, value);
  }
}
