/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Component from '@ember/component';

export default Component.extend({
  onSubmit() {},
  replicationMode: null,
  replicationModeForDisplay: null,
  model: null,

  actions: {
    onSubmit() {
      return this.onSubmit(...arguments);
    },
  },
});
