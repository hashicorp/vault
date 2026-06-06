/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Controller from '@ember/controller';
import { computed } from '@ember/object';
import { or } from '@ember/object/computed';
import { service } from '@ember/service';
import { SecretsApiSshListRolesListEnum } from '@hashicorp/vault-client-typescript';
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

    // Adds or removes the given SSH role from the zero-address config, then reloads the list.
    async toggleZeroAddress(item) {
      const backendPath = item.backend;
      this.set('loading-' + item.id, true);
      try {
        const response = await this.api.secrets.sshListRoles(
          backendPath,
          SecretsApiSshListRolesListEnum.TRUE
        );
        const allRoles = this.api.keyInfoToArray(response);
        const newValue = !item.zero_address;
        const zeroAddressRoles = allRoles
          .filter((role) => (role.id === item.id ? newValue : role.zero_address))
          .map((role) => role.id);
        if (zeroAddressRoles.length === 0) {
          await this.api.secrets.sshDeleteZeroAddressConfiguration(backendPath);
        } else {
          await this.api.secrets.sshConfigureZeroAddress(backendPath, { roles: zeroAddressRoles });
        }
        this.send('reload');
      } catch (e) {
        const { message } = await this.api.parseError(e);
        this.flashMessages.danger(message);
      } finally {
        this.set('loading-' + item.id, false);
      }
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
      } else if (this.backendType === 'totp') {
        try {
          await this.api.secrets.totpDeleteKey(name, item.backend);
          this.flashMessages.success(`${name} was successfully deleted.`);
          this.send('reload');
        } catch (e) {
          const { message } = await this.api.parseError(e);
          this.flashMessages.danger(message);
        }
      } else if (this.backendType === 'ssh') {
        try {
          await this.api.secrets.sshDeleteRole(name, item.backend);
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
