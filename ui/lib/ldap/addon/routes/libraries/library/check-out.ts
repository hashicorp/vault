/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { ModelFrom } from 'vault/route';

import type FlashMessageService from 'vault/services/flash-messages';
import type RouterService from '@ember/routing/router-service';
import type { LdapLibraryRouteModel } from 'ldap/routes/libraries/library';
import type Controller from '@ember/controller';
import type Transition from '@ember/routing/transition';
import type { Breadcrumb } from 'vault/app-types';
import { LdapLibraryCheckOutCredentials } from 'vault/vault/adapters/ldap/library';
import { ldapBreadcrumbs, libraryRoutes } from 'ldap/utils/ldap-breadcrumbs';
import type ApiService from 'vault/services/api';
import type SecretMountPath from 'vault/services/secret-mount-path';

export type LdapLibraryCheckOutRouteModel = ModelFrom<LdapLibraryCheckOutRoute>;

interface LdapLibraryCheckOutController extends Controller {
  breadcrumbs: Array<Breadcrumb>;
  model: LdapLibraryCheckOutRouteModel;
}

export default class LdapLibraryCheckOutRoute extends Route {
  @service declare readonly flashMessages: FlashMessageService;
  @service('app-router') declare readonly router: RouterService;
  @service declare readonly api: ApiService;
  @service declare readonly secretMountPath: SecretMountPath;

  accountsRoute = 'vault.cluster.secrets.backend.ldap.libraries.library.details.accounts';

  beforeModel(transition: Transition) {
    // transition must be from the details.accounts route to ensure it was initiated by the check-out action
    if (transition.from?.name !== this.accountsRoute) {
      this.router.replaceWith(this.accountsRoute);
    }
  }

  async model(_params: object, transition: Transition) {
    try {
      const ttl = transition.to?.queryParams['ttl'] as string;
      const { library } = this.modelFor('libraries.library') as LdapLibraryRouteModel;
      const { currentPath } = this.secretMountPath;

      const response = await this.api.secrets.ldapLibraryCheckOut(library.completeLibraryName, currentPath, {
        ttl,
      });
      const { lease_id, lease_duration, renewable } = response;
      const { service_account_name, password } = response.data as {
        service_account_name: string;
        password: string;
      };

      return { lease_id, lease_duration, renewable, service_account_name, password };
    } catch (error) {
      // if check-out fails, return to library details route
      const { message } = await this.api.parseError(
        error,
        'Error checking out account. Please try again or contact support.'
      );
      this.flashMessages.danger(message);
      return this.router.replaceWith(this.accountsRoute);
    }
  }
  setupController(
    controller: LdapLibraryCheckOutController,
    resolvedModel: LdapLibraryCheckOutCredentials,
    transition: Transition
  ) {
    super.setupController(controller, resolvedModel, transition);
    const { library } = this.modelFor('libraries.library') as LdapLibraryRouteModel;
    const { currentPath } = this.secretMountPath;
    const routeParams = (childResource: string) => {
      return [currentPath, childResource];
    };
    controller.breadcrumbs = [
      { label: 'Secrets', route: 'secrets', linkExternal: true },
      { label: currentPath, route: 'overview' },
      { label: 'Libraries', route: 'libraries' },
      ...ldapBreadcrumbs(library.completeLibraryName, routeParams, libraryRoutes),
      { label: 'Check-Out' },
    ];
  }
}
