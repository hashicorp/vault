/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Service from '@ember/service';

export default Service.extend({
  mode: null,

  getMode() {
    return this.mode;
  },

  setMode(mode) {
    this.set('mode', mode);
  },
});
