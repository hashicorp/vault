/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import PkiIssuerRoute from '../issuer';
import { withConfirmLeave } from 'core/decorators/confirm-leave';

@withConfirmLeave()
export default class PkiIssuerCrossSignRoute extends PkiIssuerRoute {
  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);
    controller.breadcrumbs.push(
      { label: resolvedModel.id, route: 'issuers.issuer.details' },
      { label: 'cross-sign' }
    );
  }
}
