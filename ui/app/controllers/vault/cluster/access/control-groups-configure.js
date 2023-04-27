/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Controller from '@ember/controller';

export default Controller.extend({
  actions: {
    onSave({ saveType }) {
      if (saveType === 'destroyRecord') {
        this.send('reload');
      }
    },
  },
});
