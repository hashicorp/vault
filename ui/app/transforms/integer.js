import { isEmpty } from '@ember/utils';
import DS from 'ember-data';

/*
  DS.attr('integer')
*/
export default DS.Transform.extend({
  deserialize: function(value) {
    if (isEmpty(value)) {
      return null;
    } else {
      return value;
    }
  },
  serialize: function(value) {
    if (isEmpty(value)) {
      return null;
    } else {
      return value;
    }
  },
});
