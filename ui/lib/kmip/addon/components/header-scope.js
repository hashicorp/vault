/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@ember/component';
import { service } from '@ember/service';
import layout from '../templates/components/header-scope';

export default Component.extend({
  layout,
  tagName: '',
  secretMountPath: service(),
});
