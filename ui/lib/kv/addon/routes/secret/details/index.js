/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Route from '@ember/routing/route';
import { pathIsFromDirectory, breadcrumbsForDirectory } from 'vault/lib/kv-breadcrumbs';

export default class KvSecretDetailsIndexRoute extends Route {
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
  }
}
