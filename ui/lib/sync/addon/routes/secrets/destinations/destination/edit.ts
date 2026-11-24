/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import formResolver from 'vault/forms/sync/resolver';

import type { DestinationRouteModel } from '../destination';

export default class SyncSecretsDestinationsDestinationEditRoute extends Route {
  model() {
    const { destination } = this.modelFor('secrets.destinations.destination') as DestinationRouteModel;
    const { type, name, connection_details, options } = destination;
    // granularity is returned as granularityLevel in the response but expected as granularity in the request
    const { granularity_level, ...partialOptions } = options;

    return {
      type,
      form: formResolver(type, {
        name,
        ...connection_details,
        ...partialOptions,
        granularity: granularity_level,
      }),
    };
  }
}
