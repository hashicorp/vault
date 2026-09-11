/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Controller from '@ember/controller';
import { service } from '@ember/service';
import { action } from '@ember/object';

export default class IdentityCreateController extends Controller {
  @service router;

  showRoute = 'vault.cluster.access.identity.entities.show';
  showTab = 'details';

  @action
  navAfterSave({ saveType, model, id }) {
    const isDelete = saveType === 'delete';
    const type = model.identityType;
    const formType = model.form.identityFormType;

    const listRoutes = {
      'entity-alias': 'vault.cluster.access.identity.entities.aliases.index',
      'group-alias': 'vault.cluster.access.identity.entities.aliases.index',
      group: 'vault.cluster.access.identity.entities.index',
      entity: 'vault.cluster.access.identity.entities.index',
    };

    if (!isDelete) {
      // For aliases, use the aliases.show route instead of the generic show route
      const isAlias = formType === 'alias';
      const showRoute = isAlias ? 'vault.cluster.access.identity.entities.aliases.show' : this.showRoute;
      this.router.transitionTo(showRoute, id, this.showTab);
    } else {
      const routeName = listRoutes[type];
      this.router.transitionTo(routeName);
    }
  }
}
