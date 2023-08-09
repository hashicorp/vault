/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';

export default class OidcKeyRoute extends Route {
  @service store;

  model({ name }) {
    return this.store.findRecord('oidc/key', name);
  }
}
