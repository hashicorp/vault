/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

// a map of keyCode for use in keyboard event handlers
export default {
  ENTER: 'Enter',
  // some older browsers use 'Esc' instead of 'Escape'
  ESC: ['Escape', 'Esc'],
  TAB: 'Tab',
  LEFT: 'ArrowLeft',
  UP: 'ArrowUp',
  RIGHT: 'ArrowRight',
  DOWN: 'ArrowDown',
  T: 'T',
  BACKSPACE: 'Backspace',
  SPACE: ' ',
};
