/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { kvErrorHandler } from 'kv/utils/kv-error-handler';

export default class KvSecretDetailsRoute extends Route {
  @service api;

  queryParams = {
    version: {
      refreshModel: true,
    },
  };

  async model(params) {
    const parentModel = this.modelFor('secret');
    const { backend, path } = parentModel;
    let secret;
    // if a version is selected from the dropdown it triggers a model refresh
    // and we fire off new request for that version's secret data
    try {
      const initOverride = params.version
        ? (context) => this.api.addQueryParams(context, { version: params.version })
        : undefined;
      const { data, metadata } = await this.api.secrets.kvV2Read(path, backend, initOverride);
      secret = { secretData: data, ...metadata };
    } catch (error) {
      const { status, response } = await this.api.parseError(error);
      const { data, metadata, failReadErrorCode } = kvErrorHandler(status, response);
      secret = failReadErrorCode ? { failReadErrorCode } : { secretData: data, ...metadata };
    }

    return {
      ...parentModel,
      secret,
    };
  }

  // breadcrumbs are set in details/index.js
  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);
    const { version } = this.paramsFor(this.routeName);
    controller.set('version', resolvedModel.secret.version || version);
  }

  resetController(controller, isExiting) {
    if (isExiting) {
      controller.set('version', null);
    }
  }
}
