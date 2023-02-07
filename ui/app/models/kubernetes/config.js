/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Model, { attr } from '@ember-data/model';
import { withFormFields } from 'vault/decorators/model-form-fields';

@withFormFields(['kubernetesHost', 'serviceAccountJwt', 'kubernetesCaCert'])
export default class KubernetesConfigModel extends Model {
  @attr('string') backend; // dynamic path of secret -- set on response from value passed to queryRecord
  @attr('string', {
    label: 'Kubernetes host',
    subText:
      'Kubernetes API URL to connect to. Defaults to https://$KUBERNETES_SERVICE_HOST:$KUBERNETES_SERVICE_PORT if those environment variables are set.',
  })
  kubernetesHost;
  @attr('string', {
    label: 'Service account JWT',
    subText:
      'The JSON web token of the service account used by the secret engine to manage Kubernetes roles. Defaults to the local pod’s JWT if found.',
  })
  serviceAccountJwt;
  @attr('string', {
    label: 'Kubernetes CA Certificate',
    subText:
      'PEM-encoded CA certificate to use by the secret engine to verify the Kubernetes API server certificate. Defaults to the local pod’s CA if found.',
    editType: 'textarea',
  })
  kubernetesCaCert;
  @attr('boolean', { defaultValue: false }) disableLocalCaJwt;
}
