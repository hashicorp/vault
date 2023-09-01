/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
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
    // Only fetch versioned data if selected version does not match parent (current) version
    // and parentModel.secret has failReadErrorCode since permissions aren't version specific
    if (
      params.version &&
      parentModel.secret.version !== params.version &&
      !parentModel.secret.failReadErrorCode
    ) {
      // query params have changed by selecting a different version from the dropdown
      // fire off new request for that version's secret data
      const { backend, path } = parentModel;
      return hash({
        ...parentModel,
        secret: this.store.queryRecord('kv/data', { backend, path, version: params.version }),
      });
    }
    return parentModel;
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
