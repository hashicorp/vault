/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Route from '@ember/routing/route';
import { pathIsFromDirectory, breadcrumbsForDirectory } from 'vault/lib/kv-breadcrumbs';
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
    const queryVersion =
      // first check version data exists (user has READ permission for kv/data)
      // only query version if params have changed
      parentModel.secret.version && params.version && parentModel.secret.version !== Number(params.version);

    if (queryVersion) {
      // query params have changed by selecting a different version from the dropdown
      // fire off new request for that version's secret data
      const { backend, path } = parentModel;
      return hash({
        ...parentModel,
        secret: this.store.queryRecord('kv/data', { backend, path, version: params.version }).catch(() => {
          // return empty record to access capability getters on model
          return this.store.createRecord('kv/data', { backend, path });
        }),
      });
    }
    return parentModel;
  }

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);

    let breadcrumbsArray = [
      { label: 'secrets', route: 'secrets', linkExternal: true },
      { label: resolvedModel.backend, route: 'list' },
    ];

    if (pathIsFromDirectory(resolvedModel.path)) {
      breadcrumbsArray = [...breadcrumbsArray, ...breadcrumbsForDirectory(resolvedModel.path, true)];
    } else {
      breadcrumbsArray.push({ label: resolvedModel.path });
    }

    controller.breadcrumbs = breadcrumbsArray;
    controller.set('version', resolvedModel.secret.version);
  }

  resetController(controller, isExiting) {
    if (isExiting) {
      controller.set('version', null);
    }
  }
}
