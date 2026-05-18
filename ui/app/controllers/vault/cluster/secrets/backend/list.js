/**
 * Copyright IBM Corp. 2026
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

    async delete(item) {
      const name = item.id;
      // Handle keymgmt list items (plain objects from API service)
      if (this.backendType === 'keymgmt' && item.type === 'key') {
        try {
          await this.api.secrets.keyManagementDeleteKey(name, item.backend);
          this.flashMessages.success(`${name} was successfully deleted.`);
          this.send('reload');
        } catch (e) {
          const { message } = await this.api.parseError(e);
          this.flashMessages.danger(message);
        }
      } else if (this.backendType === 'keymgmt' && item.type === 'provider') {
        try {
          await this.api.secrets.keyManagementDeleteKmsProvider(name, item.backend);
          this.flashMessages.success(`${name} was successfully deleted.`);
          this.send('reload');
        } catch (e) {
          const { message } = await this.api.parseError(e);
          this.flashMessages.danger(message);
        }
      } else {
        // Handle Ember Data models
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
      }
    },
  },
});
