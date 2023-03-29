/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';

/**
 * @module ToolHash
 * ToolHash components are components that sys/wrapping/hash functionality.  Most of the functionality is passed through as actions from the tool-actions-form and then called back with properties.
 *
 * @example
 * ```js
 * <ToolHash
 *  @onClear={{action "onClear"}}
 *  @sum={{sum}}
 *  @algorithm={{algorithm}}
 *  @format={{format}}
 *  @errors={{errors}}/>
 * ```
 * @param onClear {Function} - parent action that is passed through. Must be passed as {{action "onClear"}}
 * @param sum=null {String} - property passed from parent to child and then passed back up to parent.
 * @param algorithm {String} - property returned from parent.
 * @param format {String} - property returned from parent.
 * @param error=null {Object} - errors passed from parent as default then from child back to parent.
 */
export default class ToolHash extends Component {
  @action
  onClear() {
    this.args.onClear();
  }
}
