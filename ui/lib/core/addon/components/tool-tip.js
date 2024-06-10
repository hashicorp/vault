/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { debounce } from '@ember/runloop';

/**
 * @deprecated
 * @module ToolTip
 *
 * @example
 * <ToolTip @verticalPosition="above" as |T|>
 *  <T.Trigger>
 *    <Icon @name="check-circle" class="has-text-success" />
 *  </T.Trigger>
 *  <T.Content @defaultClass="tool-tip">
 *    <div class="box">
 *      My tooltip text
 *    </div>
 *  </T.Content>
 * </ToolTip>
 * 
 *  * Use HDS tooltip instead

 * @param {string} [verticalPosition] - vertical position specification (above, below)
 * @param {string} [horizontalPosition] - horizontal position specification (center, auto-right)
 *
 */
export default class ToolTipComponent extends Component {
  get delay() {
    return this.args.delay || 200;
  }
  get horizontalPosition() {
    return this.args.horizontalPosition || 'auto-right';
  }

  toggleState({ dropdown, action }) {
    dropdown.actions[action]();
  }

  @action
  open(dropdown) {
    debounce(this, 'toggleState', { dropdown, action: 'open' }, this.delay);
  }
  @action
  close(dropdown) {
    debounce(this, 'toggleState', { dropdown, action: 'close' }, this.delay);
  }
  @action
  prevent() {
    return false;
  }
}
