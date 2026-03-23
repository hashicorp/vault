/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Controller from '@ember/controller';
import { computed } from '@ember/object';
import { or } from '@ember/object/computed';
import { service } from '@ember/service';
import ListController from 'core/mixins/list-controller';
import { keyIsFolder } from 'core/utils/key-utils';
import BackendCrumbMixin from 'vault/mixins/backend-crumb';

export default Controller.extend(ListController, BackendCrumbMixin, {
  flashMessages: service(),
  api: service(),
  router: service(),
  queryParams: ['page', 'pageFilter', 'tab'],

  tab: '',

  // callback from HDS pagination to set the queryParams page
  get paginationQueryParams() {
    return (page) => {
      return {
        page,
      };
    };
  },

  filterIsFolder: computed('filter', function () {
    return !!keyIsFolder(this.filter);
  }),

  isConfigurableTab: or('isCertTab', 'isConfigure'),

  actions: {
    chooseAction(action) {
      this.set('selectedAction', action);
    },

    toggleZeroAddress(item, backend) {
      item.toggleProperty('zeroAddress');
      this.set('loading-' + item.id, true);
      backend
        .saveZeroAddressConfig()
        .catch((e) => {
          item.set('zeroAddress', false);
          this.flashMessages.danger(e.message);
        })
        .finally(() => {
          this.set('loading-' + item.id, false);
        });
    },

    delete(item) {
      const name = item.id;
      item
        .destroyRecord()
        .then(() => {
          this.flashMessages.success(`${name} was successfully deleted.`);
          this.send('reload');
        })
        .catch((e) => {
          const error = e.errors ? e.errors.join('. ') : e.message;
          this.flashMessages.danger(error);
        });
    },
  },
});
