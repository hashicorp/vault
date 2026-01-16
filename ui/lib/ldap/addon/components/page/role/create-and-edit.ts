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

import { LdapRolesCreateRouteModel } from 'ldap/routes/roles/create';
import { Breadcrumb, ValidationMap } from 'vault/vault/app-types';
import type FlashMessageService from 'vault/services/flash-messages';
import type RouterService from '@ember/routing/router-service';
import type Owner from '@ember/owner';
import type ApiService from 'vault/services/api';
import type SecretMountPath from 'vault/services/secret-mount-path';
import type {
  LdapWriteStaticRoleRequest,
  LdapWriteDynamicRoleRequest,
} from '@hashicorp/vault-client-typescript';

interface Args {
  model: LdapRolesCreateRouteModel;
  breadcrumbs: Array<Breadcrumb>;
}
interface RoleTypeOption {
  title: string;
  icon: string;
  description: string;
  value: 'static' | 'dynamic';
}

export default class LdapCreateAndEditRolePageComponent extends Component<Args> {
  @service declare readonly flashMessages: FlashMessageService;
  @service('app-router') declare readonly router: RouterService;
  @service declare readonly api: ApiService;
  @service declare readonly secretMountPath: SecretMountPath;

  @tracked declare roleType: 'static' | 'dynamic';
  @tracked modelValidations: ValidationMap | null = null;
  @tracked invalidFormMessage = '';
  @tracked error = '';
  @tracked isNew = true;

  constructor(owner: Owner, args: Args) {
    super(owner, args);

    const { staticForm, dynamicForm } = this.args.model;
    // if we have both form types then we are creating a new role
    // default to static role type
    if (staticForm && dynamicForm) {
      this.roleType = 'static';
      this.isNew = true;
    } else {
      this.roleType = staticForm ? 'static' : 'dynamic';
      this.isNew = false;
    }
  }

  get roleTypeOptions(): Array<RoleTypeOption> {
    return [
      {
        title: 'Static role',
        icon: 'user',
        description: 'Static roles map to existing users in an LDAP system.',
        value: 'static',
      },
      {
        title: 'Dynamic role',
        icon: 'folder-users',
        description: 'Dynamic roles allow Vault to create and delete a user in an LDAP system.',
        value: 'dynamic',
      },
    ];
  }

  get fields() {
    if (this.roleType === 'static') {
      return ['name', 'username', 'dn', 'rotation_period'];
    }
    return [
      'name',
      'default_ttl',
      'max_ttl',
      'username_template',
      'creation_ldif',
      'deletion_ldif',
      'rollback_ldif',
    ];
  }

  get form() {
    const { staticForm, dynamicForm } = this.args.model;
    return this.roleType === 'static' ? staticForm : dynamicForm;
  }

  save = task(
    waitFor(async (event: Event) => {
      event.preventDefault();

      const { isValid, state, invalidFormMessage, data } = this.form.toJSON();
      const { name, ...payload } = data;
      const { currentPath } = this.secretMountPath;

      this.modelValidations = isValid ? null : state;
      this.invalidFormMessage = isValid ? '' : invalidFormMessage;

      if (isValid) {
        try {
          if (this.roleType === 'static') {
            await this.api.secrets.ldapWriteStaticRole(
              name,
              currentPath,
              payload as LdapWriteStaticRoleRequest
            );
          } else {
            await this.api.secrets.ldapWriteDynamicRole(
              name,
              currentPath,
              payload as LdapWriteDynamicRoleRequest
            );
          }
          this.flashMessages.success(
            `Successfully ${this.form.isNew ? 'created' : 'updated'} the role ${name}`
          );
          this.router.transitionTo(
            'vault.cluster.secrets.backend.ldap.roles.role.details',
            this.roleType,
            name
          );
        } catch (error) {
          const { message } = await this.api.parseError(
            error,
            'Error saving role. Please try again or contact support.'
          );
          this.error = message;
        }
      }
    })
  );

  @action
  onTypeChange(option: RoleTypeOption) {
    this.roleType = option.value;
    this.modelValidations = null;
    this.invalidFormMessage = '';
  }

  @action
  cancel() {
    this.router.transitionTo('vault.cluster.secrets.backend.ldap.roles');
  }
}
