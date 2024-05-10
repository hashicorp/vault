/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
import { task } from 'ember-concurrency';
import { service } from '@ember/service';

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
  @tracked errorBanner;
  @tracked invalidFormAlert;
  @tracked modelValidations;
  // formatting here is purposeful so that whitespace renders correctly in JsonEditor
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
      const { isValid, state, invalidFormMessage } = this.args.model.validate();
      this.modelValidations = isValid ? null : state;
      this.invalidFormAlert = invalidFormMessage;
      if (isValid) {
        const { isNew, name } = this.args.model;
        yield this.args.model.save();
        this.flashMessages.success(`Successfully ${isNew ? 'created' : 'updated'} the scope ${name}.`);
        this.args.onSave();
      }
    } catch (error) {
      const message = error.errors ? error.errors.join('. ') : error.message;
      this.errorBanner = message;
      this.invalidFormAlert = 'There was an error submitting this form.';
    }
  }
  @action
  cancel() {
    const method = this.args.model.isNew ? 'unloadRecord' : 'rollbackAttributes';
    this.args.model[method]();
    this.args.onCancel();
  }
}
