/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { clickable, isPresent } from 'ember-cli-page-object';

export default {
  textareaIsPresent: isPresent('[data-test-textarea]'),
  copyButtonIsPresent: isPresent('[data-test-copy-button]'),
  toggleMasked: clickable('[data-test-button="toggle-masked"]'),
};
