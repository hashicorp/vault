/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { service } from '@ember/service';
import Component from '@ember/component';

export default Component.extend({
  flashMessages: service(),

  actions: {
    enable(model) {
      model.set('disabled', false);

      model
        .save()
        .then(() => {
          this.flashMessages.success(`Successfully enabled entity: ${model.id}`);
        })
        .catch((e) => {
          this.flashMessages.success(
            `There was a problem enabling the entity: ${model.id} - ${e.errors.join(' ') || e.message}`
          );
        });
    },
  },
});
