/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { inject as service } from '@ember/service';
import { action } from '@ember/object';
import { getOwner } from '@ember/application';
import errorMessage from 'vault/utils/error-message';

import type LdapRoleModel from 'vault/models/ldap/role';
import type SecretEngineModel from 'vault/models/secret-engine';
import type FlashMessageService from 'vault/services/flash-messages';
import type { Breadcrumb, EngineOwner } from 'vault/vault/app-types';

interface Args {
  roles: Array<LdapRoleModel>;
  promptConfig: boolean;
  backendModel: SecretEngineModel;
  breadcrumbs: Array<Breadcrumb>;
}

export default class LdapRolesPageComponent extends Component<Args> {
  @service declare readonly flashMessages: FlashMessageService;

  @tracked filterValue = '';

  get mountPoint(): string {
    const owner = getOwner(this) as EngineOwner;
    return owner.mountPoint;
  }

  get filteredRoles() {
    const { roles } = this.args;
    return this.filterValue
      ? roles.filter((role) => role.name.toLowerCase().includes(this.filterValue.toLowerCase()))
      : roles;
  }

  @action
  async onRotate(model: LdapRoleModel) {
    try {
      const message = `Successfully rotated credentials for ${model.name}.`;
      await model.rotateStaticPassword();
      this.flashMessages.success(message);
    } catch (error) {
      this.flashMessages.danger(`Error rotating credentials \n ${errorMessage(error)}`);
    }
  }

  @action
  async onDelete(model: LdapRoleModel) {
    try {
      const message = `Successfully deleted role ${model.name}.`;
      await model.destroyRecord();
      this.args.roles.removeObject(model);
      this.flashMessages.success(message);
    } catch (error) {
      this.flashMessages.danger(`Error deleting role \n ${errorMessage(error)}`);
    }
  }
}
