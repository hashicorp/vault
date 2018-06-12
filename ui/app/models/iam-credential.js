import DS from 'ember-data';
import Ember from 'ember';
import { expandAttributeMeta } from 'vault/utils/field-to-attrs';
const { attr } = DS;
const { computed } = Ember;
const CREATE_FIELDS = ['ttl'];

const DISPLAY_FIELDS = ['accessKey', 'secretKey', 'securityToken', 'leaseId', 'renewable', 'leaseDuration'];
export default DS.Model.extend({
  role: attr('object', {
    readOnly: true,
  }),

  withSTS: attr('boolean', {
    readOnly: true,
  }),

  ttl: attr({
    editType: 'ttl',
    defaultValue: '1h',
  }),
  leaseId: attr('string'),
  renewable: attr('boolean'),
  leaseDuration: attr('number'),
  accessKey: attr('string'),
  secretKey: attr('string'),
  securityToken: attr('string'),

  attrs: computed('accessKey', function() {
    let keys = this.get('accessKey') ? DISPLAY_FIELDS.slice(0) : CREATE_FIELDS.slice(0);
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
