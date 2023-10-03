/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { inject as service } from '@ember/service';
import Controller from '@ember/controller';
import ListController from 'core/mixins/list-controller';

export default Controller.extend(ListController, {
  flashMessages: service(),

  // callback from HDS pagination to set the queryParams page
  get paginationQueryParams() {
    return (page) => {
      return {
        page,
      };
    };
  },

  actions: {
    delete(model) {
      const type = model.get('identityType');
      const id = model.id;
      return model
        .destroyRecord()
        .then(() => {
          this.send('reload');
          this.flashMessages.success(`Successfully deleted ${type}: ${id}`);
        })
        .catch((e) => {
          this.flashMessages.success(
            `There was a problem deleting ${type}: ${id} - ${e.errors.join(' ') || e.message}`
          );
        });
    },

    toggleDisabled(model) {
      const action = model.get('disabled') ? ['enabled', 'enabling'] : ['disabled', 'disabling'];
      const type = model.get('identityType');
      const id = model.id;
      model.toggleProperty('disabled');

      model
        .save()
        .then(() => {
          this.flashMessages.success(`Successfully ${action[0]} ${type}: ${id}`);
        })
        .catch((e) => {
          this.flashMessages.success(
            `There was a problem ${action[1]} ${type}: ${id} - ${e.errors.join(' ') || e.message}`
          );
        });
    },
    reloadRecord(model) {
      model.reload();
    },
  },
});
