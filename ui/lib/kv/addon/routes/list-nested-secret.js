/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

/**
 * We have two routes for the list view. While this file does the logic the associated template is list.hbs.
 */

import Route from '@ember/routing/route';
import { action } from '@ember/object';
import { inject as service } from '@ember/service';
import { hash } from 'rsvp';
import { normalizePath } from 'vault/utils/path-encoding-helpers';

export default class KvSecretsListRoute extends Route {
  @service store;
  @service router;
  @service secretMountPath;

  getNestedSecretFromUrlParam() {
    const { nestedSecret } = this.paramsFor('list-nested-secret');
    return nestedSecret ? normalizePath(nestedSecret) : '';
  }

  // beforeModel() {
  //   const nestedSecret = this.getNestedSecretFromUrlParam();
  //   if (this.routeName === 'list' && nestedSecret.endsWith('/')) {
  //     // this.router.replaceWith('vault.cluster.secrets.backend.kv.list', nestedSecret + '/');
  //   }
  // }

  model() {
    // TODO add filtering and return model for query on kv/metadata.
    const nestedSecret = this.getNestedSecretFromUrlParam();
    const backend = this.secretMountPath.currentPath;
    const arrayOfSecretModels = this.store.query('kv/metadata', { backend, nestedSecret }).catch((err) => {
      if (err.httpStatus === 404) {
        return [];
      } else {
        throw err;
      }
    });
    return hash({
      arrayOfSecretModels,
      backend,
      routeName: this.routeName,
    });
  }

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);
    controller.set('model', resolvedModel.arrayOfSecretModels);
    controller.routeName = resolvedModel.routeName;
    controller.pageTitle = resolvedModel.backend;
    controller.breadcrumbs = [
      { label: 'secrets', route: 'secrets', linkExternal: true },
      { label: resolvedModel.backend },
    ];
  }
  @action
  willTransition() {
    window.scrollTo(0, 0);
    // ARG TODO not working
    this.store.clearDataset('kv/metadata');
  }
}
