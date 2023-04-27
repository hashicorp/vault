/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Component from '@ember/component';

export default Component.extend({
  'data-test-component': 'console/output-log',
  attributeBindings: ['data-test-component'],
  log: null,
});
