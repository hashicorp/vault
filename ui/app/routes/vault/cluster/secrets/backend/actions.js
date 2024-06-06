/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { parentKeyForKey } from 'core/utils/key-utils';
import EditBase from './secret-edit';

export default EditBase.extend({
  queryParams: {
    selectedAction: {
      replace: true,
    },
  },

  templateName: 'vault/cluster/secrets/backend/transitActionsLayout',

  beforeModel() {
    const { secret } = this.paramsFor(this.routeName);
    const parentKey = parentKeyForKey(secret);
    const { backend } = this.paramsFor('vault.cluster.secrets.backend');
    if (this.backendType(backend) !== 'transit') {
      if (parentKey) {
        return this.router.transitionTo('vault.cluster.secrets.backend.show', parentKey);
      } else {
        return this.router.transitionTo('vault.cluster.secrets.backend.show-root');
      }
    }
  },
  setupController(controller, model) {
    this._super(...arguments);
    const { selectedAction } = this.paramsFor(this.routeName);
    controller.set('selectedAction', selectedAction || model.secret.supportedActions[0]);
    controller.set('breadcrumbs', [
      {
        label: 'Secrets',
        route: 'vault.cluster.secrets',
      },
      {
        label: model.secret.backend,
        route: 'vault.cluster.secrets.backend.list-root',
        model: model.secret.backend,
      },
      {
        label: model.secret.id,
        route: 'vault.cluster.secrets.backend.show',
        models: [model.secret.backend, model.secret.id],
      },
      {
        label: 'actions',
      },
    ]);
  },
});
