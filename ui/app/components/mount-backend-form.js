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

  init() {
    this._super(...arguments);
    const type = this.mountType;
    const modelType = type === 'secret' ? 'secret-engine' : 'auth-method';
    const model = this.store.createRecord(modelType);
    this.set('mountModel', model);
  },

  mountTypes: computed('mountType', function() {
    return this.mountType === 'secret' ? ENGINES : METHODS;
  }),

  willDestroy() {
    // if unsaved, we want to unload so it doesn't show up in the auth mount list
    this.get('mountModel').rollbackAttributes();
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

  mountBackend: task(function*() {
    const mountModel = this.mountModel;
    const { type, path } = mountModel.getProperties('type', 'path');
    try {
      yield mountModel.save();
    } catch (err) {
      // err will display via model state
      return;
    }

    let mountType = this.mountType;
    mountType = mountType === 'secret' ? `${mountType}s engine` : `${mountType} method`;
    this.flashMessages.success(`Successfully mounted the ${type} ${mountType} at ${path}.`);
    yield this.onMountSuccess(type, path);
    return;
  }).drop(),

  actions: {
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
