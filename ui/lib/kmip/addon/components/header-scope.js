/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Component from '@ember/component';
import { inject as service } from '@ember/service';
import layout from '../templates/components/header-scope';

export default Component.extend({
  layout,
  tagName: '',
  secretMountPath: service(),
});
