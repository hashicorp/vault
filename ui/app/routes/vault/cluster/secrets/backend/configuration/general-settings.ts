/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import engineDisplayData from 'vault/helpers/engines-display-data';

import type SecretEngineModel from 'vault/models/secret-engine';
import type { Breadcrumb } from 'vault/vault/app-types';
import type Controller from '@ember/controller';

interface RouteModel {
  secretsEngine: SecretEngineModel;
}
interface RouteController extends Controller {
  breadcrumbs: Array<Breadcrumb>;
  model: RouteModel;
  typeDisplay: string;
  isConfigurable: boolean;
  modelId: string;
}

export default class SecretsBackendConfigurationGeneralSettingsRoute extends Route {
  setupController(controller: RouteController, resolvedModel: RouteModel) {
    super.setupController(controller, resolvedModel);
    const engine = engineDisplayData(resolvedModel.secretsEngine.type);
    controller.typeDisplay = engine.displayName;
    controller.isConfigurable = engine.isConfigurable ?? false;
    controller.modelId = resolvedModel.secretsEngine.id;
  }
}
