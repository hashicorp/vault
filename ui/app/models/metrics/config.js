import DS from 'ember-data';

const { attr } = DS;

export default DS.Model.extend({
  queriesAvailable: attr('boolean'),
  defaultReportMonths: attr('number'),
  retentionMonths: attr('number'),
  enabled: attr('string'),
});
