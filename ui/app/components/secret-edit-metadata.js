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
 * @param {Function} updateValidationErrorCount - function on parent that handle disabling the save button.
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { inject as service } from '@ember/service';
import { set } from '@ember/object';
import { tracked } from '@glimmer/tracking';
export default class SecretEditMetadata extends Component {
  @service router;
  @service store;

  @tracked validationErrorCount = 0;

  constructor() {
    super(...arguments);
    this.validationMessages = {
      customMetadata: '',
      maxVersions: '',
    };
  }

  async save() {
    let model = this.args.model;
    try {
      await model.save();
    } catch (err) {
      // error
    }
    this.router.transitionTo('vault.cluster.secrets.backend.metadata', this.args.model.id);
  }

  @action
  onSaveChanges(event) {
    event.preventDefault();
    const changed = this.args.model.hasDirtyAttributes; // ARG TODO when API done double check this is working
    if (changed) {
      this.save();
      return;
    }
    // ARG TODO else validation error?
  }
  @action onKeyUp(name, value) {
    if (value) {
      if (name === 'customMetadata') {
        // // ARG TODO for now set this to hardcoded. Need to write test coverage for this.
        // this.model.set('customMetadata', { key: 'meep', value: value });
        // cp validations won't work on an object so performing validations here
        /* eslint-disable no-useless-escape */
        let regex = /^[^\/]+$/g; // looking for a forward slash
        if (!value.match(regex)) {
          set(this.validationMessages, name, 'Custom values cannot contain a forward slash.');
        } else {
          set(this.validationMessages, name, '');
        }
      }
      if (name === 'maxVersions') {
        let number = Number(value);
        this.args.model.maxVersions = number;
        if (!this.args.model.validations.attrs.maxVersions.isValid) {
          set(this.validationMessages, name, this.args.model.validations.attrs.maxVersions.message);
        } else {
          set(this.validationMessages, name, '');
        }
      }
    }

    let values = Object.values(this.validationMessages);
    this.validationErrorCount = values.filter(Boolean).length;
    // when it's update this work, but on create we need to bubble up the count
    // so it disables the save in the secret-creat-or-update form
    this.args.updateValidationErrorCount(this.validationErrorCount);
  }
}
