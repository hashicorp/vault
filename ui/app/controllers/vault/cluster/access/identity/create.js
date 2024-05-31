/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Controller from '@ember/controller';
import { service } from '@ember/service';

export default Controller.extend({
  router: service(),
  showRoute: 'vault.cluster.access.identity.show',
  showTab: 'details',

  actions: {
    navAfterSave({ saveType, model }) {
      const isDelete = saveType === 'delete';
      const type = model.identityType;
      const listRoutes = {
        'entity-alias': 'vault.cluster.access.identity.aliases.index',
        'group-alias': 'vault.cluster.access.identity.aliases.index',
        group: 'vault.cluster.access.identity.index',
        entity: 'vault.cluster.access.identity.index',
      };
      if (!isDelete) {
        this.router.transitionTo(this.showRoute, model.id, this.showTab);
      } else {
        const routeName = listRoutes[type];
        this.router.transitionTo(routeName);
      }
    },
  },
});
