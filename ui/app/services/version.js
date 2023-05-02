/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Service, { inject as service } from '@ember/service';
import { keepLatestTask, task } from 'ember-concurrency';
import { tracked } from '@glimmer/tracking';

export default class VersionService extends Service {
  @service store;
  @tracked features = [];
  @tracked version = null;

  get hasPerfReplication() {
    return this.hasFeature('Performance Replication');
  }

  get hasDRReplication() {
    return this.hasFeature('DR Replication');
  }

  get hasSentinel() {
    return this.hasFeature('Sentinel');
  }

  get hasNamespaces() {
    return this.hasFeature('Namespaces');
  }

  get isEnterprise() {
    if (!this.version) return false;
    return this.version.includes('+');
  }

  get isOSS() {
    return !this.isEnterprise;
  }

  hasFeature(feature) {
    if (!this.features) return false;
    return this.features.includes(feature);
  }

  @task
  *getVersion() {
    if (this.version) return;
    const response = yield this.store.adapterFor('cluster').health();
    this.version = response.version;
    return;
  }

  @keepLatestTask
  *getFeatures() {
    if (this.features?.length || this.isOSS) {
      return;
    }
    try {
      const response = yield this.store.adapterFor('cluster').features();
      this.features = response.features;
      return;
    } catch (err) {
      // if we fail here, we're likely in DR Secondary mode and don't need to worry about it
    }
  }

  fetchVersion() {
    return this.getVersion.perform();
  }

  fetchFeatures() {
    return this.getFeatures.perform();
  }
}
