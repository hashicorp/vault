import DS from 'ember-data';
import Ember from 'ember';
const { attr } = DS;
const { computed, get } = Ember;
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
  toCreds: computed('key', function() {
    // todo: would this be better copied as an SSH command?
    return this.get('key');
  }),
});
