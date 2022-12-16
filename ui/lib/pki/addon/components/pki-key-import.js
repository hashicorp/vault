import { action } from '@ember/object';
import Component from '@glimmer/component';
import { task } from 'ember-concurrency';
import { tracked } from '@glimmer/tracking';
import { inject as service } from '@ember/service';
import errorMessage from 'vault/utils/error-message';

// TODO: convert to typescript after https://github.com/hashicorp/vault/pull/18387 is merged
/**
 * @module PkiKeyImport
 * PkiKeyImport components are used to import PKI keys.
 *
 * @example
 * ```js
 * <PkiKeyImport @model={{this.model}} />
 * ```
 *
 * @param {Object} model - pki/key model.
 * @callback onCancel - Callback triggered when cancel button is clicked.
 * @callback onSubmit - Callback triggered on submit success.
 */

export default class PkiKeyImport extends Component {
  @service store;
  @service flashMessages;

  @tracked file = { value: '', enterAsText: false };
  @tracked errorBanner;
  @tracked invalidFormAlert;
  @tracked modelValidations;

  @task
  *submitForm(event) {
    event.preventDefault();
    try {
      const { isValid, state, invalidFormMessage } = this.args.model.validate();
      this.modelValidations = isValid ? null : state;
      this.invalidFormAlert = invalidFormMessage;
      if (isValid) {
        const { keyName } = this.args.model;
        yield this.args.model.save();
        this.flashMessages.success(`Successfully imported key ${keyName}.`);
        this.args.onSave();
      }
    } catch (error) {
      this.errorBanner = errorMessage(error);
      this.invalidFormAlert = 'There was an error submitting this form.';
    }
  }

  @action
  onFileUpload() {
    // do something
  }

  @action
  cancel() {
    const method = this.args.model.isNew ? 'unloadRecord' : 'rollbackAttributes';
    this.args.model[method]();
    this.args.onCancel();
  }
}
