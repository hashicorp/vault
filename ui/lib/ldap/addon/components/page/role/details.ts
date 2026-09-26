/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
import { service } from '@ember/service';
import { task } from 'ember-concurrency';
import { waitFor } from '@ember/test-waiters';
import { toLabel } from 'core/helpers/to-label';

import type { LdapRolesRoleRouteModel } from 'ldap/routes/roles/role';
import type { Breadcrumb } from 'vault/vault/app-types';
import type FlashMessageService from 'vault/services/flash-messages';
import type RouterService from '@ember/routing/router-service';
import type ApiService from 'vault/services/api';
import type SecretMountPath from 'vault/services/secret-mount-path';

interface Args {
  model: LdapRolesRoleRouteModel;
  breadcrumbs: Array<Breadcrumb>;
}

export default class LdapRoleDetailsPageComponent extends Component<Args> {
  @service declare readonly flashMessages: FlashMessageService;
  @service('app-router') declare readonly router: RouterService;
  @service declare readonly api: ApiService;
  @service declare readonly secretMountPath: SecretMountPath;

  @tracked showConfirmDeleteModal = false;
  @tracked showConfirmRotateModal = false;
  @tracked cachedPassword: string | null = null;
  @tracked passwordError = '';

  isTtl = (field: string) => ['default_ttl', 'max_ttl', 'rotation_period'].includes(field);
  label = (field: string) => {
    return (
      {
        name: 'Role name',
        type: 'Role type',
        dn: 'Distinguished name',
        default_ttl: 'TTL',
        max_ttl: 'Max TTL',
        creation_ldif: 'Creation LDIF',
        deletion_ldif: 'Deletion LDIF',
        rollback_ldif: 'Rollback LDIF',
        password: 'Password',
      }[field] || toLabel([field])
    );
  };

  get displayFields() {
    const { role, capabilities } = this.args.model;
    const fields = ['name', 'type'];
    if (role.type === 'static') {
      fields.push('dn', 'username');
      if (capabilities?.canReadCreds) {
        fields.push('password');
      }
      fields.push('rotation_period');
    } else {
      fields.push(
        'default_ttl',
        'max_ttl',
        'username_template',
        'creation_ldif',
        'deletion_ldif',
        'rollback_ldif'
      );
    }
    return fields;
  }

  onPasswordToggle = task(
    waitFor(async () => {
      // Already fetched; MaskedInput owns its own show/hide state from here.
      if (this.cachedPassword !== null) {
        return;
      }

      this.passwordError = '';
      try {
        const { role } = this.args.model;
        const { currentPath } = this.secretMountPath;
        const response = await this.api.secrets.ldapRequestStaticRoleCredentials(role.name, currentPath);
        this.cachedPassword = (response.data as { password?: string }).password ?? '';
      } catch (error) {
        const { message } = await this.api.parseError(
          error,
          'Unable to retrieve password. Please try again or contact support.'
        );
        this.passwordError = message;
      }
    })
  );

  @action
  async delete() {
    try {
      const { role } = this.args.model;
      const { currentPath } = this.secretMountPath;

      if (role.type === 'static') {
        await this.api.secrets.ldapDeleteStaticRole(role.name, currentPath);
      } else {
        await this.api.secrets.ldapDeleteDynamicRole(role.name, currentPath);
      }
      this.flashMessages.success('Role deleted successfully.');
      this.router.transitionTo('vault.cluster.secrets.backend.ldap.roles');
    } catch (error) {
      const { message } = await this.api.parseError(
        error,
        'Unable to delete role. Please try again or contact support.'
      );
      this.flashMessages.danger(message);
    }
  }

  rotateCredentials = task(
    waitFor(async () => {
      try {
        const { currentPath } = this.secretMountPath;
        const { role } = this.args.model;
        await this.api.secrets.ldapRotateStaticRole(role.completeRoleName, currentPath, {});
        this.flashMessages.success('Credentials successfully rotated.');
        // The previous password is no longer valid.
        this.cachedPassword = null;
      } catch (error) {
        const { message } = await this.api.parseError(error);
        this.flashMessages.danger(`Error rotating credentials \n ${message}`);
      }
    })
  );
}
