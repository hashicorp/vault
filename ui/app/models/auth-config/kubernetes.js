import Ember from 'ember';
import DS from 'ember-data';

import AuthConfig from '../auth-config';
import fieldToAttrs from 'vault/utils/field-to-attrs';

const { attr } = DS;
const { computed } = Ember;

export default AuthConfig.extend({
  kubernetesHost: attr('string', {
    label: 'Kubernetes Host',
    helpText:
      'Host must be a host string, a host:port pair, or a URL to the base of the Kubernetes API server',
  }),

  kubernetesCaCert: attr('string', {
    label: 'Kubernetes CA Certificate',
    editType: 'file',
    helpText: 'PEM encoded CA cert for use by the TLS client used to talk with the Kubernetes API',
  }),

  tokenReviewerJwt: attr('string', {
    label: 'Token Reviewer JWT',
    helpText:
      'A service account JWT used to access the TokenReview API to validate other JWTs during login. If not set the JWT used for login will be used to access the API',
  }),

  pemKeys: attr({
    label: 'Service account verification keys',
    editType: 'stringArray',
  }),

  fieldGroups: computed(function() {
    const groups = [
      {
        default: ['kubernetesHost', 'kubernetesCaCert'],
      },
      {
        'Kubernetes Options': ['tokenReviewerJwt', 'pemKeys'],
      },
    ];
    return fieldToAttrs(this, groups);
  }),
});
