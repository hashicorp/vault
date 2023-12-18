/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { inject as service } from '@ember/service';
import { action } from '@ember/object';

import type LdapLibraryModel from 'vault/models/ldap/library';
import type SecretEngineModel from 'vault/models/secret-engine';
import type RouterService from '@ember/routing/router-service';
import type { Breadcrumb } from 'vault/vault/app-types';
import LdapRoleModel from 'vault/models/ldap/role';
import { LdapLibraryAccountStatus } from 'vault/vault/adapters/ldap/library';

interface Args {
  roles: Array<LdapRoleModel>;
  libraries: Array<LdapLibraryModel>;
  librariesStatus: Array<LdapLibraryAccountStatus>;
  promptConfig: boolean;
  backendModel: SecretEngineModel;
  breadcrumbs: Array<Breadcrumb>;
}

export default class LdapLibrariesPageComponent extends Component<Args> {
  @service declare readonly router: RouterService;

  @tracked selectedRole: LdapRoleModel | undefined;

  @action
  selectRole([roleName]: Array<string>) {
    const model = this.args.roles.find((role) => role.name === roleName);
    this.selectedRole = model;
  }

  @action
  generateCredentials() {
    const { type, name } = this.selectedRole as LdapRoleModel;
    this.router.transitionTo('vault.cluster.secrets.backend.ldap.roles.role.credentials', type, name);
  }
}
