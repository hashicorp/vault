/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import OidcScopeForm from 'vault/forms/oidc/scope';

export default class OidcScopesCreateRoute extends Route {
  model() {
    return new OidcScopeForm({}, { isNew: true });
  }
}
