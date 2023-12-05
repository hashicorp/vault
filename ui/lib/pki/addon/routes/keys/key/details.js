/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import PkiKeyRoute from '../key';

export default class PkiKeyDetailsRoute extends PkiKeyRoute {
  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);
    controller.breadcrumbs.push({ label: resolvedModel.id });
  }
}
