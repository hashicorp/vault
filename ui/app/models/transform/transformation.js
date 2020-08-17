import DS from 'ember-data';

export default DS.Model.extend({
  name: DS.attr('string'),
  template: DS.belongsTo('transform/template'),
  roles: DS.belongsTo('transform/role'),
});
