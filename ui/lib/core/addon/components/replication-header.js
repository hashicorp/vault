import Component from '@ember/component';
import layout from '../templates/components/replication-header';

export default Component.extend({
  layout,
  classNames: ['replication-header'],
  isSecondary: null,
  secondaryId: null,
  isSummaryDashboard: false,
  'data-test-replication-header': true,
});
