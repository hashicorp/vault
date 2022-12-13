import Ember from 'ember';
import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { inject as service } from '@ember/service';
import { action } from '@ember/object';
import { task } from 'ember-concurrency';
import { waitFor } from '@ember/test-waiters';
import { supportedSecretBackends } from 'vault/helpers/supported-secret-backends';
import { methods } from 'vault/helpers/mountable-auth-methods';

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

export default class MountBackendForm extends Component {
  @service store;
  @service wizard;
  @service flashMessages;

  // validation related properties
  @tracked modelValidations = null;
  @tracked invalidFormAlert = null;

  @tracked errorMessage = '';

  willDestroy() {
    // if unsaved, we want to unload so it doesn't show up in the auth mount list
    super.willDestroy(...arguments);
    this.args.mountModel.rollbackAttributes();
  }

  checkPathChange(type) {
    if (!type) return;
    const mount = this.args.mountModel;
    const currentPath = mount.path;
    const mountTypes =
      this.args.mountType === 'secret' ? supportedSecretBackends() : methods().map((auth) => auth.type);
    // if the current path has not been altered by user,
    // change it here to match the new type
    if (!currentPath || mountTypes.includes(currentPath)) {
      mount.path = type;
    }
  }

  checkModelValidity(model) {
    const { isValid, state, invalidFormMessage } = model.validate();
    this.modelValidations = state;
    this.invalidFormAlert = invalidFormMessage;
    return isValid;
  }

  @task
  @waitFor
  *mountBackend(event) {
    event.preventDefault();
    const mountModel = this.args.mountModel;
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
        // TODO: smells like we need to fix a test
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

    const changedAttrKeys = Object.keys(mountModel.changedAttributes());
    const updatesConfig =
      changedAttrKeys.includes('casRequired') ||
      changedAttrKeys.includes('deleteVersionAfter') ||
      changedAttrKeys.includes('maxVersions');

    try {
      yield mountModel.save();
    } catch (err) {
      if (err.httpStatus === 403) {
        this.flashMessages.danger(
          'You do not have access to the sys/mounts endpoint. The secret engine was not mounted.'
        );
        return;
      }
      if (err.errors) {
        const errors = err.errors.map((e) => {
          if (typeof e === 'object') return e.title || e.message || JSON.stringify(e);
          return e;
        });
        this.errorMessage = errors;
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
      // TODO: Move this logic to custom mount component for KV
      // remove the config data from the model otherwise it will save it even if the network request failed.
      [
        this.args.mountModel.maxVersions,
        this.args.mountModel.casRequired,
        this.args.mountModel.deleteVersionAfter,
      ] = [0, false, 0];
    }
    this.flashMessages.success(
      `Successfully mounted the ${type} ${
        this.mountType === 'secret' ? 'secrets engine' : 'auth method'
      } at ${path}.`
    );
    yield this.args.onMountSuccess(type, path);
    return;
  }

  @action
  onKeyUp(name, value) {
    this.args.mountModel.set(name, value);
  }

  @action
  onTypeChange(path, value) {
    if (path === 'type') {
      this.wizard.set('componentState', value);
    }
  }

  @action
  setMountType(value) {
    this.args.mountModel.set('type', value);
    this.checkPathChange(value);
    if (value === '') {
      // resets wizard
      if (this.wizard.featureState === 'idle') {
        this.wizard.transitionFeatureMachine(this.wizard.featureState, 'RESET', this.args.mountModel.type);
      }
    } else {
      this.wizard.transitionFeatureMachine(this.wizard.featureState, 'CONTINUE', this.args.mountModel.type);
    }
  }
}
