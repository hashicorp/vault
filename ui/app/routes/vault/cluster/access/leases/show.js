/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { hash } from 'rsvp';
import { action } from '@ember/object';
import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { keyIsFolder, parentKeyForKey } from 'core/utils/key-utils';

export default class LeasesShowRoute extends Route {
  @service api;
  @service capabilities;
  @service router;

  beforeModel() {
    const { lease_id: leaseId } = this.paramsFor(this.routeName);
    const parentKey = parentKeyForKey(leaseId);
    if (keyIsFolder(leaseId)) {
      if (parentKey) {
        return this.router.transitionTo('vault.cluster.access.leases.list', parentKey);
      } else {
        return this.router.transitionTo('vault.cluster.access.leases.list-root');
      }
    }
  }

  async model(params) {
    const { lease_id } = params;
    const leaseData = await this.api.sys.leasesReadLease({ lease_id });
    // isAuthLease is derived from the lease id prefix (previously a computed on the model)
    leaseData.isAuthLease = /^auth/.test(lease_id);

    return hash({
      lease: leaseData,
      capabilities: hash({
        renew: this.capabilities.fetchPathCapabilities('sys/leases/renew'),
        revoke: this.capabilities.fetchPathCapabilities('sys/leases/revoke'),
        leases: this.modelFor('vault.cluster.access.leases'),
      }),
    });
  }

  setupController(controller, model) {
    super.setupController(...arguments);
    const { lease_id: leaseId } = this.paramsFor(this.routeName);
    controller.setProperties({
      model: model.lease,
      capabilities: model.capabilities,
      baseKey: { id: leaseId },
    });
  }

  @action
  error(error) {
    const { lease_id } = this.paramsFor(this.routeName);
    error.keyId = lease_id;
    return true;
  }

  @action
  refreshModel() {
    this.refresh();
  }
}
