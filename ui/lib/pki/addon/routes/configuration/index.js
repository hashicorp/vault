/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import { hash } from 'rsvp';
import { withConfig } from 'pki/decorators/check-config';
@withConfig()
export default class ConfigurationIndexRoute extends Route {
  @service store;

  async fetchMountConfig(backend) {
    const mountConfig = await this.store.query('secret-engine', { path: backend });
    if (mountConfig) {
      return mountConfig.get('firstObject');
    }
  }

  model() {
    const { urls, crl, engine } = this.modelFor('configuration');
    return hash({
      engine,
      hasConfig: this.shouldPromptConfig,
      mountConfig: this.fetchMountConfig(engine.id),
      issuerModel: this.store.createRecord('pki/issuer'),
      urls,
      crl,
    });
  }
}
