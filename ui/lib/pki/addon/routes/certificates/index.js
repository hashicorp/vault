/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import PkiOverviewRoute from '../overview';
import { inject as service } from '@ember/service';
import { hash } from 'rsvp';

export default class PkiCertificatesIndexRoute extends PkiOverviewRoute {
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
      hasConfig: this.hasConfig(),
      certificates: this.fetchCertificates(),
      parentModel: this.modelFor('certificates'),
    });
  }
}
