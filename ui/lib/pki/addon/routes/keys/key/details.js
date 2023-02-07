/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import PkiKeysIndexRoute from '.';

export default class PkiKeyDetailsRoute extends PkiKeysIndexRoute {
  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);
    controller.breadcrumbs.push({ label: resolvedModel.id });
  }
}
