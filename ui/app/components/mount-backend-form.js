import Ember from 'ember';
import { task } from 'ember-concurrency';
import { methods } from 'vault/helpers/mountable-auth-methods';

const { inject } = Ember;
const METHODS = methods();

export default Ember.Component.extend({
  store: inject.service(),
  flashMessages: inject.service(),
  routing: inject.service('-routing'),

  /*
   * @param Function
   * @public
   *
   * Optional param to call a function upon successfully mounting a backend
   *
   */
  onMountSuccess: () => {},
  onConfigError: () => {},
  /*
   * @param String
   * @public
   * the type of backend we want to mount
   * defaults to `auth`
   *
   */
  mountType: 'auth',

  /*
   *
   * @param DS.Model
   * @private
   * Ember Data model corresponding to the `mountType`.
   * Created and set during `init`
   *
   */
  mountModel: null,

  init() {
    this._super(...arguments);
    const type = this.get('mountType');
    const modelType = type === 'secret' ? 'secret-engine' : 'auth-method';
    const model = this.get('store').createRecord(modelType);
    this.set('mountModel', model);
    this.changeConfigModel(model.get('type'));
  },

  willDestroy() {
    // if unsaved, we want to unload so it doesn't show up in the auth mount list
    this.get('mountModel').rollbackAttributes();
  },

  getConfigModelType(methodType) {
    let noConfig = ['approle'];
    if (noConfig.includes(methodType)) {
      return;
    }
    if (methodType === 'aws') {
      return 'auth-config/aws/client';
    }
    return `auth-config/${methodType}`;
  },

  changeConfigModel(methodType) {
    const mount = this.get('mountModel');
    const configRef = mount.hasMany('authConfigs').value();
    const currentConfig = configRef.get('firstObject');
    if (currentConfig) {
      // rollbackAttributes here will remove the the config model from the store
      // because `isNew` will be true
      currentConfig.rollbackAttributes();
    }
    const configType = this.getConfigModelType(methodType);
    if (!configType) return;
    const config = this.get('store').createRecord(configType);
    config.set('backend', mount);
  },

  checkPathChange(type) {
    const mount = this.get('mountModel');
    const currentPath = mount.get('path');
    // if the current path matches a type (meaning the user hasn't altered it),
    // change it here to match the new type
    const isUnchanged = METHODS.findBy('type', currentPath);
    if (isUnchanged) {
      mount.set('path', type);
    }
  },

  mountBackend: task(function*() {
    const mountModel = this.get('mountModel');
    const { type, path } = mountModel.getProperties('type', 'path');
    try {
      yield mountModel.save();
    } catch (err) {
      // err will display via model state
      return;
    }
    this.get('flashMessages').success(
      `Successfully mounted ${type} ${this.get('mountType')} method at ${path}.`
    );
    yield this.get('saveConfig').perform(mountModel);
  }).drop(),

  saveConfig: task(function*(mountModel) {
    const configRef = mountModel.hasMany('authConfigs').value();
    const config = configRef.get('firstObject');
    const { type, path } = mountModel.getProperties('type', 'path');
    try {
      if (config && Object.keys(config.changedAttributes()).length) {
        yield config.save();
        this.get('flashMessages').success(
          `The config for ${type} ${this.get('mountType')} method at ${path} was saved successfully.`
        );
      }
      yield this.get('onMountSuccess')();
    } catch (err) {
      this.get('flashMessages').danger(
        `There was an error saving the configuration for ${type} ${this.get(
          'mountType'
        )} method at ${path}. ${err.errors.join(' ')}`
      );
      yield this.get('onConfigError')(mountModel.id);
    }
  }).drop(),

  actions: {
    onTypeChange(path, value) {
      if (path === 'type') {
        this.changeConfigModel(value);
        this.checkPathChange(value);
      }
    },
  },
});
