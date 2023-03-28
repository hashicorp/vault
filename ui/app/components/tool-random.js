/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';

/**
 * @module ToolRandom
 * ToolRandom components are components that sys/wrapping/random functionality.  Most of the functionality is passed through as actions from the tool-actions-form and then called back with properties.
 *
 * @example
 * ```js
 * <ToolRandom
 *  @onClear={{action "onClear"}}
 *  @format={{format}}
 *  @bytes={{bytes}}
 *  @random_bytes={{random_bytes}}
 *  @errors={{errors}}/>
 * ```
 * @param onClear {Function} - parent action that is passed through. Must be passed as {{action "onClear"}}
 * @param format {String} - property returned from parent.
 * @param bytes {String} - property returned from parent.
 * @param random_bytes {String} - property returned from parent.
 * @param error=null {Object} - errors passed from parent as default then from child back to parent.
 */

export default class ToolRandom extends Component {
  @action
  onClear() {
    this.args.onClear();
  }
}
