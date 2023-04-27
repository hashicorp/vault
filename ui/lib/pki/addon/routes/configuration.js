/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import { hash } from 'rsvp';

export default class PkiConfigurationRoute extends Route {
  @service store;

  model() {
    const engine = this.modelFor('application');
    return hash({
      engine,
      urls: this.store.findRecord('pki/urls', engine.id).catch((e) => e.httpStatus),
      crl: this.store.findRecord('pki/crl', engine.id).catch((e) => e.httpStatus),
    });
  }
}
