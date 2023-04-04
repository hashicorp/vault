/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { withConfirmLeave } from 'core/decorators/confirm-leave';
import PkiKeyRoute from '../key';

@withConfirmLeave()
export default class PkiKeyEditRoute extends PkiKeyRoute {
  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);
    controller.breadcrumbs.push({ label: resolvedModel.id, route: 'keys.key.details' }, { label: 'edit' });
  }
}
