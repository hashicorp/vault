import Ember from 'ember';
import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { inject as service } from '@ember/service';
import { action, setProperties } from '@ember/object';
import { task } from 'ember-concurrency';
import { methods } from 'vault/helpers/mountable-auth-methods';
import { engines, KMIP, TRANSFORM, KEYMGMT } from 'vault/helpers/mountable-secret-engines';
import { waitFor } from '@ember/test-waiters';

/**
 * @module MountBackendForm
 * The `MountBackendForm` is used to mount either a secret or auth backend.
 *
 * @example ```js
 *   <MountBackendForm @mountType="secret" @onMountSuccess={{this.onMountSuccess}} />```
 *
 * @param {function} onMountSuccess - A function that transitions once the Mount has been successfully posted.
 * @param {string} [mountType=auth] - The type of backend we want to mount.
 *
 */

const METHODS = methods();
const ENGINES = engines();

export default class MountBackendForm extends Component {
  @service store;
  @service flashMessages;
  @service version;

  get mountType() {
    return this.args.mountType || 'auth';
  }

  @tracked mountModel = null;
  @tracked showEnable = false;

  // validation related properties
  @tracked modelValidations = null;
  @tracked invalidFormAlert = null;

  @tracked mountIssue = false;

  @tracked errors = '';
  @tracked errorMessage = '';

  constructor() {
    super(...arguments);
    const type = this.args.mountType || 'auth';
    const modelType = type === 'secret' ? 'secret-engine' : 'auth-method';
    const model = this.store.createRecord(modelType);
    this.mountModel = model;
  }

  get mountTypes() {
    return this.mountType === 'secret' ? this.engines : METHODS;
  }

  get engines() {
    if (this.version.isEnterprise) {
      return ENGINES.concat([KMIP, TRANSFORM, KEYMGMT]);
    }
    return ENGINES;
  }

  willDestroy() {
    // if unsaved, we want to unload so it doesn't show up in the auth mount list
    super.willDestroy(...arguments);
    if (this.mountModel) {
      const method = this.mountModel.isNew ? 'unloadRecord' : 'rollbackAttributes';
      this.mountModel[method]();
    }
  }

  checkPathChange(type) {
    let mount = this.mountModel;
    let currentPath = mount.path;
    let list = this.mountTypes;
    // if the current path matches a type (meaning the user hasn't altered it),
    // change it here to match the new type
    let isUnchanged = list.findBy('type', currentPath);
    if (!currentPath || isUnchanged) {
      mount.path = type;
    }
  }

  checkModelValidity(model) {
    const { isValid, state, invalidFormMessage } = model.validate();
    setProperties(this, {
      modelValidations: state,
      invalidFormAlert: invalidFormMessage,
    });

    return isValid;
  }

  @task
  @waitFor
  *mountBackend() {
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
        yield this.args.onMountSuccess(type, path);
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
        this.errors = errors;
      } else if (err.message) {
        this.errorMessage = err.message;
      } else {
        this.errorMessage = 'An error occurred, check the vault logs.';
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
    yield this.args.onMountSuccess(type, path);
    return;
  }

  @action
  onKeyUp(name, value) {
    this.mountModel.set(name, value);
  }

  @action
  onTypeChange(path, value) {
    if (path === 'type') {
      this.checkPathChange(value);
    }
  }

  @action
  toggleShowEnable(value) {
    this.showEnable = value;
  }
}
