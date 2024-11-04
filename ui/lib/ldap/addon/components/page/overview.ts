/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { service } from '@ember/service';
import { action } from '@ember/object';

import type LdapLibraryModel from 'vault/models/ldap/library';
import type SecretEngineModel from 'vault/models/secret-engine';
import type RouterService from '@ember/routing/router-service';
import type { Breadcrumb } from 'vault/vault/app-types';
import LdapRoleModel from 'vault/models/ldap/role';
import { LdapLibraryAccountStatus } from 'vault/vault/adapters/ldap/library';

interface Args {
  roles: Array<LdapRoleModel>;
  staticRoles: Array<LdapRoleModel>;
  dynamicRoles: Array<LdapRoleModel>;
  libraries: Array<LdapLibraryModel>;
  librariesStatus: Array<LdapLibraryAccountStatus>;
  promptConfig: boolean;
  backendModel: SecretEngineModel;
  breadcrumbs: Array<Breadcrumb>;
}

interface Option {
  type: string;
  id: string;
}

export default class LdapLibrariesPageComponent extends Component<Args> {
  @service('app-router') declare readonly router: RouterService;

  @tracked selectedRole: Option | undefined;

  get roleOptions() {
    const mapOptions = (roles: LdapRoleModel[]) =>
      roles
        .filter((r: LdapRoleModel) => !r.name.endsWith('/'))
        .map((r: LdapRoleModel) => {
          if (r.name.endsWith('/')) return;
          return { id: r.name, type: r.type };
        });
    return [
      { groupName: 'Static', options: mapOptions(this.args.staticRoles) },
      { groupName: 'Dynamic', options: mapOptions(this.args.dynamicRoles) },
    ];
  }

  @action
  selectRole([role]: Array<Option>) {
    this.selectedRole = role;
  }

  @action
  generateCredentials() {
    const { type, id: name } = this.selectedRole as Option;
    this.router.transitionTo('vault.cluster.secrets.backend.ldap.roles.role.credentials', type, name);
  }
}
