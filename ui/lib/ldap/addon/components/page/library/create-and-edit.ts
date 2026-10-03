/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
import { service } from '@ember/service';
import { task } from 'ember-concurrency';
import { waitFor } from '@ember/test-waiters';

import type LdapLibraryForm from 'vault/forms/secrets/ldap/library';
import type { Breadcrumb, ValidationMap } from 'vault/app-types';
import type FlashMessageService from 'vault/services/flash-messages';
import type RouterService from '@ember/routing/router-service';
import type ApiService from 'vault/services/api';
import type SecretMountPath from 'vault/services/secret-mount-path';

interface Args {
  form: LdapLibraryForm;
  breadcrumbs: Array<Breadcrumb>;
}

export default class LdapCreateAndEditLibraryPageComponent extends Component<Args> {
  @service declare readonly flashMessages: FlashMessageService;
  @service('app-router') declare readonly router: RouterService;
  @service declare readonly api: ApiService;
  @service declare readonly secretMountPath: SecretMountPath;

  @tracked modelValidations: ValidationMap | null = null;
  @tracked invalidFormMessage = '';
  @tracked error = '';

  save = task(
    waitFor(async (event: Event) => {
      event.preventDefault();

      const { currentPath } = this.secretMountPath;
      const { form } = this.args;
      const { isValid, state, invalidFormMessage, data } = form.toJSON();

      this.modelValidations = isValid ? null : state;
      this.invalidFormMessage = isValid ? '' : invalidFormMessage;

      if (isValid) {
        try {
          const action = form.isNew ? 'created' : 'updated';
          const { name, ...rest } = data;
          // transform disable_check_in_enforcement back to boolean
          const disable_check_in_enforcement = data.disable_check_in_enforcement === 'Enabled' ? false : true;
          const payload = { ...rest, disable_check_in_enforcement };
          await this.api.secrets.ldapLibraryConfigure(name, currentPath, payload);
          this.flashMessages.success(`Successfully ${action} the library ${name}.`);
          const libraryParam = name.includes('/') ? encodeURIComponent(name) : name;
          this.router.transitionTo(
            'vault.cluster.secrets.backend.ldap.libraries.library.details',
            libraryParam
          );
        } catch (error) {
          const { message } = await this.api.parseError(
            error,
            'Error saving library. Please try again or contact support.'
          );
          this.error = message;
        }
      }
    })
  );

  @action
  cancel() {
    this.router.transitionTo('vault.cluster.secrets.backend.ldap.libraries');
  }
}
