import DS from 'ember-data';
import { computed } from '@ember/object';
import { expandAttributeMeta } from 'vault/utils/field-to-attrs';

const { attr } = DS;

export default DS.Model.extend({
  queriesAvailable: attr('boolean'),
  defaultReportMonths: attr('number', {
    label: 'Default display',
    subText: 'The number of months we’ll display in the Vault usage dashboard by default.',
  }),
  retentionMonths: attr('number', {
    label: 'Retention period',
    subText: 'The number of months of activity logs to maintain for client tracking.',
  }),
  enabled: attr('string', {
    editType: 'boolean',
    trueValue: 'On',
    falseValue: 'Off',
    label: 'Enable usage data collection',
    helpText:
      'Enable or disable client tracking. Keep in mind that disabling tracking will delete the data for the current month.',
  }),

  configAttrs: computed(function() {
    let keys = ['enabled', 'defaultReportMonths', 'retentionMonths'];
    return expandAttributeMeta(this, keys);
  }),
});
