import DS from 'ember-data';
const { attr } = DS;

export default DS.Model.extend({
  mode: attr('string'),
  paths: attr('array', {
    defaultValue: function() {
      return [];
    },
  }),
});
