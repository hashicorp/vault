/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { action } from '@ember/object';
import errorMessage from 'vault/utils/error-message';

import type FlashMessageService from 'vault/services/flash-messages';
import type RouterService from '@ember/routing/router-service';
import type LdapLibraryModel from 'vault/models/ldap/library';
import type Controller from '@ember/controller';
import type Transition from '@ember/routing/transition';
import type { Breadcrumb } from 'vault/vault/app-types';
import { LdapLibraryCheckOutCredentials } from 'vault/vault/adapters/ldap/library';
import type AdapterError from 'ember-data/adapter'; // eslint-disable-line ember/use-ember-data-rfc-395-imports
import { ldapBreadcrumbs, libraryRoutes } from 'ldap/utils/ldap-breadcrumbs';

interface LdapLibraryCheckOutController extends Controller {
  breadcrumbs: Array<Breadcrumb>;
  model: LdapLibraryCheckOutCredentials;
}

export default class LdapLibraryCheckOutRoute extends Route {
  @service declare readonly flashMessages: FlashMessageService;
  @service('app-router') declare readonly router: RouterService;

  accountsRoute = 'vault.cluster.secrets.backend.ldap.libraries.library.details.accounts';

  beforeModel(transition: Transition) {
    // transition must be from the details.accounts route to ensure it was initiated by the check-out action
    if (transition.from?.name !== this.accountsRoute) {
      this.router.replaceWith(this.accountsRoute);
    }
  }
  model(_params: object, transition: Transition) {
    const ttl = transition.to?.queryParams['ttl'];
    const library = this.modelFor('libraries.library') as LdapLibraryModel;
    return library.checkOutAccount(ttl);
  }
  setupController(
    controller: LdapLibraryCheckOutController,
    resolvedModel: LdapLibraryCheckOutCredentials,
    transition: Transition
  ) {
    super.setupController(controller, resolvedModel, transition);
    const library = this.modelFor('libraries.library') as LdapLibraryModel;
    const routeParams = (childResource: string) => {
      return [library.backend, childResource];
    };
    controller.breadcrumbs = [
      { label: library.backend, route: 'overview' },
      { label: 'Libraries', route: 'libraries' },
      ...ldapBreadcrumbs(library.name, routeParams, libraryRoutes),
      { label: 'Check-Out' },
    ];
  }

  @action
  error(error: AdapterError) {
    // if check-out fails, return to library details route
    const message = errorMessage(error, 'Error checking out account. Please try again or contact support.');
    this.flashMessages.danger(message);
    this.router.replaceWith(this.accountsRoute);
  }
}
