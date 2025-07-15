/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { task } from 'ember-concurrency';
import { waitFor } from '@ember/test-waiters';
import { service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { ValidationMap } from 'vault/vault/app-types';

import type SshConfigForm from 'vault/forms/secrets/ssh-config';
import type Router from '@ember/routing/router';
import type FlashMessageService from 'vault/services/flash-messages';
import type ApiService from 'vault/services/api';

/**
 * @module ConfigureSshComponent is used to configure the SSH secret engine.
 *
 * @example
 * ```js
 * <SecretEngine::ConfigureSsh
 *    @configForm={{this.model.configForm}}
 *    @id={{this.model.id}}
 *  />
 * ```
 *
 * @param {string} configForm - SSH ca-config form
 * @param {string} id - name of the SSH secret engine, ex: 'ssh-123'
 */

interface Args {
  configForm: SshConfigForm;
  id: string;
}

export default class ConfigureSshComponent extends Component<Args> {
  @service declare readonly router: Router;
  @service declare readonly api: ApiService;
  @service declare readonly flashMessages: FlashMessageService;

  @tracked errorMessage: string | null = null;
  @tracked invalidFormAlert: string | null = null;
  @tracked modelValidations: ValidationMap | null = null;

  save = task(
    waitFor(async (event: Event) => {
      event.preventDefault();
      this.resetErrors();

      const { id, configForm } = this.args;
      const { isValid, state, invalidFormMessage, data } = configForm.toJSON();
      this.modelValidations = isValid ? null : state;
      this.invalidFormAlert = isValid ? '' : invalidFormMessage;

      if (isValid) {
        try {
          await this.api.secrets.sshConfigureCa(id, data);
          this.flashMessages.success(`Successfully saved ${id}'s root configuration.`);
          this.transition();
        } catch (error) {
          const { message } = await this.api.parseError(error);
          this.errorMessage = message;
          this.invalidFormAlert = 'There was an error submitting this form.';
        }
      }
    })
  );

  resetErrors() {
    this.flashMessages.clearMessages();
    this.errorMessage = null;
    this.invalidFormAlert = null;
  }

  transition(isDelete = false) {
    // deleting a key is the only case in which we want to stay on the create/edit page.
    const { id } = this.args;
    if (isDelete) {
      this.router.transitionTo('vault.cluster.secrets.backend.configuration.edit', id);
    } else {
      this.router.transitionTo('vault.cluster.secrets.backend.configuration', id);
    }
  }

  @action
  onCancel() {
    // clear errors because they're canceling out of the workflow.
    this.resetErrors();
    this.transition();
  }

  @action
  async deleteCaConfig() {
    try {
      await this.api.secrets.sshDeleteCaConfiguration(this.args.id);
      this.flashMessages.success('CA information deleted successfully.');
      this.transition(true);
    } catch (error) {
      const { message } = await this.api.parseError(error);
      this.flashMessages.danger(message);
    }
  }
}
