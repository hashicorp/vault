/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import { hash } from 'rsvp';
import { withConfig } from 'pki/decorators/check-config';

@withConfig()
export default class PkiConfigurationRoute extends Route {
  @service store;

  async fetchMountConfig(backend) {
    const mountConfig = await this.store.query('secret-engine', { path: backend });

    if (mountConfig) {
      return mountConfig.get('firstObject');
    }
  }

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
      hasConfig: this.shouldPromptConfig,
      engine,
      urls: this.fetchUrls(engine.id),
      crl: this.fetchCrl(engine.id),
      mountConfig: this.fetchMountConfig(engine.id),
      issuerModel: this.store.createRecord('pki/issuer'),
    });
  }
}
