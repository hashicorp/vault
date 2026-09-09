/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';

export default class OidcClientRoute extends Route {
  @service api;
  @service capabilities;

  async model({ name }) {
    const { data } = await this.api.identity.oidcReadClient(name);
    const capabilities = await this.capabilities.for('oidcClient', { name });
    return { client: { ...data, name }, capabilities };
  }
}
