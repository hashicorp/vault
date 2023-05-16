/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Component from '@glimmer/component';
import { inject as service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
import { task } from 'ember-concurrency';
import { waitFor } from '@ember/test-waiters';
// import errorMessage from 'vault/utils/error-message';

/**
 * @module CreateAndEditSecretPage
 * CreateAndEditRolePage component is a child component for create and edit secret pages.
 *
 * @param {object} model - secret model that contains secret record and backend
 */

export default class CreateAndEditSecretPageComponent extends Component {
  @service router;
  @service flashMessages;

  // @tracked roleRulesTemplates;
  @tracked selectedTemplateId;
  @tracked modelValidations;
  @tracked invalidFormAlert;
  @tracked errorBanner;

  constructor() {
    super(...arguments);
    // things may go here
  }

  @task
  @waitFor
  *save() {
    // try {
    //   // set generatedRoleRoles to value of selected template
    //   const selectedTemplate = this.roleRulesTemplates?.findBy('id', this.selectedTemplateId);
    //   if (selectedTemplate) {
    //     this.args.model.generatedRoleRules = selectedTemplate.rules;
    //   }
    //   yield this.args.model.save();
    //   this.router.transitionTo(
    //     'vault.cluster.secrets.backend.kubernetes.roles.role.details',
    //     this.args.model.name
    //   );
    // } catch (error) {
    //   const message = errorMessage(error, 'Error saving role. Please try again or contact support');
    //   this.errorBanner = message;
    //   this.invalidFormAlert = 'There was an error submitting this form.';
    // }
  }

  @action
  async onSave(event) {
    event.preventDefault();
    // const { isValid, state, invalidFormMessage } = await this.args.model.validate();
    // if (isValid) {
    //   this.modelValidations = null;
    //   this.save.perform();
    // } else {
    //   this.invalidFormAlert = invalidFormMessage;
    //   this.modelValidations = state;
    // }
  }

  @action
  cancel() {
    // const { model } = this.args;
    // const method = model.isNew ? 'unloadRecord' : 'rollbackAttributes';
    // model[method]();
    // this.router.transitionTo('vault.cluster.secrets.backend.kubernetes.roles');
  }
}
