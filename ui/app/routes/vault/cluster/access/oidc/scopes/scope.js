/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';

export default class OidcScopeRoute extends Route {
  @service api;
  @service capabilities;

  async model({ name }) {
    const { data } = await this.api.identity.oidcReadScope(name);
    const capabilities = await this.capabilities.for('oidcScope', { name });
    return {
      scope: { ...data, name },
      capabilities,
    };
  }
}
