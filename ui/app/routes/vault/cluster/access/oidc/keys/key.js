/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';

export default class OidcKeyRoute extends Route {
  @service api;
  @service capabilities;

  async model({ name }) {
    const { data } = await this.api.identity.oidcReadKey(name);
    const { pathFor } = this.capabilities;
    const paths = {
      key: pathFor('oidcKey', { name }),
      rotate: pathFor('oidcKeyRotate', { name }),
    };
    const capabilities = await this.capabilities.fetch(Object.values(paths));
    return {
      key: { ...data, name },
      capabilities: { ...capabilities[paths.key], canRotate: capabilities[paths.rotate].canUpdate },
    };
  }
}
