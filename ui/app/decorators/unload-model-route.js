import { inject as service } from '@ember/service';
import { removeModel } from '../utils/remove-record';

/**
 * Takes array of model names which should be unloaded when moving away from
 * the route. Automatically unloads all `capabilities` and `control-group`
 * data before unloading the specified models.
 *
 *** basic example
 *
 * import Route from '@ember/route';
 * import withUnloadModelRoute from 'vault/decorators/unload-model-route';
 *
 * @withUnloadModelRoute(['policy/acl'])
 * class PolicyViewRoute extends Route { foo = null; }
 *
 */

export function withUnloadModelRoute(models) {
  return function decorator(SuperClass) {
    return class UnloadModelRoute extends SuperClass {
      @service store;

      static _models;

      constructor() {
        super(...arguments);
        if (!models || !Array.isArray(models)) {
          throw new Error('models must be an array of model types as strings');
        }
        this._models = models;
      }

      unloadModels() {
        const { _models } = this;
        _models.forEach((modelName) => {
          removeModel(this.store, modelName);
        });
      }

      deactivate() {
        ['capabilities', 'control-group', ...this._models].forEach((modelName) => {
          this.store.unloadAll(modelName);
        });
      }
    };
  };
}
