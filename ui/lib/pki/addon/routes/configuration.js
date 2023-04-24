/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import { hash } from 'rsvp';

export default class PkiConfigurationRoute extends Route {
  @service store;

  async fetchUrls(backend) {
    try {
      return await this.store.findRecord('pki/urls', backend);
    } catch (e) {
      return e.httpStatus;
    }
  }

  async fetchCrl(backend) {
    try {
      return await this.store.findRecord('pki/crl', backend);
    } catch (e) {
      return e.httpStatus;
    }
  }

  model() {
    const engine = this.modelFor('application');
    return hash({
      engine,
      urls: this.fetchUrls(engine.id),
      crl: this.fetchCrl(engine.id),
    });
  }
}
