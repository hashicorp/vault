/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';

/**
 * @module ToolWrap
 * ToolWrap components are components that sys/wrapping/wrap functionality.  Most of the functionality is passed through as actions from the tool-actions-form and then called back with properties.
 *
 * @example
 * <ToolWrap
 *  @onClear={{action "onClear"}}
 *  @token={{token}}
 *  @selectedAction="wrap"
 *  @codemirrorUpdated={{action "codemirrorUpdated"}}
 *  @updateTtl={{action "updateTtl"}}
 *  @buttonDisabled={{buttonDisabled}}
 *  @errors={{errors}}/>
 *
 * @param onClear {Function} - parent action that is passed through. Must be passed as {{action "onClear"}}
 * @param token=null {String} - property passed from parent to child and then passed back up to parent
 * @param selectedAction="wrap" - passed in from parent.  This is the wrap action, others include hash, etc.
 * @param codemirrorUpdated {Function} - parent action that is passed through. Must be passed as {{action "codemirrorUpdated"}}.
 * @param updateTtl {Function} - parent action that is passed through. Must be passed as {{action "updateTtl"}}
 * @param buttonDisabled=false {Boolean} - false default and if there is an error on codemirror it turns to true.
 * @param error=null {Object} - errors passed from parent as default then from child back to parent.
 */

export default class ToolWrap extends Component {
  @tracked data = '{\n}';
  @tracked buttonDisabled = false;

  @action
  onClear() {
    this.args.onClear();
  }
  @action
  updateTtl(evt) {
    if (!evt) return;
    const ttl = evt.enabled ? `${evt.seconds}s` : '30m';
    this.args.updateTtl(ttl);
  }
  @action
  codemirrorUpdated(val, codemirror) {
    codemirror.performLint();
    const hasErrors = codemirror?.state.lint.marked?.length > 0;
    this.data = val;
    this.buttonDisabled = hasErrors;
    this.args.codemirrorUpdated(val, hasErrors);
  }
}
