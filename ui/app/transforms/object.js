import $ from 'jquery';
import DS from 'ember-data';
/*
  DS.attr('object')
*/
export default DS.Transform.extend({
  deserialize: function(value) {
    if (!$.isPlainObject(value)) {
      return {};
    } else {
      return value;
    }
  },
  serialize: function(value) {
    if (!$.isPlainObject(value)) {
      return {};
    } else {
      return value;
    }
  },
});
