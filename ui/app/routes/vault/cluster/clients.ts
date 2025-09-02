/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { hash } from 'rsvp';
import { service } from '@ember/service';

import type Store from '@ember-data/store';

export default class ClientsRoute extends Route {
  @service declare readonly store: Store;

  getVersionHistory(): Promise<
    Array<{ version: string; previousVersion: string; timestampInstalled: string }>
  > {
    return this.store
      .findAll('clients/version-history')
      .then((response) => {
        return response.map(({ version, previousVersion, timestampInstalled }) => {
          return {
            version,
            previousVersion,
            timestampInstalled,
          };
        });
      })
      .catch(() => []);
  }

  model() {
    // swallow config error so activity can show if no config permissions
    return hash({
      config: this.store.queryRecord('clients/config', {}).catch(() => ({})),
      versionHistory: this.getVersionHistory(),
    });
  }
}
