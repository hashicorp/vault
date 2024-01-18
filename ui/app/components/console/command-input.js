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
      switch (keyCode) {
        case keys.ENTER:
          // pass only the target so you have access to the element for repositioning the cursor
          this.onExecuteCommand.perform(event.target);
          break;
        case keys.UP:
        case keys.DOWN:
          this.onShiftCommand(keyCode);
          break;
        default:
          this.onValueUpdate(event.target.value);
      }
    },
    fullscreen() {
      this.onFullscreen();
    },
  },
});
