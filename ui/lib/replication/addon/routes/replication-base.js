/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { service } from '@ember/service';
import Route from '@ember/routing/route';

export default class ReplicationBaseRoute extends Route {
  @service('app-router') router;
  @service api;
  @service version;
  @service('replication-mode') rm;

  get replicationMode() {
    return this.rm.mode;
  }
}
