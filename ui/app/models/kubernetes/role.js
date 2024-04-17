/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Model, { attr } from '@ember-data/model';
import { withModelValidations } from 'vault/decorators/model-validations';
import { withFormFields } from 'vault/decorators/model-form-fields';
import lazyCapabilities, { apiPath } from 'vault/macros/lazy-capabilities';
import { tracked } from '@glimmer/tracking';

const validations = {
  name: [{ type: 'presence', message: 'Name is required' }],
};
const formFieldProps = [
  'name',
  'serviceAccountName',
  'kubernetesRoleType',
  'kubernetesRoleName',
  'allowedKubernetesNamespaces',
  'tokenMaxTtl',
  'tokenDefaultTtl',
  'nameTemplate',
];

@withModelValidations(validations)
@withFormFields(formFieldProps)
export default class KubernetesRoleModel extends Model {
  @attr('string') backend; // dynamic path of secret -- set on response from value passed to queryRecord
  @attr('string', {
    label: 'Role name',
    subText: 'The role’s name in Vault.',
  })
  name;

  @attr('string', {
    label: 'Service account name',
    subText: 'Vault will use the default template when generating service accounts, roles and role bindings.',
  })
  serviceAccountName;

  @attr('string', {
    label: 'Kubernetes role type',
    editType: 'radio',
    possibleValues: ['Role', 'ClusterRole'],
  })
  kubernetesRoleType;

  @attr('string', {
    label: 'Kubernetes role name',
    subText: 'Vault will use the default template when generating service accounts, roles and role bindings.',
  })
  kubernetesRoleName;

  @attr('string', {
    label: 'Allowed Kubernetes namespaces',
    subText:
      'A list of the valid Kubernetes namespaces in which this role can be used for creating service accounts. If set to "*" all namespaces are allowed.',
  })
  allowedKubernetesNamespaces;

  @attr({
    label: 'Max Lease TTL',
    editType: 'ttl',
  })
  tokenMaxTtl;

  @attr({
    label: 'Default Lease TTL',
    editType: 'ttl',
  })
  tokenDefaultTtl;

  @attr('string', {
    label: 'Name template',
    editType: 'optionalText',
    defaultSubText:
      'Vault will use the default template when generating service accounts, roles and role bindings.',
    subText: 'Vault will use the default template when generating service accounts, roles and role bindings.',
  })
  nameTemplate;

  @attr extraAnnotations;
  @attr extraLabels;

  @attr('string') generatedRoleRules;

  @tracked _generationPreference;
  get generationPreference() {
    // when the user interacts with the radio cards the value will be set to the pseudo prop which takes precedence
    if (this._generationPreference) {
      return this._generationPreference;
    }
    // for existing roles, default the value based on which model prop has value -- only one can be set
    let pref = null;
    if (this.serviceAccountName) {
      pref = 'basic';
    } else if (this.kubernetesRoleName) {
      pref = 'expanded';
    } else if (this.generatedRoleRules) {
      pref = 'full';
    }
    return pref;
  }
  set generationPreference(pref) {
    // unset model props specific to filteredFormFields when changing preference
    // only one of service_account_name, kubernetes_role_name or generated_role_rules can be set
    const props = {
      basic: ['kubernetesRoleType', 'kubernetesRoleName', 'generatedRoleRules', 'nameTemplate'],
      expanded: ['serviceAccountName', 'generatedRoleRules'],
      full: ['serviceAccountName', 'kubernetesRoleName'],
    }[pref];
    props.forEach((prop) => (this[prop] = null));
    this._generationPreference = pref;
  }

  get filteredFormFields() {
    // return different form fields based on generationPreference
    const hiddenFieldIndices = {
      basic: [2, 3, 7], // kubernetesRoleType, kubernetesRoleName and nameTemplate
      expanded: [1], // serviceAccountName
      full: [1, 3], // serviceAccountName and kubernetesRoleName
    }[this.generationPreference];

    return hiddenFieldIndices
      ? this.formFields.filter((field, index) => !hiddenFieldIndices.includes(index))
      : null;
  }

  @lazyCapabilities(apiPath`${'backend'}/roles/${'name'}`, 'backend', 'name') rolePath;
  @lazyCapabilities(apiPath`${'backend'}/creds/${'name'}`, 'backend', 'name') credsPath;
  @lazyCapabilities(apiPath`${'backend'}/roles`, 'backend') rolesPath;

  get canCreate() {
    return this.rolePath.get('canCreate');
  }
  get canDelete() {
    return this.rolePath.get('canDelete');
  }
  get canEdit() {
    return this.rolePath.get('canUpdate');
  }
  get canRead() {
    return this.rolePath.get('canRead');
  }
  get canList() {
    return this.rolesPath.get('canList');
  }
  get canGenerateCreds() {
    return this.credsPath.get('canCreate');
  }
}
