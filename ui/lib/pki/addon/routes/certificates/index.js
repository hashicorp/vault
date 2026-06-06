/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { withConfig } from 'pki/decorators/check-issuers';
import { getCliMessage } from 'pki/routes/overview';
import { SecretsApiPkiListCertsListEnum } from '@hashicorp/vault-client-typescript';
import { paginate } from 'core/utils/paginate-list';

@withConfig()
export default class PkiCertificatesIndexRoute extends Route {
  @service secretMountPath;
  @service api;

  queryParams = {
    page: {
      refreshModel: true,
    },
  };

  async model(params) {
    const model = {
      hasConfig: this.pkiMountHasConfig,
      parentModel: this.modelFor('certificates'),
      pageFilter: params.pageFilter,
      certificates: [],
    };

    try {
      const page = Number(params.page) || 1;
      const { keys: certificates } = await this.api.secrets.pkiListCerts(
        this.secretMountPath.currentPath,
        SecretsApiPkiListCertsListEnum.TRUE
      );
      model.certificates = paginate(certificates, { page });
    } catch (e) {
      if (e.response.status !== 404) {
        throw e;
      }
    }

    return model;
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
