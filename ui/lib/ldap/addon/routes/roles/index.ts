/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import LdapRolesRoute from '../roles';

import type Transition from '@ember/routing/transition';
import type { LdapRole } from 'vault/secrets/ldap';
import type Controller from '@ember/controller';
import type { Breadcrumb } from 'vault/app-types';
import type SecretsEngineResource from 'vault/resources/secrets/engine';
import { LdapApplicationModel } from '../application';

interface RouteModel {
  secretsEngine: SecretsEngineResource;
  promptConfig: boolean;
  roles: Array<LdapRole>;
}

interface RouteController extends Controller {
  breadcrumbs: Array<Breadcrumb>;
  model: RouteModel;
}

export default class LdapRolesIndexRoute extends LdapRolesRoute {
  queryParams = {
    pageFilter: {
      refreshModel: true,
    },
    page: {
      refreshModel: true,
    },
  };

  async model(params: { page?: string; pageFilter: string }) {
    const { page, pageFilter: filter } = params;
    const { secretsEngine, promptConfig } = this.modelFor('application') as LdapApplicationModel;
    const { roles, capabilities } = await this.fetchRolesAndCapabilities({ page: Number(page) || 1, filter });
    return {
      secretsEngine,
      promptConfig,
      roles,
      capabilities,
    };
  }

  setupController(controller: RouteController, resolvedModel: RouteModel, transition: Transition) {
    super.setupController(controller, resolvedModel, transition);

    controller.breadcrumbs = [
      { label: 'Secrets', route: 'secrets', linkExternal: true },
      { label: resolvedModel.secretsEngine.id, route: 'overview' },
      { label: 'Roles' },
    ];
  }

  resetController(controller: RouteController, isExiting: boolean) {
    if (isExiting) {
      controller.set('pageFilter', undefined);
      controller.set('page', undefined);
    }
  }
}
