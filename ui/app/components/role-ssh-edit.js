/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import RoleEdit from './role-edit';
import { computed } from '@ember/object';

export default RoleEdit.extend({
  init() {
    this._super(...arguments);
    this.set('backendType', 'ssh');
  },

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

  subtitle: computed('mode', 'model.id', function () {
    if (this.mode === 'create' || this.mode === 'edit') return;

    return this.model.id;
  }),

  actions: {
    updateTtl(path, val) {
      const model = this.model;
      const valueToSet = val.enabled === true ? `${val.seconds}s` : undefined;
      model.set(path, valueToSet);
    },
  },
});
