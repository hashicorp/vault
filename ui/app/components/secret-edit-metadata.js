/**
 * @module SecretEditMetadata
 *
 * @example
 * ```js
 * <SecretEditMetadata
 * @model={{model}}
 * @validationMessages={{validationMessages}}
 * @mode={{mode}}
 * />
 * ```
 *
 * @param {object} model - name of the current cluster, passed from the parent.
 * @param {object} [validationMessages] - Object that contains form validation errors. keys are the field names and values are the messages.
 * @param {Function} mode - if the mode is create, show, edit.
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { inject as service } from '@ember/service';
import { set } from '@ember/object';
import { action } from '@ember/object';
export default class SecretEditMetadata extends Component {
  @service router;
  @service store;

  @tracked validationMesssages;

  constructor() {
    super(...arguments);
    this.validationMessages = {
      customMetadata: '',
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
      // // ARG TODO for now set this to hardcoded.
      // this.model.set('customMetadata', { key: 'meep', value: value });
      // cp validations won't work on an object so performing validations here
      let regex = /^[^\/]+$/g; // looking for a forward slash
      if (!value.match(regex)) {
        set(this.validationMessages, name, 'Custom values cannot contain a forward slash.');
      } else {
        set(this.validationMessages, name, '');
      }
    }
  }
}
