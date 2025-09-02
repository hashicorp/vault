/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { service } from '@ember/service';
import Helper from '@ember/component/helper';

export default Helper.extend({
  flashMessages: service(),

  compute([message, type]) {
    return () => {
      this.flashMessages[type || 'success'](message);
    };
  },
});
