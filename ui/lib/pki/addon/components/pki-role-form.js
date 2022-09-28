import Component from '@glimmer/component';
import { action } from '@ember/object';
import { inject as service } from '@ember/service';
import { task } from 'ember-concurrency';
import { tracked } from '@glimmer/tracking';

/**
 * @module PkiRoleForm
 * PkiRoleForm components are used to create and update PKI roles.
 *
 * @example
 * ```js
 * <PkiRoleForm @model={{this.model}}/>
 * ```
 * @callback onCancel
 * @callback onSave
 * @param {Object} model - Pki-role-engine model.
 * @param {function} onChange - Handle the input coming from the custom yield field.
 * @param {onCancel} onCancel - Callback triggered when cancel button is clicked.
 * @param {onSave} onSave - Callback triggered on save success.
 */

export default class PkiRoleForm extends Component {
  @service store;
  @service flashMessages;

  @tracked errorBanner;
  @tracked invalidFormAlert;
  @tracked keyBits;
  @tracked modelValidations;
  @tracked notAfter;
  @tracked ttl;

  defaultKeyBits() {
    if (this.args.model.keyType === 'rsa') return 2048;
    if (this.args.model.keyType === 'ec') return 256;
    return 0;
  }

  @task
  *save(event) {
    event.preventDefault();
    try {
      const { isValid, state, invalidFormMessage } = this.args.model.validate();
      this.modelValidations = isValid ? null : state;
      this.invalidFormAlert = invalidFormMessage;
      if (isValid) {
        const { isNew, name } = this.args.model;
        // custom set params
        this.args.model.notAfter = this.notAfter;
        this.args.model.ttl = this.ttl;
        // user hasn't triggered the select by clicking keyBits so need to set a default
        if (!this.args.model.keyBits) {
          this.args.model.keyBits = this.defaultKeyBits();
        }
        yield this.args.model.save();
        this.flashMessages.success(`Successfully ${isNew ? 'created' : 'updated'} the role ${name}.`);
        this.args.onSave();
      }
    } catch (error) {
      const message = error.errors ? error.errors.join('. ') : error.message;
      this.errorBanner = message;
      this.invalidFormAlert = 'There was an error submitting this form.';
    }
  }

  @action onNotValidAfterChange(modelParam, value) {
    if (modelParam === 'ttl') {
      this.notAfter = '';
      this.ttl = value;
    }
    if (modelParam === 'not_after') {
      this.ttl = '';
      this.notAfter = value;
    }
  }

  @action
  cancel() {
    const method = this.args.model.isNew ? 'unloadRecord' : 'rollbackAttributes';
    this.args.model[method]();
    this.args.onCancel();
  }
}
