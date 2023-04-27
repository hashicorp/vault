/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Controller from '@ember/controller';
import ListController from 'core/mixins/list-controller';

export default Controller.extend(ListController, {
  actions: {
    onDelete() {
      this.send('reload');
    },
  },
});
