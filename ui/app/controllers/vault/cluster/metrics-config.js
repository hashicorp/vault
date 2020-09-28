import Controller from '@ember/controller';
import { computed } from '@ember/object';

export default Controller.extend({
  infoRows: computed(function() {
    return [
      {
        label: 'Usage data collection',
        helperText: 'Enable or disable collecting data to track clients.',
        valueKey: 'enabled',
      },
      {
        label: 'Retention period',
        helperText: 'The number of months of activity logs to maintain for  client tracking.',
        valueKey: 'retentionMonths',
      },
      {
        label: 'Default display',
        helperText: 'The number of months we’ll display in the Vault usage dashboard by default.',
        valueKey: 'defaultReportMonths',
      },
    ];
  }),
});
