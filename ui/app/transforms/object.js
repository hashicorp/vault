import DS from 'ember-data';
import { typeOf } from '@ember/utils';
/*
  DS.attr('object')
*/
export default DS.Transform.extend({
  deserialize: function(value) {
    if (typeOf(value) !== 'object') {
      return {};
    } else {
      return value;
    }
  },
  serialize: function(value) {
    if (typeOf(value) !== 'object') {
      return {};
    } else {
      return value;
    }
  },
});
