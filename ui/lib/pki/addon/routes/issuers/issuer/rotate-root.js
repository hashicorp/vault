/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import PkiIssuerIndexRoute from './index';
import { inject as service } from '@ember/service';
import { withConfirmLeave } from 'core/decorators/confirm-leave';

@withConfirmLeave()
export default class PkiIssuerRotateRootRoute extends PkiIssuerIndexRoute {
  @service secretMountPath;
  @service store;

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);
    controller.breadcrumbs = [
      { label: 'secrets', route: 'secrets', linkExternal: true },
      { label: this.secretMountPath.currentPath, route: 'overview' },
      { label: 'issuers', route: 'issuers.index' },
      { label: resolvedModel.id, route: 'issuers.issuer.details' },
      { label: 'rotate root' },
    ];
  }
}
