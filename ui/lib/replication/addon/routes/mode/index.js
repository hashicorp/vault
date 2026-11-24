/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { setProperties } from '@ember/object';
import Route from '@ember/routing/route';
import { service } from '@ember/service';

export default Route.extend({
  router: service('app-router'),
  capabilities: service(),
  async model() {
    const replicationMode = this.paramsFor('mode').replication_mode;
    const cluster = this.modelFor('mode');
    const tokenPath = `sys/replication/${replicationMode}/primary/secondary-token`;
    const revokePath = `sys/replication/${replicationMode}/primary/revoke-secondary`;
    const perms = await this.capabilities.fetch([tokenPath, revokePath]);

    setProperties(cluster, {
      canRevokeSecondary: perms[revokePath].canUpdate,
      canAddSecondary: perms[tokenPath].canUpdate,
    });
    return cluster;
  },
});
