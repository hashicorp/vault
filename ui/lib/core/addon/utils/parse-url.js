/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

// adapted from https://gist.github.com/jed/964849
const fn = (function (anchor) {
  return function (url) {
    anchor.href = url;
    const parts = {};
    for (const prop in anchor) {
      if ('' + anchor[prop] === anchor[prop]) {
        parts[prop] = anchor[prop];
      }
    }

    return parts;
  };
})(document.createElement('a'));

export default fn;
