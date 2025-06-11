import Route from '@ember/routing/route';
import { CONFIGURABLE_SECRET_ENGINES, allEngines } from 'vault/helpers/mountable-secret-engines';

export default class SecretsBackendConfigurationIndexRoute extends Route {
  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);
    controller.typeDisplay = allEngines().find(
      (engine) => engine.type === resolvedModel.secretsEngine.type
    )?.displayName;
    controller.isConfigurable = CONFIGURABLE_SECRET_ENGINES.includes(resolvedModel.secretsEngine.type);
    controller.modelId = resolvedModel.secretsEngine.id;
  }
}
