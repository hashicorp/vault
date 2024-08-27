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
    if (params.version) {
      // we don't send a version param when the route initially loads. if there is a version,
      // it has been selected from the dropdown. fire off new request for that version's secret data
      return hash({
        ...parentModel,
        secret: this.store.queryRecord('kv/data', { backend, path, version: params.version }),
      });
    }
    return hash({
      ...parentModel,
      secret: this.store.queryRecord('kv/data', { backend, path }),
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
