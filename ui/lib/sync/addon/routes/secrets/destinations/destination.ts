/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import apiMethodResolver from 'sync/utils/api-method-resolver';

import type RouterService from '@ember/routing/router-service';
import type FlashMessageService from 'vault/services/flash-messages';
import type Transition from '@ember/routing/transition';
import type ApiService from 'vault/services/api';
import type { Destination } from 'vault/sync';
import type CapabilitiesService from 'vault/services/capabilities';
import type { CapabilitiesMap } from 'vault/app-types';

type Params = {
  name: string;
  type: Destination['type'];
};

export type DestinationRouteModel = {
  destination: Destination;
  capabilities: CapabilitiesMap;
};

export default class SyncSecretsDestinationsDestinationRoute extends Route {
  @service declare readonly api: ApiService;
  @service('app-router') declare readonly router: RouterService;
  @service declare readonly flashMessages: FlashMessageService;
  @service declare readonly capabilities: CapabilitiesService;

  async model(params: Params) {
    const { name, type } = params;
    const method = apiMethodResolver('read', type);

    const requests = [
      this.api.sys[method](name, {}),
      this.capabilities.fetch([
        this.capabilities.pathFor('syncDestination', { name, type }),
        this.capabilities.pathFor('syncSetAssociation', { name, type }),
      ]),
    ];
    const [destination, capabilities] = await Promise.all(requests);

    return {
      destination,
      capabilities,
    };
  }

  afterModel({ destination }: DestinationRouteModel, transition: Transition) {
    // handles the case where the user attempts to perform actions on a destination when a purge has been initiated
    // editing is available from the list view and syncing secrets is available from the overview
    // the list endpoint does not return the full model so we don't have access to purgeInitiatedAt to disable or hide the actions
    // if transitioning from either of the mentioned routes and a purge has been initiated redirect to the secrets view
    const baseRoute = 'vault.cluster.sync.secrets.destinations.destination';
    const routes = [`${baseRoute}.edit`, `${baseRoute}.sync`];
    const toRoute = transition.to?.name;
    if (toRoute && routes.includes(toRoute) && destination.purgeInitiatedAt) {
      const action = transition.to?.localName === 'edit' ? 'Editing a destination' : 'Syncing secrets';
      this.flashMessages.info(`${action} is not permitted once a purge has been initiated.`);
      this.router.replaceWith('vault.cluster.sync.secrets.destinations.destination.secrets');
    }
  }
}
