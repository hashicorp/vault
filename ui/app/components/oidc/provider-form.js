import Component from '@glimmer/component';
import { action } from '@ember/object';
import { inject as service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { task } from 'ember-concurrency';

/**
 * @module OidcProviderForm
 * OidcProviderForm components are used to create and update OIDC providers
 *
 * @example
 * ```js
 * <OidcProviderForm @model={{this.model}} />
 * ```
 * @callback onCancel
 * @callback onSave
 * @param {Object} model - oidc client model
 * @param {onCancel} onCancel - callback triggered when cancel button is clicked
 * @param {onSave} onSave - callback triggered on save success
 */

export default class OidcProviderForm extends Component {
  @service store;
  @service flashMessages;

  @tracked modelValidations;
  @tracked radioCardGroupValue =
    !this.args.model.allowedClientIds || this.args.model.allowedClientIds.includes('*')
      ? 'allow_all'
      : 'limited';

  @action
  handleClientSelection(selection) {
    // if array then coming from search-select component, set selection as model clients
    if (Array.isArray(selection)) {
      this.args.model.allowedClientIds = selection;
    } else {
      // otherwise update radio button value and reset clients so
      // UI always reflects a user's selection (including when no clients are selected)
      this.radioCardGroupValue = selection;
      this.args.model.allowedClientIds = [];
    }
  }

  @task
  *save(event) {
    event.preventDefault();
    try {
      const { isValid, state } = this.args.model.validate();
      this.modelValidations = isValid ? null : state;
      if (isValid) {
        if (this.radioCardGroupValue === 'allow_all') {
          this.args.model.allowedClientIds = ['*'];
        }
        const { isNew, name } = this.args.model;
        yield this.args.model.save();
        this.flashMessages.success(
          `Successfully ${isNew ? 'created' : 'updated'} the OIDC provider 
          ${name}.`
        );
        this.args.onSave();
      }
    } catch (error) {
      const message = error.errors ? error.errors.join('. ') : error.message;
      this.flashMessages.danger(message);
    }
  }

  @action
  cancel() {
    const method = this.args.model.isNew ? 'unloadRecord' : 'rollbackAttributes';
    this.args.model[method]();
    this.args.onCancel();
  }
}
