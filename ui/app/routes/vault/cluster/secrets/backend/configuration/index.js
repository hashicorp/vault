/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import engineDisplayData from 'vault/helpers/engines-display-data';

export default class SecretsBackendConfigurationIndexRoute extends Route {
  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);
    const engine = engineDisplayData(resolvedModel.secretsEngine.type);
    controller.typeDisplay = engine.displayName;
    controller.isConfigurable = engine.isConfigurable ?? false;
    controller.modelId = resolvedModel.secretsEngine.id;
  }
}
