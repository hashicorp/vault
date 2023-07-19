/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import { hash } from 'rsvp';
import { pathIsFromDirectory, breadcrumbsForDirectory } from 'vault/lib/kv-breadcrumbs';

export default class KvSecretMetadataVersionsRoute extends Route {
  @service store;
  @service secretMountPath;

  model() {
    const backend = this.secretMountPath.get();
    const parentModel = this.modelFor('secret.metadata');
    return hash({
      backend,
      ...parentModel,
    });
  }

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);
    let breadcrumbsArray = [
      { label: 'secrets', route: 'secrets', linkExternal: true },
      { label: resolvedModel.backend, route: 'list' },
    ];

    if (pathIsFromDirectory(resolvedModel.name)) {
      breadcrumbsArray = [...breadcrumbsArray, ...breadcrumbsForDirectory(resolvedModel.name)];
    } else {
      breadcrumbsArray.push({ label: resolvedModel.name });
    }
    breadcrumbsArray.push({ label: 'version history' });
    controller.set('breadcrumbs', breadcrumbsArray);
  }
}
