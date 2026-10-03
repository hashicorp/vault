/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { service } from '@ember/service';
import errorMessage from 'vault/utils/error-message';
import { tracked } from '@glimmer/tracking';

import type { LdapLibraryRouteModel } from 'ldap/routes/libraries/library';
import { Breadcrumb } from 'vault/vault/app-types';
import type FlashMessageService from 'vault/services/flash-messages';
import type RouterService from '@ember/routing/router-service';
import type SecretMountPath from 'vault/services/secret-mount-path';
import type ApiService from 'vault/services/api';

interface Args {
  model: LdapLibraryRouteModel;
  breadcrumbs: Array<Breadcrumb>;
}

export default class LdapLibraryDetailsPageComponent extends Component<Args> {
  @service declare readonly flashMessages: FlashMessageService;
  @service('app-router') declare readonly router: RouterService;
  @service declare readonly secretMountPath: SecretMountPath;
  @service declare readonly api: ApiService;

  @tracked showConfirmModal = false;

  @action
  async delete() {
    try {
      const { completeLibraryName } = this.args.model.library;
      await this.api.secrets.ldapLibraryDelete(completeLibraryName, this.secretMountPath.currentPath);
      this.flashMessages.success('Library deleted successfully.');
      this.router.transitionTo('vault.cluster.secrets.backend.ldap.libraries');
    } catch (error) {
      const message = errorMessage(error, 'Unable to delete library. Please try again or contact support.');
      this.flashMessages.danger(message);
    }
  }
}
