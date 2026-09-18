/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';

export default class StorageRoute extends Route {
  @service api;

  // raft/configuration isn't in the OpenAPI spec/generated client, so use the raw request interface
  async model() {
    const response = await this.api.request.get('/sys/storage/raft/configuration');
    const { data } = await response.json();
    return data?.config?.servers ?? [];
  }
}
