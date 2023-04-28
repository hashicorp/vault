/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { inject as service } from '@ember/service';
import Component from '@ember/component';

export default Component.extend({
  tagName: '',
  version: service(),
});
