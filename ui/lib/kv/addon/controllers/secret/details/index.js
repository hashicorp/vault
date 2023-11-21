/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Controller from '@ember/controller';
import Ember from 'ember';
import { tracked } from '@glimmer/tracking';
import { task, timeout } from 'ember-concurrency';
import { inject as service } from '@ember/service';

export default class KvSecretDetailsIndexController extends Controller {
  @service store; // this.store is referenced in this.fetchSyncStatus()

  @tracked syncDestinations;

  // task is cancelled by resetController() upon leaving the kv v2 details route
  @task
  *pollSyncStatus() {
    while (true) {
      if (Ember.testing) return;

      yield timeout(10000);
      try {
        this.syncDestinations = yield this.fetchSyncStatus(this.model);
      } catch (e) {
        // otherwise keep polling
      }
    }
  }
}
