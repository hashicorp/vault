/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { service } from '@ember/service';
import Controller from '@ember/controller';
export default Controller.extend({
  namespaceService: service('namespace'),
  router: service(),
  actions: {
    onSave({ saveType }) {
      if (saveType === 'save') {
        // fetch new namespaces for the namespace picker
        this.namespaceService.findNamespacesForUser.perform();
        return this.router.transitionTo('vault.cluster.access.namespaces.index');
      }
    },
  },
});
