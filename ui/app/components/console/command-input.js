/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@ember/component';
import keys from 'core/utils/key-codes';

export default Component.extend({
  onExecuteCommand() {},
  onFullscreen() {},
  onValueUpdate() {},
  onShiftCommand() {},
  value: null,
  isFullscreen: null,

  actions: {
    handleKeyUp(event) {
      const keyCode = event.keyCode;
      const val = event.target.value;
      switch (keyCode) {
        case keys.ENTER:
          this.onExecuteCommand(val);
          break;
        case keys.UP:
        case keys.DOWN:
          this.onShiftCommand(keyCode);
          break;
        default:
          this.onValueUpdate(val);
      }
    },
    fullscreen() {
      this.onFullscreen();
    },
  },
});
