/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { helper } from '@ember/component/helper';
import { stringify } from 'core/helpers/stringify';

export function stringifiedSecretData(secretData) {
  // must pass the third param called "space" in JSON.stringify to structure object with whitespace
  // otherwise the following codemirror modifier check will pass `this._editor.getValue() !== namedArgs.content` and _setValue will be called.
  // the method _setValue moves the cursor to the beginning of the text field.
  // the effect is that the cursor jumps after the first key input.
  const startingValue = JSON.stringify({ '': '' }, null, 2);
  return secretData ? stringify([secretData], {}) : startingValue;
}

export default helper(([secretData]) => stringifiedSecretData(secretData));
