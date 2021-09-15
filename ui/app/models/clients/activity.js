import Model, { attr } from '@ember-data/model';

export default Model.extend({
  total: attr('object'),
  byNamespace: attr('array'),
  endTime: attr('string'),
  startTime: attr('string'),
});
