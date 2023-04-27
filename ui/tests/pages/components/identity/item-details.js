/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { clickable } from 'ember-cli-page-object';

export default {
  enable: clickable('[data-test-enable]'),
};
