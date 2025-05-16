/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import formResolver from 'vault/forms/sync/resolver';

import type { DestinationRouteModel } from '../destination';

// originally this route was inheriting the model (Ember Data destination model) from the destination parent route
// an explicit route will be necessary since we will be passing in a Form instance to edit
// this will be done in a follow up PR but for now the Ember Data model will be returned to preserver functionality
export default class SyncSecretsDestinationsDestinationEditRoute extends Route {
  model() {
    const { destination } = this.modelFor('secrets.destinations.destination') as DestinationRouteModel;
    const { type, name, connectionDetails, options } = destination;
    // granularity is returned as granularityLevel in the response but expected as granularity in the request
    const { granularityLevel, ...partialOptions } = options;

    return {
      type,
      form: formResolver(type, {
        name,
        ...connectionDetails,
        ...partialOptions,
        granularity: granularityLevel,
      }),
    };
  }
}
