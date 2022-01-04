import Transform from '@ember-data/serializer/transform';
import { typeOf } from '@ember/utils';
/*
  DS.attr('object')
*/
export default Transform.extend({
  deserialize: function (value) {
    if (typeOf(value) !== 'object') {
      return {};
    } else {
      return value;
    }
  },
  serialize: function (value) {
    if (typeOf(value) !== 'object') {
      return {};
    } else {
      return value;
    }
  },
});
