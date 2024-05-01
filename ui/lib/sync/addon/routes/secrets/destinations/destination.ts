/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';

import type Store from '@ember-data/store';
import type RouterService from '@ember/routing/router-service';
import type FlashMessageService from 'vault/services/flash-messages';
import type Transition from '@ember/routing/transition';
import type SyncDestinationModel from 'vault/models/sync/destination';

interface RouteParams {
  name: string;
  type: string;
}

export default class SyncSecretsDestinationsDestinationRoute extends Route {
  @service declare readonly store: Store;
  @service declare readonly router: RouterService;
  @service declare readonly flashMessages: FlashMessageService;

  model(params: RouteParams) {
    const { name, type } = params;
    return this.store.findRecord(`sync/destinations/${type}`, name);
  }

  afterModel(model: SyncDestinationModel, transition: Transition) {
    // handles the case where the user attempts to perform actions on a destination when a purge has been initiated
    // editing is available from the list view and syncing secrets is available from the overview
    // the list endpoint does not return the full model so we don't have access to purgeInitiatedAt to disable or hide the actions
    // if transitioning from either of the mentioned routes and a purge has been initiated redirect to the secrets view
    const baseRoute = 'vault.cluster.sync.secrets.destinations.destination';
    const routes = [`${baseRoute}.edit`, `${baseRoute}.sync`];
    const toRoute = transition.to?.name;
    if (toRoute && routes.includes(toRoute) && model.purgeInitiatedAt) {
      const action = transition.to?.localName === 'edit' ? 'Editing a destination' : 'Syncing secrets';
      this.flashMessages.info(`${action} is not permitted once a purge has been initiated.`);
      this.router.replaceWith('vault.cluster.sync.secrets.destinations.destination.secrets');
    }
  }
}
