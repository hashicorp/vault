import DS from 'ember-data';
import Ember from 'ember';
const { attr } = DS;
const { computed, get } = Ember;
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
    get(this.constructor, 'attributes').forEach((meta, name) => {
      const index = keys.indexOf(name);
      if (index === -1) {
        return;
      }
      keys.replace(index, 1, {
        type: meta.type,
        name,
        options: meta.options,
      });
    });
    return keys;
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
