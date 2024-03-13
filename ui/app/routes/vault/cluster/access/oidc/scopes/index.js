/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';

export default class OidcScopesRoute extends Route {
  @service store;

  model() {
    return this.store.query('oidc/scope', {}).catch((err) => {
      if (err.httpStatus === 404) {
        return [];
      } else {
        throw err;
      }
    });
  }
}
