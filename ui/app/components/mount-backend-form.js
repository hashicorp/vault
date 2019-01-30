import { inject as service } from '@ember/service';
import { computed } from '@ember/object';
import Component from '@ember/component';
import { task } from 'ember-concurrency';
import { methods } from 'vault/helpers/mountable-auth-methods';
import { engines } from 'vault/helpers/mountable-secret-engines';

const METHODS = methods();
const ENGINES = engines();

export default Component.extend({
  store: service(),
  wizard: service(),
  flashMessages: service(),

  /*
   * @param Function
   * @public
   *
   * Optional param to call a function upon successfully mounting a backend
   *
   */
  onMountSuccess() {},
  onConfigError() {},
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

  showConfig: false,

  init() {
    this._super(...arguments);
    const type = this.get('mountType');
    const modelType = type === 'secret' ? 'secret-engine' : 'auth-method';
    const model = this.get('store').createRecord(modelType);
    this.set('mountModel', model);
  },

  mountTypes: computed('mountType', function() {
    return this.get('mountType') === 'secret' ? ENGINES : METHODS;
  }),

  willDestroy() {
    // if unsaved, we want to unload so it doesn't show up in the auth mount list
    this.get('mountModel').rollbackAttributes();
  },

  getConfigModelType(methodType) {
    let mountType = this.get('mountType');
    // will be something like secret-aws
    // or auth-azure
    let key = `${mountType}-${methodType}`;
    let noConfig = ['auth-approle', 'auth-alicloud'];
    if (mountType === 'secret' || noConfig.includes(key)) {
      return;
    }
    if (methodType === 'aws') {
      return 'auth-config/aws/client';
    }
    return `auth-config/${methodType}`;
  },

  changeConfigModel(methodType) {
    let mount = this.get('mountModel');
    if (this.get('mountType') === 'secret') {
      return;
    }
    let configRef = mount.hasMany('authConfigs').value();
    let currentConfig = configRef && configRef.get('firstObject');
    if (currentConfig) {
      // rollbackAttributes here will remove the the config model from the store
      // because `isNew` will be true
      currentConfig.rollbackAttributes();
      currentConfig.unloadRecord();
    }
    let configType = this.getConfigModelType(methodType);
    if (!configType) return;
    let config = this.get('store').createRecord(configType);
    config.set('backend', mount);
  },

  checkPathChange(type) {
    let mount = this.get('mountModel');
    let currentPath = mount.get('path');
    let list = this.get('mountTypes');
    // if the current path matches a type (meaning the user hasn't altered it),
    // change it here to match the new type
    let isUnchanged = list.findBy('type', currentPath);
    if (!currentPath || isUnchanged) {
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

    let mountType = this.get('mountType');
    mountType = mountType === 'secret' ? `${mountType}s engine` : `${mountType} method`;
    this.get('flashMessages').success(`Successfully mounted the ${type} ${mountType} at ${path}.`);
    if (this.get('mountType') === 'secret') {
      yield this.get('onMountSuccess')(type, path);
      return;
    }
    yield this.get('saveConfig').perform(mountModel);
  }).drop(),

  advanceWizard() {
    this.get('wizard').transitionFeatureMachine(
      this.get('wizard.featureState'),
      'CONTINUE',
      this.get('mountModel').get('type')
    );
  },
  saveConfig: task(function*(mountModel) {
    const configRef = mountModel.hasMany('authConfigs').value();
    const { type, path } = mountModel.getProperties('type', 'path');
    if (!configRef) {
      this.advanceWizard();
      yield this.get('onMountSuccess')(type, path);
      return;
    }
    const config = configRef.get('firstObject');
    try {
      if (config && Object.keys(config.changedAttributes()).length) {
        yield config.save();
        this.advanceWizard();
        this.get('flashMessages').success(
          `The config for ${type} ${this.get('mountType')} method at ${path} was saved successfully.`
        );
      }
      yield this.get('onMountSuccess')(type, path);
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
        this.get('wizard').set('componentState', value);
        this.changeConfigModel(value);
        this.checkPathChange(value);
      }
    },

    toggleShowConfig(value) {
      this.set('showConfig', value);
      if (value === true && this.get('wizard.featureState') === 'idle') {
        this.get('wizard').transitionFeatureMachine(
          this.get('wizard.featureState'),
          'CONTINUE',
          this.get('mountModel').get('type')
        );
      } else {
        this.get('wizard').transitionFeatureMachine(
          this.get('wizard.featureState'),
          'RESET',
          this.get('mountModel').get('type')
        );
      }
    },
  },
});
