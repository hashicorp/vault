/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';

export default class OidcClientsCreateRoute extends Route {
  @service store;

  model() {
    return this.store.createRecord('oidc/client');
  }
}
