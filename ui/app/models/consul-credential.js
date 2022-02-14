import Model, { attr } from '@ember-data/model';
import { computed } from '@ember/object';
import { expandAttributeMeta } from 'vault/utils/field-to-attrs';
const CREDENTIAL_TYPES = [
  {
    value: 'client',
    displayName: 'Service Token',
  },
  {
    value: 'management',
    displayName: 'Management Token',
  },
];

const DISPLAY_FIELDS = ['accessor', 'token', 'local'];
export default Model.extend({
  helpText:
    'There are no inputs, just submit the form.',
  role: attr('object', {
    readOnly: true,
  }),

  credentialType: attr('string', {
    defaultValue: 'client',
    possibleValues: CREDENTIAL_TYPES,
    readOnly: true,
  }),

  leaseId: attr('string'),
  renewable: attr('boolean'),
  leaseDuration: attr('number'),
  local: attr('boolean'),
  accessor: attr('string'),
  token: attr('string'),

  attrs: computed('accessor', 'token', 'local', function() {
    if (this.accessKey || this.securityToken) {
      return expandAttributeMeta(this, DISPLAY_FIELDS.slice(0));
    }
    return expandAttributeMeta(this, []);
  }),

  toCreds: computed('accessor', 'token', 'local', 'leaseId', function() {
    const props = {
      accessor: this.accessor,
      token: this.token,
      local: this.local,
      leaseId: this.leaseId
    };
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
