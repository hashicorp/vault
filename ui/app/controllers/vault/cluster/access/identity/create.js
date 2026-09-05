/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Controller from '@ember/controller';
import { service } from '@ember/service';
import { action } from '@ember/object';

export default class IdentityCreateController extends Controller {
  @service router;

  showRoute = 'vault.cluster.access.identity.show';
  showTab = 'details';

  @action
  navAfterSave({ saveType, model, id }) {
    const isDelete = saveType === 'delete';
    const type = model.identityType;
    const formType = model.form.identityFormType;

    const listRoutes = {
      'entity-alias': 'vault.cluster.access.identity.aliases.index',
      'group-alias': 'vault.cluster.access.identity.aliases.index',
      group: 'vault.cluster.access.identity.index',
      entity: 'vault.cluster.access.identity.index',
    };

    if (!isDelete) {
      // For aliases, use the aliases.show route instead of the generic show route
      const isAlias = formType === 'alias';
      const showRoute = isAlias ? 'vault.cluster.access.identity.aliases.show' : this.showRoute;
      this.router.transitionTo(showRoute, id, this.showTab);
    } else {
      const routeName = listRoutes[type];
      this.router.transitionTo(routeName);
    }
  }
}
