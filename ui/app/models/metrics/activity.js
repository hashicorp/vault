import Model, { attr } from '@ember-data/model';

export default Model.extend({
  total: attr('object'),
  by_namespace: attr('array'),
  endTime: attr('string'),
  startTime: attr('string'),
});
