import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';

/**
 * For replacing unload-model-route mixin without converting to modern type.
 * Takes array of model names which should be unloaded when moving away from
 * the route. Automatically unloads all `capabilities` and `control-group`
 * data before unloading the specified models.
 */
export default class UnloadModelRoute extends Route {
  @service store;

  constructor() {
    super(...arguments);
    if (!this.modelTypes || !Array.isArray(this.modelTypes)) {
      throw new Error('Set modelTypes on the route. Must be an array of strings');
    }
  }

  deactivate() {
    ['capabilities', 'control-group', ...this.modelTypes].forEach((modelName) => {
      this.store.unloadAll(modelName);
    });
  }
}
