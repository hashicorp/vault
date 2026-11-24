/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { task } from 'ember-concurrency';
import { waitFor } from '@ember/test-waiters';

/**
 * @module Oidc::AssignmentForm
 * Oidc::AssignmentForm components are used to display the create view for OIDC providers assignments.
 *
 * @example
 * ```js
 * <Oidc::AssignmentForm @form={this.form}
 * @onCancel={transition-to "vault.cluster.access.oidc.assignment"} @param1={{param1}}
 * @onSave={transition-to "vault.cluster.access.oidc.assignments.assignment.details" this.model.name}
 * />
 * ```

 * @param {object} form - Form class
 * @callback onCancel - callback triggered when cancel button is clicked
 * @callback onSave - callback triggered when save button is clicked*
 */

export default class OidcAssignmentFormComponent extends Component {
  @service api;
  @service flashMessages;

  @tracked modelValidations;
  @tracked invalidFormMessage;
  @tracked errorBanner;

  save = task(
    waitFor(async (event) => {
      event.preventDefault();
      try {
        const { isValid, state, invalidFormMessage, data } = this.args.form.toJSON();
        this.modelValidations = isValid ? null : state;
        this.invalidFormMessage = invalidFormMessage;

        if (isValid) {
          const { name, ...payload } = data;
          await this.api.identity.oidcWriteAssignment(name, payload);
          this.flashMessages.success(
            `Successfully ${this.args.form.isNew ? 'created' : 'updated'} the assignment ${name}.`
          );
          // this form is sometimes used in modal, passing the model notifies
          // the parent if the save was successful
          this.args.onSave(this.args.form);
        }
      } catch (error) {
        const { message } = await this.api.parseError(error);
        this.errorBanner = message;
      }
    })
  );
}
