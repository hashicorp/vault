/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';

import type Store from '@ember-data/store';

type RouteParams = {
  name: string;
  type: string;
};

// originally this route was inheriting the model (Ember Data destination model) from the destination parent route
// an explicit route will be necessary since we will be passing in a Form instance to edit
// this will be done in a follow up PR but for now the Ember Data model will be returned to preserver functionality
export default class SyncSecretsDestinationsDestinationEditRoute extends Route {
  @service declare readonly store: Store;

  model() {
    const { name, type } = this.paramsFor('secrets.destinations.destination') as RouteParams;
    return this.store.findRecord(`sync/destinations/${type}`, name);
  }
}
