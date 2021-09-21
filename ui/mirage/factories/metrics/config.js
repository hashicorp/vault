import Mirage from 'ember-cli-mirage';

export default Mirage.Factory.extend({
  default_report_months: 12,
  enabled: 'default-enable',
  queries_available: false,
  retention_months: 24,
});
