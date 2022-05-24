import Ember from 'ember';
import { inject as service } from '@ember/service';
import { computed } from '@ember/object';
import Component from '@ember/component';
import { task } from 'ember-concurrency';
import { methods } from 'vault/helpers/mountable-auth-methods';
import { engines, KMIP, TRANSFORM, KEYMGMT } from 'vault/helpers/mountable-secret-engines';
import { waitFor } from '@ember/test-waiters';

const METHODS = methods();
const ENGINES = engines();

export default Component.extend({
  store: service(),
  wizard: service(),
  flashMessages: service(),
  version: service(),

  /*
   * @param Function
   * @public
   *
   * Optional param to call a function upon successfully mounting a backend
   *
   */
  onMountSuccess() {},
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

  showEnable: false,

  // validation related properties
  modelValidations: null,
  invalidFormAlert: null,

  mountIssue: false,

  init() {
    this._super(...arguments);
    const type = this.mountType;
    const modelType = type === 'secret' ? 'secret-engine' : 'auth-method';
    const model = this.store.createRecord(modelType);
    this.set('mountModel', model);
  },

  mountTypes: computed('engines', 'mountType', function () {
    return this.mountType === 'secret' ? this.engines : METHODS;
  }),

  engines: computed('version.{features[],isEnterprise}', function () {
    if (this.version.isEnterprise) {
      return ENGINES.concat([KMIP, TRANSFORM, KEYMGMT]);
    }
    return ENGINES;
  }),

  willDestroy() {
    this._super(...arguments);
    // if unsaved, we want to unload so it doesn't show up in the auth mount list
    this.mountModel.rollbackAttributes();
  },

  checkPathChange(type) {
    let mount = this.mountModel;
    let currentPath = mount.path;
    let list = this.mountTypes;
    // if the current path matches a type (meaning the user hasn't altered it),
    // change it here to match the new type
    let isUnchanged = list.findBy('type', currentPath);
    if (!currentPath || isUnchanged) {
      mount.set('path', type);
    }
  },

  checkModelValidity(model) {
    const { isValid, state, invalidFormMessage } = model.validate();
    this.setProperties({
      modelValidations: state,
      invalidFormAlert: invalidFormMessage,
    });

    return isValid;
  },

  mountBackend: task(
    waitFor(function* () {
      const mountModel = this.mountModel;
      const { type, path } = mountModel;
      // only submit form if validations pass
      if (!this.checkModelValidity(mountModel)) {
        return;
      }
      let capabilities = null;
      try {
        capabilities = yield this.store.findRecord('capabilities', `${path}/config`);
      } catch (err) {
        if (Ember.testing) {
          //captures mount-backend-form component test
          yield mountModel.save();
          let mountType = this.mountType;
          mountType = mountType === 'secret' ? `${mountType}s engine` : `${mountType} method`;
          this.flashMessages.success(`Successfully mounted the ${type} ${mountType} at ${path}.`);
          yield this.onMountSuccess(type, path);
          return;
        } else {
          throw err;
        }
      }

      let changedAttrKeys = Object.keys(mountModel.changedAttributes());
      let updatesConfig =
        changedAttrKeys.includes('casRequired') ||
        changedAttrKeys.includes('deleteVersionAfter') ||
        changedAttrKeys.includes('maxVersions');

      try {
        yield mountModel.save();
      } catch (err) {
        if (err.httpStatus === 403) {
          this.mountIssue = true;
          this.flashMessages.danger(
            'You do not have access to the sys/mounts endpoint. The secret engine was not mounted.'
          );
          return;
        }
        if (err.errors) {
          let errors = err.errors.map((e) => {
            if (typeof e === 'object') return e.title || e.message || JSON.stringify(e);
            return e;
          });
          this.set('errors', errors);
        } else if (err.message) {
          this.set('errorMessage', err.message);
        } else {
          this.set('errorMessage', 'An error occurred, check the vault logs.');
        }
        return;
      }
      // mountModel must be after the save
      if (mountModel.isV2KV && updatesConfig && !capabilities.get('canUpdate')) {
        // config error is not thrown from secret-engine adapter, so handling here
        this.flashMessages.warning(
          'You do not have access to the config endpoint. The secret engine was mounted, but the configuration settings were not saved.'
        );
        // remove the config data from the model otherwise it will save it even if the network request failed.
        [this.mountModel.maxVersions, this.mountModel.casRequired, this.mountModel.deleteVersionAfter] = [
          0,
          false,
          0,
        ];
      }
      let mountType = this.mountType;
      mountType = mountType === 'secret' ? `${mountType}s engine` : `${mountType} method`;
      this.flashMessages.success(`Successfully mounted the ${type} ${mountType} at ${path}.`);
      yield this.onMountSuccess(type, path);
      return;
    })
  ).drop(),

  actions: {
    onKeyUp(name, value) {
      this.mountModel.set(name, value);
    },

    onTypeChange(path, value) {
      if (path === 'type') {
        this.wizard.set('componentState', value);
        this.checkPathChange(value);
      }
    },

    toggleShowEnable(value) {
      this.set('showEnable', value);
      if (value === true && this.wizard.featureState === 'idle') {
        this.wizard.transitionFeatureMachine(this.wizard.featureState, 'CONTINUE', this.mountModel.type);
      } else {
        this.wizard.transitionFeatureMachine(this.wizard.featureState, 'RESET', this.mountModel.type);
      }
    },
  },
});
