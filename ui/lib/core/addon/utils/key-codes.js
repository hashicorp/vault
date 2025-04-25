/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

/*
 * DEPRCATED (see: https://developer.mozilla.org/en-US/docs/Web/API/KeyboardEvent/keyCode).
 *
 * TODO: Replace all instances with `event.key` (use lib/core/addon/utils/keyboard-keys.ts).
 * `event.keyCode` is deprecated and will be removed in future versions of browsers.
 */

// a map of keyCode for use in keyboard event handlers
export default {
  ENTER: 13,
  ESC: 27,
  TAB: 9,
  LEFT: 37,
  UP: 38,
  RIGHT: 39,
  DOWN: 40,
  T: 116,
  BACKSPACE: 8,
};
