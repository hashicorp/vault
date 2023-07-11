/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Route from '@ember/routing/route';
import { action } from '@ember/object';
import { inject as service } from '@ember/service';
import { hash } from 'rsvp';
import { normalizePath } from 'vault/utils/path-encoding-helpers';

export default class KvSecretsListRoute extends Route {
  @service store;
  @service secretMountPath;

  nestedSecretName() {
    const { name } = this.paramsFor('list');
    return name ? normalizePath(name) : '';
  }

  // beforeModel() {
  //   const secret = this.secretParam();
  //   if (secret.endsWith('/')) {
  //     this.router.replaceWith('secrets', secret + '/');
  //   }
  // }
  model() {
    // TODO add filtering and return model for query on kv/metadata.
    const nestedSecret = this.nestedSecretName() || '';
    const backend = this.secretMountPath.currentPath;
    const secrets = this.store.query('kv/metadata', { backend, nestedSecret }).catch((err) => {
      if (err.httpStatus === 404) {
        return [];
      } else {
        throw err;
      }
    });
    return hash({
      nestedSecret,
      secrets,
      backend,
    });
  }

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);
    controller.set('model', resolvedModel.secrets);
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
