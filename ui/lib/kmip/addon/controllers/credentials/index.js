/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import ListController from 'core/mixins/list-controller';
import Controller from '@ember/controller';
import { computed } from '@ember/object';
import { getOwner } from '@ember/application';

export default Controller.extend(ListController, {
  mountPoint: computed(function () {
    return getOwner(this).mountPoint;
  }),
});
