import Model, { attr } from '@ember-data/model';
import { withFormFields } from 'vault/decorators/model-form-fields';
import { withModelValidations } from 'vault/decorators/model-validations';

const validations = {
  kubernetesHost: [
    {
      validator: (model) => (model.disableLocalCaJwt && !model.kubernetesHost ? false : true),
      message: 'Kubernetes host is required',
    },
  ],
};
@withModelValidations(validations)
@withFormFields(['kubernetesHost', 'serviceAccountJwt', 'kubernetesCaCert'])
export default class KubernetesConfigModel extends Model {
  @attr('string') backend; // dynamic path of secret -- set on response from value passed to queryRecord
  @attr('string', {
    label: 'Kubernetes host',
    subText: 'Kubernetes API URL to connect to.',
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
