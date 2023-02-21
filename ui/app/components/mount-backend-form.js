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
  @service flashMessages;

  // validation related properties
  @tracked modelValidations = null;
  @tracked invalidFormAlert = null;

  @tracked errorMessage = '';

  willDestroy() {
    // if unsaved, we want to unload so it doesn't show up in the auth mount list
    super.willDestroy(...arguments);
    if (this.args.mountModel) {
      const method = this.args.mountModel.isNew ? 'unloadRecord' : 'rollbackAttributes';
      this.args.mountModel[method]();
    }
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

  async showWarningsForKvv2() {
    try {
      const capabilities = await this.store.findRecord('capabilities', `${this.args.mountModel.path}/config`);
      if (!capabilities?.canUpdate) {
        // config error is not thrown from secret-engine adapter, so handling here
        this.flashMessages.warning(
          'You do not have access to the config endpoint. The secret engine was mounted, but the configuration settings were not saved.'
        );
        // remove the config data from the model otherwise it will persist in the store even though network request failed.
        [
          this.args.mountModel.maxVersions,
          this.args.mountModel.casRequired,
          this.args.mountModel.deleteVersionAfter,
        ] = [0, false, 0];
      }
    } catch (e) {
      // Show different warning if we're not sure the config saved
      this.flashMessages.warning(
        'You may not have access to the config endpoint. The secret engine was mounted, but the configuration settings may not be saved.'
      );
    }
    return;
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
    if (mountModel.isV2KV && updatesConfig) {
      yield this.showWarningsForKvv2();
    }
    this.flashMessages.success(
      `Successfully mounted the ${type} ${
        this.args.mountType === 'secret' ? 'secrets engine' : 'auth method'
      } at ${path}.`
    );
    yield this.args.onMountSuccess(type, path);
    return;
  }

  @action
  onKeyUp(name, value) {
    this.args.mountModel[name] = value;
  }

  @action
  setMountType(value) {
    this.args.mountModel.type = value;
    this.checkPathChange(value);
  }
}
