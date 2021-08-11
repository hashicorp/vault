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
import { inject as service } from '@ember/service';
import { action } from '@ember/object';
export default class SecretEditMetadata extends Component {
  @service router;
  @service store;

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
}
