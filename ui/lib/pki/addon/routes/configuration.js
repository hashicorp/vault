/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { hash } from 'rsvp';

export default class PkiConfigurationRoute extends Route {
  @service store;

  model() {
    const engine = this.modelFor('application');
    return hash({
      engine,
      acme: this.store.findRecord('pki/config/acme', engine.id).catch((e) => e.httpStatus),
      cluster: this.store.findRecord('pki/config/cluster', engine.id).catch((e) => e.httpStatus),
      urls: this.store.findRecord('pki/config/urls', engine.id).catch((e) => e.httpStatus),
      crl: this.store.findRecord('pki/config/crl', engine.id).catch((e) => e.httpStatus),
    });
  }
}
