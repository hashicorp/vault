/**
 * Copyright IBM Corp. 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { service } from '@ember/service';
import Route from '@ember/routing/route';

export default class ReplicationDetailsRoute extends Route {
  @service replicationMode;
  beforeModel() {
    this.replicationMode.setMode('dr');
  }
}
