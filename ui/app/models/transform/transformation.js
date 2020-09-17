import Model, { attr, belongsTo } from '@ember-data/model';

export default Model.extend({
  name: attr('string'),
  template: belongsTo('transform/template'),
  roles: belongsTo('transform/role'),
});
