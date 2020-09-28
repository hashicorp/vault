import DS from 'ember-data';

const { attr } = DS;

export default DS.Model.extend({
  queriesAvailable: attr('boolean'),
  defaultMonths: attr('number'),
  retentionMonths: attr('number'),
  enabled: attr('boolean'),
});
