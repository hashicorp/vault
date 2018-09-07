import { computed } from '@ember/object';
import DS from 'ember-data';
import { expandAttributeMeta } from 'vault/utils/field-to-attrs';
const { attr } = DS;
const CREATE_FIELDS = ['ttl', 'roleArn'];

const DISPLAY_FIELDS = ['accessKey', 'secretKey', 'securityToken', 'leaseId', 'renewable', 'leaseDuration'];
export default DS.Model.extend({
  helpText:
    'For Vault roles that have a credential type of iam_user, these attributes are optional and you may simply submit the form.',
  role: attr('object', {
    readOnly: true,
  }),

  roleArn: attr('string', {
    label: 'Role ARN',
    helpText:
      'The ARN of the role to assume if credential_type on the Vault role is assumed_role. Optional if the role has a single role ARN; required otherwise.',
  }),
  ttl: attr({
    editType: 'ttl',
    defaultValue: '3600s',
    label: 'TTL',
    helpText:
      'Specifies the TTL for the use of the STS token. Valid only when credential_type is assumed_role or federation_token.',
  }),
  leaseId: attr('string'),
  renewable: attr('boolean'),
  leaseDuration: attr('number'),
  accessKey: attr('string'),
  secretKey: attr('string'),
  securityToken: attr('string'),

  attrs: computed('accessKey', 'securityToken', function() {
    let keys =
      this.get('accessKey') || this.get('securityToken') ? DISPLAY_FIELDS.slice(0) : CREATE_FIELDS.slice(0);
    return expandAttributeMeta(this, keys);
  }),

  toCreds: computed('accessKey', 'secretKey', 'securityToken', 'leaseId', function() {
    const props = this.getProperties('accessKey', 'secretKey', 'securityToken', 'leaseId');
    const propsWithVals = Object.keys(props).reduce((ret, prop) => {
      if (props[prop]) {
        ret[prop] = props[prop];
        return ret;
      }
      return ret;
    }, {});
    return JSON.stringify(propsWithVals, null, 2);
  }),
});
