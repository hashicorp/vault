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
 *  @errors={{@errors}}
 *  @onBack={{action "onBack" (array "token")}}
 *  @onChange={{action "onChange"}}
 *  @onClear={{action "onClear"}}
 *  @token={{@token}}
 * />
 *
 * @param {object} errors=null - errors returned if wrap fails
 * @param {function} onBack - callback that only clears specific values so the action can be repeated. Must be passed as `{{action "onBack"}}`
 * @param {function} onChange - callback that fires when inputs change and passes value and param name back to the parent
 * @param {function} onClear - callback that resets all of values to defaults. Must be passed as `{{action "onClear"}}`
 * @param {string} token=null - returned after user clicks "Wrap data", if there is a token value it displays instead of the JsonEditor
 */

export default class ToolWrap extends Component {
  @tracked buttonDisabled = false;

  @action
  updateTtl(evt) {
    if (!evt) return;
    const ttl = evt.enabled ? `${evt.seconds}s` : '30m';
    this.args.onChange('wrapTTL', ttl);
  }

  @action
  codemirrorUpdated(val, codemirror) {
    codemirror.performLint();
    const hasErrors = codemirror?.state.lint.marked?.length > 0;
    this.buttonDisabled = hasErrors;
    this.args.onChange('data', val);
  }
}
