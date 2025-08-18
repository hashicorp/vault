/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@ember/component';
import { service } from '@ember/service';
import { or } from '@ember/object/computed';

export default Component.extend({
  tagName: '',
  secretMountPath: service(),
  shouldShowPath: or('showPath', 'scope', 'role'),
  showPath: false,
  path: null,
  scope: null,
  role: null,
});
