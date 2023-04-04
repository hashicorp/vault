/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';

/**
 * @module ToolLookup
 * ToolLookup components are components that sys/wrapping/lookup functionality.  Most of the functionality is passed through as actions from the tool-actions-form and then called back with properties.
 *
 * @example
 * ```js
 * <ToolLookup
 *  @creation_time={{creation_time}}
 *  @creation_ttl={{creation_ttl}}
 *  @creation_path={{creation_path}}
 *  @expirationDate={{expirationDate}}
 *  @selectedAction={{selectedAction}}
 *  @token={{token}}
 *  @onClear={{action "onClear"}}
 *  @errors={{errors}}/>
 * ```
 * @param creation_time {Function} - parent action that is passed through.
 * @param creation_ttl {Function} - parent action that is passed through.
 * @param creation_path {Function} - parent action that is passed through.
 * @param expirationDate='' {String} - value returned from lookup.
 * @param selectedAction="wrap" - passed in from parent.  This is the wrap action, others include hash, etc.
 * @param token=null {String} - property passed from parent to child and then passed back up to parent
 * @param onClear {Function} - parent action that is passed through. Must be passed as {{action "onClear"}}
 * @param error=null {Object} - errors passed from parent as default then from child back to parent.
 */
export default class ToolLookup extends Component {
  @action
  onClear() {
    this.args.onClear();
  }
}
