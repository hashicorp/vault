/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { service } from '@ember/service';
import errorMessage from 'vault/utils/error-message';
import type SecretMountPath from 'vault/services/secret-mount-path';
import type FlashMessageService from 'vault/services/flash-messages';
import type RouterService from '@ember/routing/router-service';
import type PkiRoleModel from 'vault/models/pki/role';

interface Args {
  role: PkiRoleModel;
}

export default class DetailsPage extends Component<Args> {
  @service('app-router') declare readonly router: RouterService;
  @service declare readonly flashMessages: FlashMessageService;
  @service declare readonly secretMountPath: SecretMountPath;

  get breadcrumbs() {
    return [
      { label: 'Secrets', route: 'secrets', linkExternal: true },
      { label: this.secretMountPath.currentPath, route: 'overview' },
      { label: 'Roles', route: 'roles.index' },
      { label: this.args.role.id },
    ];
  }

  get arrayAttrs() {
    return ['keyUsage', 'extKeyUsage', 'extKeyUsageOids'];
  }

  @action
  async deleteRole() {
    try {
      await this.args.role.destroyRecord();
      this.flashMessages.success('Role deleted successfully');
      this.router.transitionTo('vault.cluster.secrets.backend.pki.roles.index');
    } catch (error) {
      this.args.role.rollbackAttributes();
      this.flashMessages.danger(errorMessage(error));
    }
  }
}
