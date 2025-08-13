/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import engineDisplayData from 'vault/helpers/engines-display-data';
export default class SecretsBackendConfigurationIndexRoute extends Route {
  @service router;

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);
    const engine = engineDisplayData(resolvedModel.secretsEngine.type);
    controller.typeDisplay = engine?.displayName;
    controller.isConfigurable = engine?.isConfigurable ?? false;
    controller.glyph = engine?.glyph;
    controller.modelId = resolvedModel.secretsEngine.id;
  }

  afterModel(model) {
    const engine = engineDisplayData(model.secretsEngine.type);

    if (!engine?.isOldEngine) {
      return this.router.replaceWith('vault.cluster.secrets.backend.configuration.general-settings');
    }
  }
}
