/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
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
