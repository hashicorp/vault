import Model, { attr } from '@ember-data/model';

// TODO CMB add model attrs here

export default Model.extend({
  total: attr('object'),
  byNamespace: attr('array'),
  endTime: attr('string'),
  startTime: attr('string'),
  clients: attr('number'),
  distinct_entities: attr('number'),
  non_entity_tokens: attr('number'),
});
