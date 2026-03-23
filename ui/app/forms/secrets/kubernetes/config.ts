/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Form from 'vault/forms/form';
import FormField from 'vault/utils/forms/field';

import type { KubernetesConfigureRequest } from '@hashicorp/vault-client-typescript';
import type { Validations } from 'vault/app-types';

export default class KubernetesConfigForm extends Form<KubernetesConfigureRequest> {
  formFields = [
    new FormField('kubernetes_host', 'string', {
      label: 'Kubernetes host',
      subText: 'Kubernetes API URL to connect to.',
    }),
    new FormField('service_account_jwt', 'string', {
      label: 'Service account JWT',
      subText:
        'The JSON web token of the service account used by the secret engine to manage Kubernetes roles. Defaults to the local pod’s JWT if found.',
    }),
    new FormField('kubernetes_ca_cert', 'string', {
      label: 'Kubernetes CA Certificate',
      subText:
        'PEM-encoded CA certificate to use by the secret engine to verify the Kubernetes API server certificate. Defaults to the local pod’s CA if found.',
      editType: 'textarea',
    }),
  ];

  validations: Validations = {
    kubernetes_host: [
      {
        validator: (data: KubernetesConfigForm['data']) =>
          data.disable_local_ca_jwt && !data.kubernetes_host ? false : true,
        message: 'Kubernetes host is required',
      },
    ],
  };

  toJSON() {
    // ensure that values from a previous manual configuration are unset
    const { disable_local_ca_jwt } = this.data;
    const data = disable_local_ca_jwt ? this.data : { disable_local_ca_jwt };
    return super.toJSON(data);
  }
}
