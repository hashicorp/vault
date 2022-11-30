import { Factory } from 'ember-cli-mirage';

export default Factory.extend({
  name: (i) => `role-${i}`,
  allowed_kubernetes_namespaces: () => ['*'],
  allowed_kubernetes_namespace_selector: '',
  token_max_ttl: 86400,
  token_default_ttl: 0,
  service_account_name: 'default',
  kubernetes_role_name: '',
  kubernetes_role_type: 'Role',
  generated_role_rules: '',
  name_template: '',
  extra_annotations: null,
  extra_labels: null,
});
