/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import { withConfig } from 'pki/decorators/check-issuers';
import { hash } from 'rsvp';

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
    const { acme, cluster, urls, crl, engine } = this.modelFor('configuration');
    return hash({
      hasConfig: this.shouldPromptConfig,
      engine,
      acme,
      cluster,
      urls,
      crl,
      mountConfig: this.fetchMountConfig(engine.id),
      issuerModel: this.store.createRecord('pki/issuer'),
    });
  }
}
