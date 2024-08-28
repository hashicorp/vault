/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { hash } from 'rsvp';

export default class KvSecretDetailsRoute extends Route {
  @service store;

  queryParams = {
    version: {
      refreshModel: true,
    },
  };

  model(params) {
    const parentModel = this.modelFor('secret');
    const { backend, path } = parentModel;
    const query = { backend, path };
    // if a version is selected from the dropdown it triggers a model refresh
    // and we fire off new request for that version's secret data
    if (params.version) {
      query.version = params.version;
    }
    return hash({
      ...parentModel,
      secret: this.store.queryRecord('kv/data', query),
    });
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
