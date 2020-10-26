import DS from 'ember-data';

export default DS.Model.extend({
  total: DS.attr('object'),
  endTime: DS.attr('string'),
  startTime: DS.attr('string'),
});
