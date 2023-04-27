/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Controller, { inject as controller } from '@ember/controller';
import BackendCrumbMixin from 'vault/mixins/backend-crumb';

export default Controller.extend(BackendCrumbMixin, {
  backendController: controller('vault.cluster.secrets.backend'),
  queryParams: ['version', 'itemType'],
  version: '',
  itemType: '',

  reset() {
    this.set('version', '');
    this.set('itemType', '');
  },
  actions: {
    refresh: function () {
      // closure actions don't bubble to routes,
      // so we have to manually bubble here
      this.send('refreshModel');
    },

    toggleAdvancedEdit(bool) {
      this.set('preferAdvancedEdit', bool);
      this.backendController.set('preferAdvancedEdit', bool);
    },
  },
});
