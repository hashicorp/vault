/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';

export default class OidcProviderRoute extends Route {
  @service store;

  model({ name }) {
    return this.store.findRecord('oidc/provider', name);
  }
}
