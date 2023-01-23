import { action } from '@ember/object';
import { service } from '@ember/service';
import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import errorMessage from 'vault/utils/error-message';

/**
 * @module PkiGenerateRoot
 * PkiGenerateRoot shows only the fields valid for the generate root endpoint.
 * This form handles the model save and rollback actions, and will call the passed
 * onSave and onCancel args for transition (passed from parent).
 * NOTE: this component is not TS because decorator-added parameters (eg validator and
 * formFields) aren't recognized on the model.
 *
 * @example
 * ```js
 * <PkiGenerateRoot @model={{this.model}} @onCancel={{transition-to "vault.cluster"}} @onSave={{transition-to "vault.cluster.secrets"}} @adapterOptions={{hash actionType="import" useIssuer=false}} />
 * ```
 *
 * @param {Object} model - pki/action model.
 * @callback onCancel - Callback triggered when cancel button is clicked, after model is unloaded
 * @callback onSave - Callback triggered after model save success.
 * @param {Object} adapterOptions - object passed as adapterOptions on the model.save method
 */
export default class PkiGenerateRootComponent extends Component {
  @service flashMessages;
  @tracked showGroup = null;
  @tracked modelValidations = null;
  @tracked errorBanner = '';
  @tracked invalidFormAlert = '';

  get defaultFields() {
    return [
      'type',
      'commonName',
      'issuerName',
      'customTtl',
      'notBeforeDuration',
      'format',
      'permittedDnsDomains',
      'maxPathLength',
    ];
  }

  @action cancel() {
    // Generate root form will always have a new model
    this.args.model.unloadRecord();
    this.args.onCancel();
  }

  @action
  checkFormValidity() {
    if (this.args.model.validate) {
      const { isValid, state, invalidFormMessage } = this.args.model.validate();
      this.modelValidations = state;
      this.invalidFormAlert = invalidFormMessage;
      return isValid;
    }
    return true;
  }

  @action
  async generateRoot(event) {
    event.preventDefault();
    const continueSave = this.checkFormValidity();
    if (!continueSave) return;
    try {
      await this.args.model.save({ adapterOptions: this.args.adapterOptions });
      this.flashMessages.success('Successfully generated root.');
      this.args.onSave();
    } catch (e) {
      this.errorBanner = errorMessage(e);
      this.invalidFormAlert = 'There was a problem generating the root.';
    }
  }
}
