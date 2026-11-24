/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { hash } from 'rsvp';
import Route from '@ember/routing/route';
import { service } from '@ember/service';

export default class SealRoute extends Route {
  @service capabilities;

  model() {
    return hash({
      cluster: this.modelFor('vault.cluster'),
      seal: this.capabilities.fetchPathCapabilities('sys/seal'),
    });
  }
}
