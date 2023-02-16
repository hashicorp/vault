import Component from '@glimmer/component';
import { inject as service } from '@ember/service';
import { task } from 'ember-concurrency';
import { tracked } from '@glimmer/tracking';
import errorMessage from 'vault/utils/error-message';
import { waitFor } from '@ember/test-waiters';

/**
 * @module PkiKeyForm
 * PkiKeyForm components are used to create and update PKI keys.
 *
 * @example
 * ```js
 * <PkiKeyForm @model={{this.model}} @onCancel={{transition-to "vault.cluster"}} @onSave={{transition-to "vault.cluster"}} />
 * ```
 *
 * @param {Object} model - pki/key model.
 * @callback onCancel - Callback triggered when cancel button is clicked.
 * @callback onSave - Callback triggered on save success.
 */

export default class PkiKeyForm extends Component {
  @service flashMessages;

  @tracked errorBanner;
  @tracked invalidFormAlert;
  @tracked modelValidations;

  @task
  @waitFor
  *save(event) {
    event.preventDefault();
    try {
      const { isNew, keyName } = this.args.model;
      const { isValid, state, invalidFormMessage } = this.args.model.validate();
      if (isNew) {
        this.modelValidations = isValid ? null : state;
        this.invalidFormAlert = invalidFormMessage;
      }
      if (!isValid && isNew) return;
      yield this.args.model.save({ adapterOptions: { import: false } });
      this.flashMessages.success(
        `Successfully ${isNew ? 'generated' : 'updated'} key${keyName ? ` ${keyName}.` : '.'}`
      );
      this.args.onSave();
    } catch (error) {
      this.errorBanner = errorMessage(error);
      this.invalidFormAlert = 'There was an error submitting this form.';
    }
  }
}
