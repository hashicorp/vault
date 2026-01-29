/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';

export default class StorageRoute extends Route {
  @service store;

  model() {
    // findAll method will return all records in store as well as response from server
    // when removing a peer via the cli, stale records would continue to appear until refresh
    // query method will only return records from response
    return this.store.query('server', {});
  }
}
