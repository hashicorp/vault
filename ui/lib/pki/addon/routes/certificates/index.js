/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import { withConfig } from 'pki/decorators/check-issuers';
import { hash } from 'rsvp';
import { getCliMessage } from 'pki/routes/overview';

@withConfig()
export default class PkiCertificatesIndexRoute extends Route {
  @service store;
  @service secretMountPath;

  async fetchCertificates() {
    try {
      return await this.store.query('pki/certificate/base', { backend: this.secretMountPath.currentPath });
    } catch (e) {
      if (e.httpStatus === 404) {
        return { parentModel: this.modelFor('certificates') };
      } else {
        throw e;
      }
    }
  }

  model() {
    return hash({
      hasConfig: this.shouldPromptConfig,
      certificates: this.fetchCertificates(),
      parentModel: this.modelFor('certificates'),
    });
  }

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);
    const certificates = resolvedModel.certificates;

    if (certificates?.length) controller.notConfiguredMessage = getCliMessage('certificates');
    else controller.notConfiguredMessage = getCliMessage();
  }
}
