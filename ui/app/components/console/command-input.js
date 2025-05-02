/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@ember/component';

export default Component.extend({
  onExecuteCommand() {},
  onFullscreen() {},
  onValueUpdate() {},
  onShiftCommand() {},
  value: null,
  isFullscreen: null,

  actions: {
    handleKeyUp(event) {
      const val = event.target.value;
      switch (event.key) {
        case 'Enter':
          this.onExecuteCommand(val);
          break;
        case 'ArrowUp':
        case 'ArrowDown':
          this.onShiftCommand(event.key);
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
