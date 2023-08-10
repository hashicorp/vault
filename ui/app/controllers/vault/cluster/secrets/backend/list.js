/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { or } from '@ember/object/computed';
import { computed } from '@ember/object';
import { inject as service } from '@ember/service';
import Controller from '@ember/controller';
import BackendCrumbMixin from 'vault/mixins/backend-crumb';
import WithNavToNearestAncestor from 'vault/mixins/with-nav-to-nearest-ancestor';
import ListController from 'core/mixins/list-controller';
import { keyIsFolder } from 'core/utils/key-utils';

export default Controller.extend(ListController, BackendCrumbMixin, WithNavToNearestAncestor, {
  flashMessages: service(),
  queryParams: ['page', 'pageFilter', 'tab'],

  tab: '',

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

    delete(item, type) {
      const name = item.id;
      item
        .destroyRecord()
        .then(() => {
          this.flashMessages.success(`${name} was successfully deleted.`);
          this.send('reload');
          if (type === 'secret') {
            this.navToNearestAncestor.perform(name);
          }
        })
        .catch((e) => {
          const error = e.errors ? e.errors.join('. ') : e.message;
          this.flashMessages.danger(error);
        });
    },
  },
});
