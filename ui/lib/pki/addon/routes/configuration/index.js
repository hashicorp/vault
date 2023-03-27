/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import { withConfig } from 'pki/decorators/check-config';
import { hash } from 'rsvp';

@withConfig()
export default class ConfigurationIndexRoute extends Route {
  @service store;
  @service secretMountPath;

  model() {
    const backend = this.secretMountPath.currentPath;
    return hash({
      hasConfig: this.shouldPromptConfig,
      engine: this.modelFor('application'),
      urls: this.store.findRecord('pki/urls', backend),
      crl: this.store.findRecord('pki/crl', backend),
      mountConfig: this.fetchMountConfig(backend),
    });
  }

  async fetchMountConfig(path) {
    const mountConfig = await this.store.query('secret-engine', { path });

    if (mountConfig) {
      return mountConfig.get('firstObject');
    }
  }
}
