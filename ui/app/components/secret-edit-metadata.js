import Component from '@glimmer/component';
import { action } from '@ember/object';
import { inject as service } from '@ember/service';
import { tracked } from '@glimmer/tracking';

/**
 * @module SecretEditMetadata
 *
 * @example
 * ```js
 * <SecretEditMetadata
 * @model={{model}}
 * @mode={{mode}}
 * @updateValidationErrorCount={{updateValidationErrorCount}}
 * />
 * ```
 *
 * @param {object} model - name of the current cluster, passed from the parent.
 * @param {string} mode - if the mode is create, show, edit.
 * @param {Function} [updateValidationErrorCount] - function on parent that handles disabling the save button.
 */

export default class SecretEditMetadata extends Component {
  @service router;
  @service store;

  @tracked validationErrorCount = 0;
  @tracked modelValidations;

  async save() {
    const model = this.args.model;
    try {
      await model.save();
    } catch (e) {
      this.error = e;
      return;
    }
    this.router.transitionTo('vault.cluster.secrets.backend.metadata', this.args.model.id);
  }

  @action
  onSaveChanges(event) {
    event.preventDefault();
    return this.save();
  }
  @action onKeyUp(name, value) {
    let state = {};
    if (value) {
      if (name === 'customMetadata') {
        // atypical case where property is not set on model on change - validate independently
        /* eslint-disable no-useless-escape */
        const regex = /^[^\\]+$/g; // looking for a backward slash
        if (!value.match(regex)) {
          state[name] = {
            errors: ['Custom values cannot contain a backward slash.'],
            isValid: false,
          };
        }
      }
      if (name === 'maxVersions') {
        this.args.model.maxVersions = value;
        state = this.args.model.validate().state;
      }
    }
    let count = 0;
    for (const key in state) {
      if (!state[key].isValid) {
        count++;
      }
    }
    this.modelValidations = state;
    this.validationErrorCount = count;
    // when mode is "update" this works, but on mode "create" we need to bubble up the count
    if (this.args.updateValidationErrorCount) {
      this.args.updateValidationErrorCount(this.validationErrorCount);
    }
  }
}
