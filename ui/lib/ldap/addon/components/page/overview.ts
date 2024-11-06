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
  libraries: Array<LdapLibraryModel>;
  librariesStatus: Array<LdapLibraryAccountStatus>;
  promptConfig: boolean;
  backendModel: SecretEngineModel;
  breadcrumbs: Array<Breadcrumb>;
}

interface Option {
  id: string;
  name: string;
  type: string;
}

export default class LdapLibrariesPageComponent extends Component<Args> {
  @service('app-router') declare readonly router: RouterService;

  @tracked selectedRole: LdapRoleModel | undefined;

  get roleOptions() {
    const options = this.args.roles
      // hierarchical roles are not selectable
      .filter((r: LdapRoleModel) => !r.name.endsWith('/'))
      // *hack alert* - type is set as id so it renders beside name in search select
      // this is to avoid more changes to search select and is okay here because
      // we use the type and name to select the item below, not the id
      .map((r: LdapRoleModel) => ({ id: r.type, name: r.name, type: r.type }));
    return options;
  }

  @action
  async selectRole([option]: Array<Option>) {
    if (option) {
      const { name, type } = option;
      const model = this.args.roles.find((role) => role.name === name && role.type === type);
      this.selectedRole = model;
    }
  }

  @action
  generateCredentials() {
    const { type, name } = this.selectedRole as LdapRoleModel;
    this.router.transitionTo('vault.cluster.secrets.backend.ldap.roles.role.credentials', type, name);
  }
}
