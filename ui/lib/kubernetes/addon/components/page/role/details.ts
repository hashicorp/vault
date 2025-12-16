/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { service } from '@ember/service';

import type { KubernetesRole } from 'vault/vault/secrets/kubernetes';
import type { Breadcrumb } from 'vault/app-types';
import type RouterService from '@ember/routing/router-service';
import type FlashMessageService from 'vault/services/flash-messages';
import type ApiService from 'vault/services/api';
import type SecretMountPath from 'vault/services/secret-mount-path';
import { toLabel } from 'core/helpers/to-label';

/**
 * @module RoleDetailsPage
 * RoleDetailsPage component is a child component for create and edit role pages.
 *
 * @param {object} role - kubernetes role
 * @param {array} breadcrumbs - breadcrumbs as an array of objects that contain label and route
 */

interface Args {
  role: KubernetesRole;
  breadcrumbs: Array<Breadcrumb>;
}

export default class RoleDetailsPageComponent extends Component<Args> {
  @service('app-router') declare readonly router: RouterService;
  @service declare readonly flashMessages: FlashMessageService;
  @service declare readonly api: ApiService;
  @service declare readonly secretMountPath: SecretMountPath;

  label = (field: string) => {
    return (
      {
        name: 'Role name',
        allowed_kubernetes_namespaces: 'Allowed Kubernetes namespaces',
        token_max_ttl: 'Max lease TTL',
        token_default_ttl: 'Default lease TTL',
      }[field] || toLabel([field])
    );
  };

  get displayFields() {
    const { role } = this.args;
    const fields = [
      'name',
      'service_account_name',
      'kubernetes_role_type',
      'kubernetes_role_name',
      'allowed_kubernetes_namespaces',
      'token_max_ttl',
      'token_default_ttl',
      'name_template',
    ];
    // return different fields based on generation preference selected during creation
    const hiddenFieldIndices: number[] = [];
    if (role.service_account_name) {
      // hide kubernetes_role_type, kubernetes_role_name and name_template
      hiddenFieldIndices.push(2, 3, 7);
    } else if (role.kubernetes_role_name) {
      // hide service_account_name
      hiddenFieldIndices.push(1);
    } else if (role.generated_role_rules) {
      // hide service_account_name and kubernetes_role_name
      hiddenFieldIndices.push(1, 3);
    }

    return fields.filter((_field, index) => !hiddenFieldIndices.includes(index));
  }

  get extraFields() {
    const fields = [];
    if (this.args.role.extra_annotations) {
      fields.push({ label: 'Annotations', key: 'extra_annotations' });
    }
    if (this.args.role.extra_labels) {
      fields.push({ label: 'Labels', key: 'extra_labels' });
    }
    return fields;
  }

  @action
  async delete() {
    try {
      await this.api.secrets.kubernetesDeleteRole(this.args.role.name, this.secretMountPath.currentPath);
      this.router.transitionTo('vault.cluster.secrets.backend.kubernetes.roles');
    } catch (error) {
      const { message } = await this.api.parseError(
        error,
        'Unable to delete role. Please try again or contact support'
      );
      this.flashMessages.danger(message);
    }
  }
}
