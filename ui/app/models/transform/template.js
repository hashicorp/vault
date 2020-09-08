import DS from 'ember-data';

export default DS.Model.extend({
  name: DS.attr('string'),
  alphabet: DS.belongsTo('transform/alphabet'),
  transformations: DS.hasMany('transformation'),
});
