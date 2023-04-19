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

  async fetchMountConfig(path) {
    const mountConfig = await this.store.query('secret-engine', { path });

    if (mountConfig) {
      return mountConfig.get('firstObject');
    }
  }

  async model() {
    const backend = this.secretMountPath.currentPath;

    return hash({
      hasConfig: this.shouldPromptConfig,
      engine: this.modelFor('application'),
      urls: this.fetchUrls(backend),
      crl: this.fetchCrl(backend),
      mountConfig: this.fetchMountConfig(backend),
      issuerModel: this.store.createRecord('pki/issuer'),
    });
  }

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);
  }
}
