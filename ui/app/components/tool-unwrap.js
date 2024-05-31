/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';

/**
 * @module ToolUnwrap
 * ToolUnwrap components are components that sys/wrapping/unwrap functionality.  Most of the functionality is passed through as actions from the tool-actions-form and then called back with properties.
 *
 * @example
 * ```js
 * <ToolUnwrap
 *  @onClear={{action "onClear"}}
 *  @token={{token}}
 *  @unwrap_data={{unwrap_data}}
 *  @details={{details}}
 *  @errors={{errors}}/>
 * ```
 * @param onClear {Function} - parent action that is passed through. Must be passed as {{action "onClear"}}
 * @param token=null {String} - property passed from parent to child and then passed back up to parent
 * @param unwrap_data {String} - property returned from parent.
 * @param details {String} - property returned from parent.
 * @param error=null {Object} - errors passed from parent as default then from child back to parent.
 */

export default class ToolUnwrap extends Component {
  @action
  onClear() {
    this.args.onClear();
  }
}
