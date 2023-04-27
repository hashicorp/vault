/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { debounce } from '@ember/runloop';

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
