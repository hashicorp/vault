/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { Factory, trait } from 'miragejs';

const generated_role_rules = `rules:
- apiGroups: [""]
  resources: ["secrets", "services"]
  verbs: ["get", "watch", "list", "create", "delete", "deletecollection", "patch", "update"]
`;
const name_template = '{{.FieldName | lowercase}}';
const extra_annotations = { foo: 'bar', baz: 'qux' };
const extra_labels = { foobar: 'baz', barbaz: 'foo' };

export default Factory.extend({
  name: (i) => `role-${i}`,
  allowed_kubernetes_namespaces: '*',
  allowed_kubernetes_namespace_selector: '',
  token_max_ttl: 86400,
  token_default_ttl: 600,
  service_account_name: 'default',
  kubernetes_role_name: '',
  kubernetes_role_type: 'Role',
  generated_role_rules: '',
  name_template: '',
  extra_annotations: null,
  extra_labels: null,

  afterCreate(record) {
    // only one of these three props can be defined
    if (record.generated_role_rules) {
      record.service_account_name = null;
      record.kubernetes_role_name = null;
    } else if (record.kubernetes_role_name) {
      record.service_account_name = null;
      record.generated_role_rules = null;
    } else if (record.service_account_name) {
      record.generated_role_rules = null;
      record.kubernetes_role_name = null;
    }
  },
  withRoleName: trait({
    service_account_name: null,
    generated_role_rules: null,
    kubernetes_role_name: 'vault-k8s-secrets-role',
    extra_annotations,
    name_template,
  }),
  withRoleRules: trait({
    service_account_name: null,
    kubernetes_role_name: null,
    generated_role_rules,
    extra_annotations,
    extra_labels,
    name_template,
  }),
});
