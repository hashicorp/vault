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

import type { LdapConfigureModel } from 'ldap/routes/configure';
import type { Breadcrumb, ValidationMap } from 'vault/vault/app-types';
import type FlashMessageService from 'vault/services/flash-messages';
import type RouterService from '@ember/routing/router-service';
import type ApiService from 'vault/services/api';
import type SecretMountPath from 'vault/services/secret-mount-path';

interface Args {
  model: LdapConfigureModel;
  breadcrumbs: Array<Breadcrumb>;
}
interface SchemaOption {
  title: string;
  icon: string;
  description: string;
  value: string;
}

export default class LdapConfigurePageComponent extends Component<Args> {
  @service declare readonly flashMessages: FlashMessageService;
  @service('app-router') declare readonly router: RouterService;
  @service declare readonly api: ApiService;
  @service declare readonly secretMountPath: SecretMountPath;

  @tracked showRotatePrompt = false;
  @tracked modelValidations: ValidationMap | null = null;
  @tracked invalidFormMessage = '';
  @tracked error = '';

  get schemaOptions(): Array<SchemaOption> {
    return [
      {
        title: 'OpenLDAP',
        icon: 'folder',
        description:
          'OpenLDAP is one of the most popular open source directory service developed by the OpenLDAP Project.',
        value: 'openldap',
      },
      {
        title: 'AD',
        icon: 'microsoft',
        description:
          'Active Directory is a directory service developed by Microsoft for Windows domain networks.',
        value: 'ad',
      },
      {
        title: 'RACF',
        icon: 'users',
        description:
          "For managing IBM's Resource Access Control Facility (RACF) security system, the generated passwords must be 8 characters or less.",
        value: 'racf',
      },
    ];
  }

  leave(route: string) {
    this.router.transitionTo(`vault.cluster.secrets.backend.ldap.${route}`);
  }

  async rotateRoot() {
    try {
      await this.api.secrets.ldapRotateRootCredentials(this.secretMountPath.currentPath);
    } catch (error) {
      // since config save was successful at this point we only want to show the error in a flash message
      const { message } = await this.api.parseError(error);
      this.flashMessages.danger(`Error rotating root password \n ${message}`);
    }
  }

  async saveConfigModelAndRotateRoot(data: LdapConfigureModel['form']['data'], rotate: boolean) {
    try {
      await this.api.secrets.ldapConfigure(this.secretMountPath.currentPath, data);
      // if save was triggered from confirm action in rotate password prompt we need to make an additional request
      if (rotate) {
        await this.rotateRoot();
      }
      this.flashMessages.success('Successfully configured LDAP engine');
      this.leave('configuration');
    } catch (error) {
      const { message } = await this.api.parseError(
        error,
        'Error saving configuration. Please try again or contact support.'
      );
      this.error = message;
    }
  }

  save = task(
    waitFor(async (event: Event | null, rotate: boolean) => {
      if (event) {
        event.preventDefault();
      }
      const { form } = this.args.model;
      const { isValid, state, invalidFormMessage, data } = form.toJSON();

      this.modelValidations = isValid ? null : state;
      this.invalidFormMessage = isValid ? '' : invalidFormMessage;
      // show rotate creds prompt for new models when form state is valid
      this.showRotatePrompt = isValid && form.isNew && !this.showRotatePrompt;

      if (isValid && !this.showRotatePrompt) {
        await this.saveConfigModelAndRotateRoot(data, rotate);
      }
    })
  );

  @action
  cancel() {
    const { isNew } = this.args.model.form;
    const transitionRoute = isNew ? 'overview' : 'configuration';
    this.leave(transitionRoute);
  }

  saveAndClose = task(
    waitFor(async (rotate: boolean, close: () => void) => {
      close();
      const { data } = this.args.model.form.toJSON();
      await this.saveConfigModelAndRotateRoot(data, rotate);
    })
  );
}
