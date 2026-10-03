/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

export type KubernetesRole = {
  allowed_kubernetes_namespace_selector: string;
  allowed_kubernetes_namespaces: string[];
  extra_annotations: Record<string, string>;
  extra_labels: Record<string, string>;
  generated_role_rules: string;
  kubernetes_role_name: string;
  kubernetes_role_type: string;
  name: string;
  name_template: string;
  service_account_name: string;
  token_default_ttl: string;
  token_max_ttl: string;
};

export type KubernetesCredentials = {
  service_account_name: string;
  service_account_namespace: string;
  service_account_token: string;
  lease_duration: number;
  lease_id: string;
};
