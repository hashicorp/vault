import Ember from 'ember';
import DS from 'ember-data';

const { decamelize } = Ember.String;

export default DS.RESTSerializer.extend({
  keyForAttribute: function(attr) {
    return decamelize(attr);
  },
});
