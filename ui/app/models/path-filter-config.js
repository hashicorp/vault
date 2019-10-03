import DS from 'ember-data';
const { attr } = DS;

export default DS.Model.extend({
  mode: attr('string', {
    defaultValue: 'whitelist',
  }),
  paths: attr('array', {
    defaultValue: function() {
      return [];
    },
  }),
});
