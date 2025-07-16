/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { clickable } from 'ember-cli-page-object';

export default {
  enable: clickable('[data-test-enable-identity]'),
};
