import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
import { task } from 'ember-concurrency';
import { inject as service } from '@ember/service';

/**
 * @module OidcScopeForm
 * Oidc scope form components are used to create and edit oidc scopes
 *
 * @example
 * ```js
 * <Oidc::ScopeForm @model={{this.model}} />
 * ```
 * @callback onCancel
 * @callback onSave
 * @param {Object} model - oidc scope model
 * @param {onCancel} onCancel - callback triggered when cancel button is clicked
 * @param {onSave} onSave - callback triggered on save success
 */

export default class OidcScopeFormComponent extends Component {
  @service flashMessages;

  @tracked modelValidations;

  exampleTemplate = `{
  "username": {{identity.entity.aliases.$MOUNT_ACCESSOR.name}},
  "contact": {
    "email": {{identity.entity.metadata.email}},
    "phone_number": {{identity.entity.metadata.phone_number}}
  },
  "groups": {{identity.entity.groups.names}}
}`;

  @task
  *save(event) {
    event.preventDefault();
    try {
      const { isValid, state } = this.args.model.validate();
      this.modelValidations = isValid ? null : state;
      if (isValid) {
        yield this.args.model.save();
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
