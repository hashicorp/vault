import { computed } from '@ember/object';
import DS from 'ember-data';
import { expandAttributeMeta } from 'vault/utils/field-to-attrs';
const { attr } = DS;
const CREATE_FIELDS = ['username', 'ip'];

const DISPLAY_FIELDS = ['username', 'ip', 'key', 'keyType', 'port'];
export default DS.Model.extend({
  role: attr('object', {
    readOnly: true,
  }),
  ip: attr('string', {
    label: 'IP Address',
  }),
  username: attr('string'),
  key: attr('string'),
  keyType: attr('string'),
  port: attr('number'),
  attrs: computed('key', function() {
    let keys = this.get('key') ? DISPLAY_FIELDS.slice(0) : CREATE_FIELDS.slice(0);
    return expandAttributeMeta(this, keys);
  }),
  toCreds: computed('key', function() {
    // todo: would this be better copied as an SSH command?
    return this.get('key');
  }),
});
