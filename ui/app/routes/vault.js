/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { later } from '@ember/runloop';
import { Promise } from 'rsvp';
import { service } from '@ember/service';
import Route from '@ember/routing/route';
import Ember from 'ember';
import ClusterResource from 'vault/resources/cluster';

const SPLASH_DELAY = 300;

export default class VaultRoute extends Route {
  @service router;
  @service version;

  cluster = null;

  beforeModel() {
    // So we can know what type (Enterprise/Community) we're running
    return this.version.fetchType();
  }

  model() {
    const delay = Ember.testing ? 0 : SPLASH_DELAY;
    // hardcode single cluster
    this.cluster ??= new ClusterResource({ id: '1', name: 'vault' }, this);
    return new Promise((resolve) => {
      later(() => {
        resolve([this.cluster]);
      }, delay);
    });
  }

  redirect(model, transition) {
    if (model.length === 1 && transition.targetName === 'vault.index') {
      return this.router.transitionTo('vault.cluster', model[0].name);
    }
  }
}
