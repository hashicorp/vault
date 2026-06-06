/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { ModelFrom } from 'vault/route';
import {
  SecretsApiLdapListStaticRolesListEnum,
  SecretsApiLdapListDynamicRolesListEnum,
} from '@hashicorp/vault-client-typescript';

import type SecretMountPath from 'vault/services/secret-mount-path';
import type Transition from '@ember/routing/transition';
import type Controller from '@ember/controller';
import type { Breadcrumb } from 'vault/app-types';
import type { LdapApplicationModel } from './application';
import type ApiService from 'vault/services/api';

export type LdapOverviewRouteModel = ModelFrom<LdapOverviewRoute>;

interface RouteController extends Controller {
  breadcrumbs: Array<Breadcrumb>;
  model: LdapOverviewRouteModel;
}

export default class LdapOverviewRoute extends Route {
  @service declare readonly api: ApiService;
  @service declare readonly secretMountPath: SecretMountPath;

  async model() {
    const { promptConfig, secretsEngine } = this.modelFor('application') as LdapApplicationModel;
    const { currentPath } = this.secretMountPath;
    const requests = [
      this.api.secrets.ldapListStaticRoles(currentPath, SecretsApiLdapListStaticRolesListEnum.TRUE),
      this.api.secrets.ldapListDynamicRoles(currentPath, SecretsApiLdapListDynamicRolesListEnum.TRUE),
    ];
    const results = await Promise.allSettled(requests);
    const roles = [];
    for (const result of results) {
      if (result.status === 'fulfilled') {
        if (result.value.keys) {
          const type = results.indexOf(result) === 0 ? 'static' : 'dynamic';
          roles.push(...result.value.keys.map((name) => ({ name, type, completeRoleName: name })));
        }
      }
    }
    return {
      promptConfig,
      secretsEngine,
      roles,
    };
  }

  setupController(
    controller: RouteController,
    resolvedModel: LdapOverviewRouteModel,
    transition: Transition
  ) {
    super.setupController(controller, resolvedModel, transition);

    controller.breadcrumbs = [
      { label: 'Vault', route: 'vault', icon: 'vault', linkExternal: true },
      { label: 'Secrets engines', route: 'secrets', linkExternal: true },
      { label: this.secretMountPath.currentPath },
    ];
  }
}
