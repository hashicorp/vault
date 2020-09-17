import Model, { attr, belongsTo, hasMany } from '@ember-data/model';

export default Model.extend({
  name: attr('string'),
  alphabet: belongsTo('transform/alphabet'),
  transformations: hasMany('transformation'),
});
