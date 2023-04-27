/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { reads } from '@ember/object/computed';
import Component from '@ember/component';

export default Component.extend({
  content: null,
  list: reads('content.keys'),
});
