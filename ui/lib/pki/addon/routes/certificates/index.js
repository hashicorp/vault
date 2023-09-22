/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import { withConfig } from 'pki/decorators/check-config';
import { hash } from 'rsvp';
import { getCliMessage } from 'pki/routes/overview';

@withConfig()
export default class PkiCertificatesIndexRoute extends Route {
  @service store;
  @service secretMountPath;

  queryParams = {
    page: {
      refreshModel: true,
    },
  };

  async fetchCertificates(params) {
    try {
      const page = Number(params.page) || 1;
      return await this.store.lazyPaginatedQuery('pki/certificate/base', {
        backend: this.secretMountPath.currentPath,
        responsePath: 'data.keys',
        page,
        skipCache: page === 1,
      });
    } catch (e) {
      if (e.httpStatus === 404) {
        return { parentModel: this.modelFor('certificates') };
      }
      throw e;
    }
  }

  model(params) {
    return hash({
      hasConfig: this.shouldPromptConfig,
      certificates: this.fetchCertificates(params),
      parentModel: this.modelFor('certificates'),
      pageFilter: params.pageFilter,
    });
  }

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);
    const certificates = resolvedModel.certificates;

    if (certificates?.length) controller.notConfiguredMessage = getCliMessage('certificates');
    else controller.notConfiguredMessage = getCliMessage();
  }

  resetController(controller, isExiting) {
    if (isExiting) {
      controller.set('page', undefined);
    }
  }
}
