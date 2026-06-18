/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@ember/component';
import { service } from '@ember/service';
import { computed } from '@ember/object';
import { isBlank } from '@ember/utils';

const LIST_ROOT_ROUTE = 'vault.cluster.secrets.backend.list-root';
const SHOW_ROUTE = 'vault.cluster.secrets.backend.show';

export default Component.extend({
  api: service(),
  router: service(),
  flashMessages: service(),

  mode: null,
  form: null,

  breadcrumbs: computed('root', 'title', function () {
    return [
      { label: 'Vault', text: 'Vault', icon: 'vault', path: 'vault.cluster.dashboard' },
      { text: 'Secrets engines', path: 'vault.cluster.secrets.backends' },
      this.root,
      { label: this.title, text: this.title },
    ];
  }),

  title: computed('mode', function () {
    if (this.mode === 'create') {
      return 'Create SSH Role';
    } else if (this.mode === 'edit') {
      return 'Edit SSH Role';
    } else {
      return 'SSH Role';
    }
  }),

  subtitle: computed('mode', 'form.name', function () {
    if (this.mode === 'create' || this.mode === 'edit') return;
    return this.form.name;
  }),

  displayFields: computed('form.{data.key_type,displayFields}', function () {
    return this.form?.displayFields ?? [];
  }),

  actions: {
    async createOrUpdate(type, event) {
      event.preventDefault();

      const { form } = this;
      const { name, backend, ...roleData } = form.data;

      if (type === 'create' && isBlank(name)) {
        this.flashMessages.danger('Role name is required');
        return;
      }

      try {
        await this.api.secrets.sshWriteRole(name, backend, roleData);
        this.router.transitionTo(SHOW_ROUTE, name);
      } catch (error) {
        const { message } = await this.api.parseError(error);
        this.flashMessages.danger(message);
      }
    },

    async delete() {
      const { form } = this;
      const { name, backend } = form.data;

      try {
        await this.api.secrets.sshDeleteRole(name, backend);
        this.router.transitionTo(LIST_ROOT_ROUTE);
      } catch (error) {
        const { message } = await this.api.parseError(error);
        this.flashMessages.danger(message);
      }
    },

    updateTtl(path, val) {
      this.form[path] = val.enabled === true ? `${val.seconds}s` : undefined;
    },
  },
});
