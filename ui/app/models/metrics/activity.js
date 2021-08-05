import Model, { attr } from '@ember-data/model';

export default Model.extend({
  total: attr('object'),
  endTime: attr('string'),
  startTime: attr('string'),
});
