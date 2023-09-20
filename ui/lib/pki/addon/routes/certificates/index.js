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

  queryParams = {
    pageFilter: {
      refreshModel: true,
    },
    currentPage: {
      refreshModel: true,
    },
  };

  async fetchCertificates(params) {
    try {
      return await this.store.lazyPaginatedQuery('pki/certificate/base', {
        backend: this.secretMountPath.currentPath,
        responsePath: 'data.keys',
        page: Number(params.currentPage) || 1,
        pageFilter: params.pageFilter,
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
      controller.set('pageFilter', undefined);
      controller.set('currentPage', undefined);
    }
  }
}
