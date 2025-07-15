/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { findDestination } from 'core/helpers/sync-destinations';
import formResolver from 'vault/forms/sync/resolver';

import type { DestinationType } from 'vault/sync';

type Params = {
  type: DestinationType;
};

export default class SyncSecretsDestinationsCreateDestinationRoute extends Route {
  model(params: Params) {
    const { type } = params;
    const { defaultValues } = findDestination(type);
    return {
      type,
      form: formResolver(type, defaultValues, { isNew: true }),
    };
  }
}
