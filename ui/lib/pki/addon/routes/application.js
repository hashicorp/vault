/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';

export default class PkiRoute extends Route {
  @service pathHelp;
  @service secretMountPath;

  beforeModel() {
    // Must call this promise before the model hook otherwise the model doesn't hydrate from OpenAPI correctly.
    // only needs to be called once to add the openAPI attributes to the model prototype
    return this.pathHelp.getNewModel('pki/role', this.secretMountPath.currentPath);
  }
}
