/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Form from 'vault/forms/form';
import FormField from 'vault/utils/forms/field';
import { tracked } from '@glimmer/tracking';

import type { KubernetesWriteRoleRequest } from '@hashicorp/vault-client-typescript';
import type { Validations } from 'vault/app-types';

type KubernetesRoleFormData = KubernetesWriteRoleRequest & {
  name?: string;
};
type GenerationPreference = 'basic' | 'expanded' | 'full' | null;

export default class KubernetesRoleForm extends Form<KubernetesRoleFormData> {
  @tracked declare _generationPreference: GenerationPreference;

  get generationPreference() {
    // when the user interacts with the radio cards the value will be set to the pseudo prop which takes precedence
    if (this._generationPreference) {
      return this._generationPreference;
    }
    // for existing roles, default the value based on which model prop has value -- only one can be set
    let pref: GenerationPreference = null;
    if (this.data.service_account_name) {
      pref = 'basic';
    } else if (this.data.kubernetes_role_name) {
      pref = 'expanded';
    } else if (this.data.generated_role_rules) {
      pref = 'full';
    }
    return pref;
  }
  set generationPreference(pref: GenerationPreference) {
    // unset model props specific to filteredFormFields when changing preference
    // only one of service_account_name, kubernetes_role_name or generated_role_rules can be sent in payload
    if (pref) {
      const props = {
        basic: ['kubernetes_role_type', 'kubernetes_role_name', 'generated_role_rules', 'name_template'],
        expanded: ['service_account_name', 'generated_role_rules'],
        full: ['service_account_name', 'kubernetes_role_name'],
      }[pref];
      props.forEach((prop) => {
        delete this.data[prop as keyof typeof this.data];
      });
    }
    this._generationPreference = pref;
  }

  get formFields() {
    const fields = [
      new FormField('name', 'string', {
        label: 'Role name',
        subText: 'The role’s name in Vault.',
      }),
      new FormField('service_account_name', 'string', {
        label: 'Service account name',
        subText:
          'Vault will use the default template when generating service accounts, roles and role bindings.',
      }),
      new FormField('kubernetes_role_type', 'string', {
        label: 'Kubernetes role type',
        editType: 'radio',
        possibleValues: ['Role', 'ClusterRole'],
      }),
      new FormField('kubernetes_role_name', 'string', {
        label: 'Kubernetes role name',
        subText:
          'Vault will use the default template when generating service accounts, roles and role bindings.',
      }),
      new FormField('allowed_kubernetes_namespaces', 'string', {
        label: 'Allowed Kubernetes namespaces',
        subText:
          'A list of the valid Kubernetes namespaces in which this role can be used for creating service accounts. If set to "*" all namespaces are allowed.',
      }),
      new FormField('token_max_ttl', 'string', {
        label: 'Max Lease TTL',
        editType: 'ttl',
      }),
      new FormField('token_default_ttl', 'string', {
        label: 'Default Lease TTL',
        editType: 'ttl',
      }),
      new FormField('name_template', 'string', {
        label: 'Name template',
        editType: 'optionalText',
        defaultSubText:
          'Vault will use the default template when generating service accounts, roles and role bindings.',
        subText:
          'Vault will use the default template when generating service accounts, roles and role bindings.',
      }),
    ];
    // return different form fields based on generationPreference
    if (this.generationPreference) {
      const hiddenFieldIndices = {
        basic: [2, 3, 7], // kubernetes_role_type, kubernetes_role_name and name_template
        expanded: [1], // service_account_name
        full: [1, 3], // service_account_name and kubernetes_role_name
      }[this.generationPreference];

      return hiddenFieldIndices ? fields.filter((_field, index) => !hiddenFieldIndices.includes(index)) : [];
    }
    return [];
  }

  validations: Validations = {
    name: [{ type: 'presence', message: 'Name is required' }],
  };
}
